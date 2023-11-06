package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/apioo/sdkgen-cli/sdk"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func Resolve(name string, version string) []byte {
	match, _ := regexp.MatchString("^([A-Za-z0-9]+)/([A-Za-z0-9]+)$", name)
	if match {
		var payload = ExportRequest{
			Format:  "spec-typeapi",
			Version: version,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			log.Fatal(err)
		}

		res, err := http.Post("https://api.typehub.cloud/document/"+name+"/export", "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode != 200 {
			log.Fatal("Could not export specification")
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		export := ExportResponse{}
		json.Unmarshal(resBody, &export)

		return Resolve(export.Href, version)
	} else if strings.HasPrefix(name, "https://") || strings.HasPrefix(name, "http://") {
		res, err := http.Get(name)
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode != 200 {
			log.Fatal("")
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		return resBody
	} else {
		jsonFile, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Fatal(err)
		}

		return byteValue
	}
}

func Generate(client *sdk.Client, generatorType string, schema []byte, outputDir string, namespace string, baseUrl string) {
	stat, err := os.Stat(outputDir)
	if err != nil {
		log.Fatal("Provided output directory does not exist")
	}

	if !stat.IsDir() {
		log.Fatal("Provided output directory does not exist")
	}

	payload := sdk.Passthru{}
	payload["raw"] = json.RawMessage(schema)

	response, err := client.Generate(generatorType, payload, namespace, baseUrl)
	if err != nil {
		log.Fatal(err)
	}

	if response.Chunks != nil {
		for file, code := range response.Chunks {
			err := os.WriteFile(outputDir+"/"+file, []byte(code), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if response.Output != "" {
		err := os.WriteFile(outputDir+"/output", []byte(response.Output), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

type Schema struct {
	Type       string             `json:"type"`
	Repository string             `json:"repository"`
	Require    map[string]Package `json:"require"`
}

type Package struct {
	Version   string `json:"version"`
	Target    string `json:"target"`
	Namespace string `json:"namespace"`
	BaseUrl   string `json:"baseUrl"`
}

type ExportRequest struct {
	Format  string `json:"format"`
	Version string `json:"version"`
}

type ExportResponse struct {
	Href string `json:"href"`
}
