package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install [flags]",
	Short: "Install all dependencies defined in the sdkgen.json file",
	Long:  `Install all dependencies defined in the sdkgen.json file, more information about the schema at https://sdkgen.app/schema`,
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

		var lock = make(map[string]interface{})
		sdkgenLock, err := readFile("./sdkgen.lock")
		if err == nil {
			err = json.Unmarshal(sdkgenLock, &lock)
			if err != nil {
				log.Fatal(err)
			}
		}

		var resolved = make(map[string]interface{})
		var hasChanged = false

		for name, pkg := range schema.Require {
			var spec []byte

			value, ok := lock[name]
			if ok {
				spec, err = json.Marshal(value)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				spec = Resolve(name, pkg.Version)
				hasChanged = true
			}

			Generate(client, schema.Type, spec, pkg.Target, pkg.Namespace, pkg.BaseUrl)

			resolved[name] = json.RawMessage(spec)

			fmt.Println("* Generated " + name)
		}

		if hasChanged {
			lockContent, err := json.Marshal(resolved)
			if err != nil {
				log.Fatal(err)
			}

			err = os.WriteFile("./sdkgen.lock", lockContent, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("Install successful!")
		os.Exit(0)
	},
}
