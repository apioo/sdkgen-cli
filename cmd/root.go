package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "sdkgen",
	Short: "SDKgen is a powerful code generator to automatically build client SDKs for your REST API.",
	Long:  `SDKgen is a powerful code generator to automatically build client SDKs for your REST API. Complete documentation is available at https://sdkgen.app/`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var sdkClient = SdkClient{
	ClientId:     "",
	ClientSecret: "",
	Namespace:    "",
	BaseUrl:      "",
	Remove:       false,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&sdkClient.ClientId, "client-id", "", "The client id is either your username or an app key which you can create at our sdkgen.app backend.")
	rootCmd.PersistentFlags().StringVar(&sdkClient.ClientSecret, "client-secret", "", "This client secret is either your password or an app secret which you can create at our sdkgen.app backend.")

	generateCmd.PersistentFlags().StringVar(&sdkClient.Namespace, "namespace", "", "Optional a namespace for the generated code.")
	generateCmd.PersistentFlags().StringVar(&sdkClient.BaseUrl, "base-url", "", "Optional a base url for your SDK.")

	generateCmd.PersistentFlags().BoolVar(&sdkClient.Remove, "remove", false, "Whether to remove all existing files at the target directory before generation.")
	installCmd.PersistentFlags().BoolVar(&sdkClient.Remove, "remove", false, "Whether to remove all existing files at the target directory before generation.")
	updateCmd.PersistentFlags().BoolVar(&sdkClient.Remove, "remove", false, "Whether to remove all existing files at the target directory before generation.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
