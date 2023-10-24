package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	rootCmd.AddCommand(typesCmd)
}

var typesCmd = &cobra.Command{
	Use:   "types",
	Short: "Shows all available types",
	Long:  `Shows all available types which can be used for code generation, you can use one of the returned type at the generate command`,
	Run: func(cmd *cobra.Command, args []string) {
		var client = sdkClient.GetClient()

		types, err := client.GetTypes()
		if err != nil {
			log.Fatal(err)
		}

		for _, _type := range types.Types {
			fmt.Println("* " + _type.Name)
		}
	},
}
