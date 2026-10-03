package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var (
	inspectTarget string
	inspectSearch string
)

func init() {
	inspectCmd.Flags().StringVarP(&inspectTarget, "target", "t", "", "Target entity to inspect (e.g. 'operations.getUsers' or 'types.User')")
	inspectCmd.Flags().StringVarP(&inspectSearch, "search", "s", "", "Filter operations or types by substring match")
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

		// Direct entity lookup
		if inspectTarget != "" {
			parts := strings.Split(inspectTarget, ".")
			if len(parts) < 2 {
				log.Fatalf("Invalid target '%s'. Format must be 'operations.<name>' or 'types.<name>'", inspectTarget)
			}

			domain := parts[0]
			entityName := strings.Join(parts[1:], ".")
			if domain == "types" {
				domain = "definitions"
			}

			section, ok := api[domain].(map[string]interface{})
			if !ok {
				log.Fatalf("Section '%s' not found in schema", domain)
			}

			targetEntity, exists := section[entityName]
			if !exists {
				log.Fatalf("Target '%s' not found in %s", entityName, domain)
			}

			targetJSON, _ := json.MarshalIndent(targetEntity, "", "  ")
			fmt.Println(string(targetJSON))
			return
		}

		// Filtered/Sorted Overview
		fmt.Println("Operations:")
		if operations, ok := api["operations"].(map[string]interface{}); ok && len(operations) > 0 {
			opKeys := make([]string, 0, len(operations))
			for k := range operations {
				if inspectSearch == "" || strings.Contains(strings.ToLower(k), strings.ToLower(inspectSearch)) {
					opKeys = append(opKeys, k)
				}
			}
			sort.Strings(opKeys)

			if len(opKeys) == 0 {
				fmt.Println("- None (no match)")
			} else {
				for _, opName := range opKeys {
					op, _ := operations[opName].(map[string]interface{})
					method, _ := op["method"].(string)
					path, _ := op["path"].(string)

					argNames := make([]string, 0)
					if argsMap, ok := op["arguments"].(map[string]interface{}); ok {
						for argName := range argsMap {
							argNames = append(argNames, argName)
						}
						sort.Strings(argNames)
					}

					fmt.Printf("- %s [%s %s] (arguments: %s)\n", opName, method, path, strings.Join(argNames, ", "))
				}
			}
		} else {
			fmt.Println("- None")
		}

		fmt.Println()
		fmt.Println("Types:")
		if definitions, ok := api["definitions"].(map[string]interface{}); ok && len(definitions) > 0 {
			typeKeys := make([]string, 0, len(definitions))
			for k := range definitions {
				if inspectSearch == "" || strings.Contains(strings.ToLower(k), strings.ToLower(inspectSearch)) {
					typeKeys = append(typeKeys, k)
				}
			}
			sort.Strings(typeKeys)

			if len(typeKeys) == 0 {
				fmt.Println("- None (no match)")
			} else {
				for _, typeName := range typeKeys {
					typeDef, _ := definitions[typeName].(map[string]interface{})
					propNames := make([]string, 0)
					if propsMap, ok := typeDef["properties"].(map[string]interface{}); ok {
						for propName := range propsMap {
							propNames = append(propNames, propName)
						}
						sort.Strings(propNames)
					}

					fmt.Printf("- %s (properties: %s)\n", typeName, strings.Join(propNames, ", "))
				}
			}
		} else {
			fmt.Println("- None")
		}
	},
}
