package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
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

		sdkgenFile, err := readFile("./sdkgen.json")
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

			Generate(client, schema.Type, spec, pkg.Target, pkg.Namespace, pkg.BaseUrl)

			resolved[name] = json.RawMessage(spec)

			fmt.Println("* Generated " + name)
		}

		lockContent, err := json.Marshal(resolved)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile("./sdkgen.lock", lockContent, 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Update successful!")
		os.Exit(0)
	},
}
