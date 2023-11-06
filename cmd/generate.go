package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
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

		var generatorType = args[0]
		var schemaFile = args[1]
		var outputDir = args[2]

		jsonFile, err := os.Open(schemaFile)
		if err != nil {
			log.Fatal(err)
		}

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Fatal(err)
		}

		Generate(client, generatorType, byteValue, outputDir, sdkClient.Namespace, sdkClient.BaseUrl)

		fmt.Println("Generation successful!")
		os.Exit(0)
	},
}
