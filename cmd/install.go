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
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install [flags]",
	Short: "Install all dependencies defined in the sdkgen.json file",
	Long:  `Install all dependencies defined in the sdkgen.json file, more information about the schema at https://sdkgen.app/schema`,
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

		var lock = make(map[string]interface{})
		var sdkgenLockPath = filepath.Join(cwd, "sdkgen.lock")
		sdkgenLock, err := readFile(sdkgenLockPath)
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

			var outputDir = filepath.Join(cwd, pkg.Target)

			Generate(client, schema.Type, spec, outputDir, pkg.Namespace, pkg.BaseUrl, sdkClient.Remove)

			resolved[name] = json.RawMessage(spec)

			fmt.Println("* Generated " + name)
		}

		if hasChanged {
			lockContent, err := json.Marshal(resolved)
			if err != nil {
				log.Fatal(err)
			}

			err = os.WriteFile(sdkgenLockPath, lockContent, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("Install successful!")
		os.Exit(0)
	},
}
