// MetaTag automatically sdk by SDKgen please do not edit this file manually
// @see https://sdkgen.app

package sdk

import (
	"encoding/json"
	"errors"
	"github.com/apioo/sdkgen-go"
	"io"
	"net/http"
	"net/url"
)

type MetaTag struct {
	internal *sdkgen.TagAbstract
}

// GetAbout
func (client *MetaTag) GetAbout() (SystemAbout, error) {
	pathParams := make(map[string]interface{})

	queryParams := make(map[string]interface{})

	u, err := url.Parse(client.internal.Parser.Url("/", pathParams))
	if err != nil {
		return SystemAbout{}, errors.New("could not parse url")
	}

	u.RawQuery = client.internal.Parser.Query(queryParams).Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return SystemAbout{}, errors.New("could not create request")
	}

	resp, err := client.internal.HttpClient.Do(req)
	if err != nil {
		return SystemAbout{}, errors.New("could not send request")
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return SystemAbout{}, errors.New("could not read response body")
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response SystemAbout
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return SystemAbout{}, errors.New("could not unmarshal JSON response")
		}

		return response, nil
	}

	switch resp.StatusCode {
	default:
		return SystemAbout{}, errors.New("the server returned an unknown status code")
	}
}

func NewMetaTag(httpClient *http.Client, parser *sdkgen.Parser) *MetaTag {
	return &MetaTag{
		internal: &sdkgen.TagAbstract{
			HttpClient: httpClient,
			Parser:     parser,
		},
	}
}
