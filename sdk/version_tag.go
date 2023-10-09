// VersionTag automatically sdk by SDKgen please do not edit this file manually
// @see https://sdkgen.app

package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/apioo/sdkgen-go"
	"io"
	"net/http"
	"net/url"
)

type VersionTag struct {
	internal *sdkgen.TagAbstract
}

// Create Creates a new version
func (client *VersionTag) Create(projectId string, payload VersionCreate) (Message, error) {
	pathParams := make(map[string]interface{})
	pathParams["project_id"] = projectId

	queryParams := make(map[string]interface{})

	u, err := url.Parse(client.internal.Parser.Url("/project/:project_id/version", pathParams))
	if err != nil {
		return Message{}, errors.New("could not parse url")
	}

	u.RawQuery = client.internal.Parser.Query(queryParams).Encode()

	raw, err := json.Marshal(payload)
	if err != nil {
		return Message{}, errors.New("could not marshal provided JSON data")
	}

	var reqBody = bytes.NewReader(raw)

	req, err := http.NewRequest("POST", u.String(), reqBody)
	if err != nil {
		return Message{}, errors.New("could not create request")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.internal.HttpClient.Do(req)
	if err != nil {
		return Message{}, errors.New("could not send request")
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Message{}, errors.New("could not read response body")
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response Message
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return Message{}, errors.New("could not unmarshal JSON response")
		}

		return response, nil
	}

	switch resp.StatusCode {
	case 400:
		var response Message
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return Message{}, errors.New("could not unmarshal JSON response")
		}

		return Message{}, &MessageException{
			Payload: response,
		}
	case 404:
		var response Message
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return Message{}, errors.New("could not unmarshal JSON response")
		}

		return Message{}, &MessageException{
			Payload: response,
		}
	case 500:
		var response Message
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return Message{}, errors.New("could not unmarshal JSON response")
		}

		return Message{}, &MessageException{
			Payload: response,
		}
	default:
		return Message{}, errors.New("the server returned an unknown status code")
	}
}

// Get Returns a version
func (client *VersionTag) Get(projectId string, versionId string) (Version, error) {
	pathParams := make(map[string]interface{})
	pathParams["project_id"] = projectId
	pathParams["version_id"] = versionId

	queryParams := make(map[string]interface{})

	u, err := url.Parse(client.internal.Parser.Url("/project/:project_id/version/:version_id", pathParams))
	if err != nil {
		return Version{}, errors.New("could not parse url")
	}

	u.RawQuery = client.internal.Parser.Query(queryParams).Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Version{}, errors.New("could not create request")
	}

	resp, err := client.internal.HttpClient.Do(req)
	if err != nil {
		return Version{}, errors.New("could not send request")
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Version{}, errors.New("could not read response body")
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response Version
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return Version{}, errors.New("could not unmarshal JSON response")
		}

		return response, nil
	}

	switch resp.StatusCode {
	case 400:
		var response Message
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return Version{}, errors.New("could not unmarshal JSON response")
		}

		return Version{}, &MessageException{
			Payload: response,
		}
	case 404:
		var response Message
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return Version{}, errors.New("could not unmarshal JSON response")
		}

		return Version{}, &MessageException{
			Payload: response,
		}
	case 500:
		var response Message
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return Version{}, errors.New("could not unmarshal JSON response")
		}

		return Version{}, &MessageException{
			Payload: response,
		}
	default:
		return Version{}, errors.New("the server returned an unknown status code")
	}
}

// GetAll Returns all versions
func (client *VersionTag) GetAll(projectId string, startIndex int, count int, search string) (VersionCollection, error) {
	pathParams := make(map[string]interface{})
	pathParams["project_id"] = projectId

	queryParams := make(map[string]interface{})
	queryParams["startIndex"] = startIndex
	queryParams["count"] = count
	queryParams["search"] = search

	u, err := url.Parse(client.internal.Parser.Url("/project/:project_id/version", pathParams))
	if err != nil {
		return VersionCollection{}, errors.New("could not parse url")
	}

	u.RawQuery = client.internal.Parser.Query(queryParams).Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return VersionCollection{}, errors.New("could not create request")
	}

	resp, err := client.internal.HttpClient.Do(req)
	if err != nil {
		return VersionCollection{}, errors.New("could not send request")
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return VersionCollection{}, errors.New("could not read response body")
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response VersionCollection
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return VersionCollection{}, errors.New("could not unmarshal JSON response")
		}

		return response, nil
	}

	switch resp.StatusCode {
	case 400:
		var response Message
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return VersionCollection{}, errors.New("could not unmarshal JSON response")
		}

		return VersionCollection{}, &MessageException{
			Payload: response,
		}
	case 404:
		var response Message
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return VersionCollection{}, errors.New("could not unmarshal JSON response")
		}

		return VersionCollection{}, &MessageException{
			Payload: response,
		}
	case 500:
		var response Message
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return VersionCollection{}, errors.New("could not unmarshal JSON response")
		}

		return VersionCollection{}, &MessageException{
			Payload: response,
		}
	default:
		return VersionCollection{}, errors.New("the server returned an unknown status code")
	}
}

func NewVersionTag(httpClient *http.Client, parser *sdkgen.Parser) *VersionTag {
	return &VersionTag{
		internal: &sdkgen.TagAbstract{
			HttpClient: httpClient,
			Parser:     parser,
		},
	}
}
