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
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update [flags]",
	Short: "Update all dependencies defined in the sdkgen.json file",
	Long:  `Update all dependencies defined in the sdkgen.json file, more information about the schema at https://sdkgen.app/schema`,
	Run: func(cmd *cobra.Command, args []string) {
		var client = sdkClient.GetClient()

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		var sdkgenFilePath = filepath.Join(cwd, "sdkgen.json")
		sdkgenFile, err := readFile(sdkgenFilePath)
		if err != nil {
			log.Fatal(err)
		}

		schema := Schema{}
		err = json.Unmarshal(sdkgenFile, &schema)
		if err != nil {
			log.Fatal(err)
		}

		var resolved = make(map[string]interface{})

		for name, pkg := range schema.Require {
			spec := Resolve(name, pkg.Version)

			var outputDir = filepath.Join(cwd, pkg.Target)

			Generate(client, schema.Type, spec, outputDir, pkg.Namespace, pkg.BaseUrl, sdkClient.Remove)

			resolved[name] = json.RawMessage(spec)

			fmt.Println("* Generated " + name)
		}

		lockContent, err := json.Marshal(resolved)
		if err != nil {
			log.Fatal(err)
		}

		var sdkgenLockPath = filepath.Join(cwd, "sdkgen.lock")
		err = os.WriteFile(sdkgenLockPath, lockContent, 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Update successful!")
		os.Exit(0)
	},
}
