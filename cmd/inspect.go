package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(inspectCmd)
}

var inspectCmd = &cobra.Command{
	Use:   "inspect SCHEMA_FILE",
	Short: "Inspects an existing TypeAPI specification",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		schemaFile := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		schemaFilePath := filepath.Join(cwd, schemaFile)
		fileBytes, err := os.ReadFile(schemaFilePath)
		if err != nil {
			log.Fatalf("Error reading file %s: %v\n", schemaFile, err)
		}

		var api map[string]interface{}
		if err := json.Unmarshal(fileBytes, &api); err != nil {
			log.Fatalf("Failed to unmarshal TypeAPI specification: %v\n", err)
		}

		fmt.Println("Operations:")
		if operations, ok := api["operations"].(map[string]interface{}); ok && len(operations) > 0 {
			for opName, opRaw := range operations {
				op, ok := opRaw.(map[string]interface{})
				if !ok {
					continue
				}

				method, _ := op["method"].(string)
				path, _ := op["path"].(string)

				argNames := make([]string, 0)
				if argsMap, ok := op["arguments"].(map[string]interface{}); ok {
					for argName := range argsMap {
						argNames = append(argNames, argName)
					}
				}

				fmt.Printf("- %s [%s %s] (arguments: %s)\n", opName, method, path, strings.Join(argNames, ", "))
			}
		} else {
			fmt.Println("- None")
		}

		fmt.Println()
		fmt.Println("Types:")
		if definitions, ok := api["definitions"].(map[string]interface{}); ok && len(definitions) > 0 {
			for typeName, typeRaw := range definitions {
				typeDef, ok := typeRaw.(map[string]interface{})
				if !ok {
					continue
				}

				propNames := make([]string, 0)
				if propsMap, ok := typeDef["properties"].(map[string]interface{}); ok {
					for propName := range propsMap {
						propNames = append(propNames, propName)
					}
				}

				fmt.Printf("- %s (properties: %s)\n", typeName, strings.Join(propNames, ", "))
			}
		} else {
			fmt.Println("- None")
		}
	},
}
