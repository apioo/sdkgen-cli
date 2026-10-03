package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/apioo/sdkgen-cli/spec"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate SCHEMA_FILE",
	Short: "Validates an existing TypeAPI specification",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		schemaFile := args[0]

		schemaFilePath, err := filepath.Abs(schemaFile)
		if err != nil {
			log.Fatal(err)
		}

		fileBytes, err := readFile(schemaFilePath)
		if err != nil {
			log.Fatalf("Error reading file %s: %v\n", schemaFile, err)
		}

		var schemaDoc interface{}
		if err := json.Unmarshal([]byte(spec.TypeAPISchema), &schemaDoc); err != nil {
			log.Fatalf("Failed to unmarshal embedded TypeAPI schema: %v\n", err)
		}

		compiler := jsonschema.NewCompiler()
		if err := compiler.AddResource("schema.json", schemaDoc); err != nil {
			log.Fatalf("Failed to register schema resource: %v\n", err)
		}

		schema, err := compiler.Compile("schema.json")
		if err != nil {
			log.Fatalf("Failed to compile TypeAPI validator: %v\n", err)
		}

		var targetDoc interface{}
		if err := json.Unmarshal(fileBytes, &targetDoc); err != nil {
			log.Fatalf("Invalid JSON file syntax in %s: %v\n", schemaFile, err)
		}

		// 1. JSON Schema Structural Validation
		if err := schema.Validate(targetDoc); err != nil {
			fmt.Printf("❌ Validation failed for %s:\n", schemaFile)
			if validationErr, ok := err.(*jsonschema.ValidationError); ok {
				for _, cause := range validationErr.Causes {
					fmt.Printf("- %s\n", cause.Error())
				}
			} else {
				fmt.Println(err)
			}
			os.Exit(1)
		}

		// 2. Logical Reference Validation
		if apiMap, ok := targetDoc.(map[string]interface{}); ok {
			var missingRefs []string

			// Collect available definitions
			definitions := make(map[string]bool)
			if defsMap, ok := apiMap["definitions"].(map[string]interface{}); ok {
				for defName := range defsMap {
					definitions[defName] = true
				}
			}

			// Recursively collect and verify reference types
			validateReferences(targetDoc, definitions, &missingRefs)

			if len(missingRefs) > 0 {
				fmt.Printf("❌ Logical validation failed for %s:\n", schemaFile)
				for _, refErr := range missingRefs {
					fmt.Printf("- %s\n", refErr)
				}
				os.Exit(1)
			}

			// 3. Convention Validation
			if problems := validateConventions(apiMap); len(problems) > 0 {
				fmt.Printf("❌ Convention validation failed for %s:\n", schemaFile)
				for _, problem := range problems {
					fmt.Printf("- %s\n", problem)
				}
				os.Exit(1)
			}
		}

		fmt.Printf("✅ %s is a valid TypeAPI specification.\n", schemaFile)
	},
}

// Recursively walks the TypeAPI JSON tree to inspect reference types
func validateReferences(node interface{}, definitions map[string]bool, missingRefs *[]string) {
	switch v := node.(type) {
	case map[string]interface{}:
		// Check if current node is a reference type object: {"type": "reference", "target": "..."}
		if typeStr, ok := v["type"].(string); ok && typeStr == "reference" {
			if target, ok := v["target"].(string); ok && target != "" {
				if !definitions[target] {
					*missingRefs = append(*missingRefs, fmt.Sprintf("Referenced type '%s' does not exist in definitions", target))
				}
			}
		}

		for _, val := range v {
			validateReferences(val, definitions, missingRefs)
		}

	case []interface{}:
		for _, elem := range v {
			validateReferences(elem, definitions, missingRefs)
		}
	}
}

var (
	requiredOperationKeys = []string{"path", "method", "return"}
	operationNamePattern  = regexp.MustCompile(`^[a-z][A-Za-z0-9]*(\.[a-z][A-Za-z0-9]*)+$`)
	typeNamePattern       = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*(_[A-Z0-9][A-Za-z0-9]*)*$`)
	pathParamPattern      = regexp.MustCompile(`:([A-Za-z_][A-Za-z0-9_]*)`)
)

// Checks naming and operation conventions which are not covered by the JSON schema
func validateConventions(apiMap map[string]interface{}) []string {
	var problems []string

	operations, _ := apiMap["operations"].(map[string]interface{})
	for _, name := range sortedKeys(operations) {
		operation, ok := operations[name].(map[string]interface{})
		if !ok {
			continue
		}

		where := "operations." + name
		if !operationNamePattern.MatchString(name) {
			problems = append(problems, fmt.Sprintf("%s: name should be 'resource.action' in camelCase (e.g. records.getAll)", where))
		}

		for _, key := range requiredOperationKeys {
			if _, ok := operation[key]; !ok {
				problems = append(problems, fmt.Sprintf("%s: missing '%s'", where, key))
			}
		}

		path, _ := operation["path"].(string)
		if strings.ContainsAny(path, "{}") {
			problems = append(problems, fmt.Sprintf("%s: path uses {param} style, use :param instead (%s)", where, path))
		}

		// Every path parameter needs a path argument and vice versa
		arguments, _ := operation["arguments"].(map[string]interface{})
		pathParams := make(map[string]bool)
		for _, match := range pathParamPattern.FindAllStringSubmatch(path, -1) {
			param := match[1]
			pathParams[param] = true

			argument, ok := arguments[param].(map[string]interface{})
			if !ok {
				problems = append(problems, fmt.Sprintf("%s: path parameter ':%s' has no argument", where, param))
			} else if argument["in"] != "path" {
				problems = append(problems, fmt.Sprintf("%s: argument '%s' must have \"in\": \"path\"", where, param))
			}
		}

		var bodies []string
		for _, argName := range sortedKeys(arguments) {
			argument, ok := arguments[argName].(map[string]interface{})
			if !ok {
				continue
			}

			switch argument["in"] {
			case "path":
				if !pathParams[argName] {
					problems = append(problems, fmt.Sprintf("%s: path argument '%s' does not appear in path", where, argName))
				}
			case "body":
				bodies = append(bodies, argName)
			}
		}

		if len(bodies) > 1 {
			problems = append(problems, fmt.Sprintf("%s: more than one body argument (%s)", where, strings.Join(bodies, ", ")))
		}

		method, _ := operation["method"].(string)
		if len(bodies) > 0 && (method == "GET" || method == "DELETE") {
			problems = append(problems, fmt.Sprintf("%s: %s operation has a body argument", where, method))
		}
	}

	definitions, _ := apiMap["definitions"].(map[string]interface{})
	for _, name := range sortedKeys(definitions) {
		if !typeNamePattern.MatchString(name) {
			problems = append(problems, fmt.Sprintf("types.%s: name should be Pascal_Snake case (e.g. Record_Collection)", name))
		}
	}

	return problems
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
