package cmd

import (
	"log"
	"os"

	"github.com/apioo/sdkgen-cli/sdk"
	"github.com/apioo/sdkgen-go/v2"
)

type SdkClient struct {
	ClientId     string
	ClientSecret string
	Namespace    string
	BaseUrl      string
	Remove       bool
}

func (sdkClient *SdkClient) GetClient() *sdk.Client {
	var tokenStore = sdkgen.MemoryTokenStore{}

	if sdkClient.ClientId == "" {
		sdkClient.ClientId = os.Getenv("SDKGEN_CLIENT_ID")
	}

	if sdkClient.ClientSecret == "" {
		sdkClient.ClientSecret = os.Getenv("SDKGEN_CLIENT_SECRET")
	}

	var credentials sdkgen.CredentialsInterface
	if sdkClient.ClientId != "" && sdkClient.ClientSecret != "" {
		credentials = sdkgen.OAuth2{
			ClientId:     sdkClient.ClientId,
			ClientSecret: sdkClient.ClientSecret,
			TokenUrl:     "https://api.sdkgen.app/authorization/token",
			TokenStore:   tokenStore,
		}
	} else {
		credentials = sdkgen.Anonymous{}
	}

	client, err := sdk.NewClient("https://api.sdkgen.app", credentials)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
