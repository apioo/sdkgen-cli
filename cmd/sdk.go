package cmd

import (
	"github.com/apioo/sdkgen-cli/sdk"
	"github.com/apioo/sdkgen-go"
	"log"
)

type SdkClient struct {
	ClientId     string
	ClientSecret string
}

func (sdkClient *SdkClient) GetClient() *sdk.Client {
	credentials := sdkgen.OAuth2{
		ClientId:     sdkClient.ClientId,
		ClientSecret: sdkClient.ClientSecret,
		TokenUrl:     "https://api.sdkgen.app/authorization/token",
	}
	client, err := sdk.NewClient("", credentials)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
