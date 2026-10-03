package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		schemaFilePath := filepath.Join(cwd, schemaFile)
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
