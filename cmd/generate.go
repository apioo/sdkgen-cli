package cmd

import (
	"encoding/json"
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
	Use:   "generate",
	Short: "Converts a provided ",
	Long:  `All software has versions. This is Hugo's`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		var client = sdkClient.GetClient()

		var generatorType = args[0]
		var file = args[1]
		var outputDir = args[2]
		var namespace = ""
		var baseUrl = ""

		stat, err := os.Stat(outputDir)
		if err != nil {
			log.Fatal(err)
		}

		if !stat.IsDir() {
			log.Fatal("Provided output directory does not exist")
		}

		jsonFile, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Fatal(err)
		}

		payload := map[string]any{}
		err = json.Unmarshal(byteValue, &payload)
		if err != nil {
			log.Fatal(err)
		}

		response, err := client.Generator().Generate(generatorType, payload, namespace, baseUrl)
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
