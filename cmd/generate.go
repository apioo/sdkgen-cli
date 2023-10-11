package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/apioo/sdkgen-cli/sdk"
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

		stat, err := os.Stat(outputDir)
		if err != nil {
			log.Fatal("Provided output directory does not exist")
		}

		if !stat.IsDir() {
			log.Fatal("Provided output directory does not exist")
		}

		jsonFile, err := os.Open(schemaFile)
		if err != nil {
			log.Fatal(err)
		}

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Fatal(err)
		}

		payload := sdk.Passthru{}
		payload["raw"] = json.RawMessage(byteValue)

		response, err := client.Generator().Generate(generatorType, payload, sdkClient.Namespace, sdkClient.BaseUrl)
		if err != nil {
			log.Fatal(err)
		}

		if response.Chunks != nil {
			for file, code := range response.Chunks {
				err := os.WriteFile(outputDir+"/"+file, []byte(code), 0644)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else if response.Output != "" {
			err := os.WriteFile(outputDir+"/output", []byte(response.Output), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("Generation successful!")
		os.Exit(0)
	},
}
