package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate [flags] TYPE SCHEMA_FILE OUTPUT_DIR",
	Short: "Generates an SDK based on a provided TypeAPI specification",
	Long:  `Generates an SDK based on a provided TypeAPI specification, more information about TypeAPI at https://typeapi.org/`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		var client = sdkClient.GetClient()

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		var generatorType = args[0]
		var schemaFile = args[1]
		var outputDir = args[2]

		jsonFile, err := readFile(schemaFile)
		if err != nil {
			log.Fatal(err)
		}

		var schema = json.RawMessage{}
		err = json.Unmarshal(jsonFile, &schema)
		if err != nil {
			log.Fatal(err)
		}

		var targetDir = filepath.Join(cwd, outputDir)
		var mapping = make(map[string]string)

		Generate(client, generatorType, schema, targetDir, sdkClient.Namespace, sdkClient.BaseUrl, mapping, sdkClient.Remove)

		fmt.Println("Generation successful!")
		os.Exit(0)
	},
}
