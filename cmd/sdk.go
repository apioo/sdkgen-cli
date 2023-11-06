package cmd

import (
	"github.com/apioo/sdkgen-cli/sdk"
	"github.com/apioo/sdkgen-go"
	"log"
	"os"
)

type SdkClient struct {
	ClientId     string
	ClientSecret string
	Namespace    string
	BaseUrl      string
}

func (sdkClient *SdkClient) GetClient() *sdk.Client {
	var tokenStore = sdkgen.MemoryTokenStore{}

	if sdkClient.ClientId == "" {
		sdkClient.ClientId = os.Getenv("SDKGEN_CLIENT_ID")
	}

	if sdkClient.ClientSecret == "" {
		sdkClient.ClientSecret = os.Getenv("SDKGEN_CLIENT_SECRET")
	}

	credentials := sdkgen.OAuth2{
		ClientId:     sdkClient.ClientId,
		ClientSecret: sdkClient.ClientSecret,
		TokenUrl:     "https://api.sdkgen.app/authorization/token",
		TokenStore:   tokenStore,
	}

	client, err := sdk.NewClient("https://api.sdkgen.app", credentials)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
