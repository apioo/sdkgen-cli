package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "sdkgen",
	Short: "SDKgen is a powerful code generator to automatically build client SDKs for your REST API.",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at https://sdkgen.app/`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var sdkClient = SdkClient{
	ClientId:     "",
	ClientSecret: "",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&sdkClient.ClientId, "client-id", "", "The client id is either your username or an app key which you can create at our sdkgen.app backend.")
	rootCmd.PersistentFlags().StringVar(&sdkClient.ClientSecret, "client-secret", "", "This client secret is either your password or an app secret which you can create at our sdkgen.app backend.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
