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

		fmt.Printf("✅ %s is a valid TypeAPI specification.\n", schemaFile)
	},
}
