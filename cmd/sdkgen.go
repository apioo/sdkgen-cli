package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/apioo/sdkgen-cli/sdk"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func Resolve(name string, version string) json.RawMessage {
	match, _ := regexp.MatchString("^([a-zA-Z0-9-_.]+)/([A-Za-z0-9_]+)$", name)
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
			log.Fatal("TypeHub: Could not export specification")
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		export := ExportResponse{}
		err = json.Unmarshal(resBody, &export)
		if err != nil {
			log.Fatal(err)
		}

		return Resolve(export.Href, version)
	} else if strings.HasPrefix(name, "https://") || strings.HasPrefix(name, "http://") {
		res, err := http.Get(name)
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode != 200 {
			log.Fatal("Url " + name + " returned a non-successful response code")
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

func Generate(client *sdk.Client, generatorType string, schema json.RawMessage, outputDir string, namespace string, baseUrl string, mapping map[string]string, remove bool) {
	stat, err := os.Stat(outputDir)
	if err != nil {
		log.Fatal("Provided output directory does not exist")
	}

	if !stat.IsDir() {
		log.Fatal("Provided output directory is not a directory")
	}

	payload := sdk.Passthru{}
	payload["raw"] = schema

	var config = make(map[string]any)
	config["namespace"] = namespace
	config["mapping"] = mapping
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}

	var base64Config = base64.StdEncoding.EncodeToString(jsonConfig)

	response := tryToGenerate(client, generatorType, payload, namespace, base64Config, baseUrl, 0)

	if response.Chunks != nil {
		if remove {
			deleteAllFilesInFolder(outputDir)
		}

		for file, code := range response.Chunks {
			err := os.WriteFile(filepath.Join(outputDir, file), []byte(code), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if response.Output != "" {
		err := os.WriteFile(filepath.Join(outputDir, "output"), []byte(response.Output), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func tryToGenerate(client *sdk.Client, generatorType string, payload sdk.Passthru, namespace string, base64Config string, baseUrl string, retryCount int) sdk.GeneratorResponse {
	response, err := client.Generate(generatorType, payload, namespace, base64Config, baseUrl)
	if err != nil {
		if retryCount > 3 {
			log.Fatal(err)
		}

		time.Sleep(2 * time.Second)

		return tryToGenerate(client, generatorType, payload, namespace, base64Config, baseUrl, retryCount+1)
	}

	return response
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

func deleteAllFilesInFolder(path string) {
	dir, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	files, err := dir.ReadDir(0)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		err = os.Remove(filepath.Join(path, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
	}
}

type Schema struct {
	Type       string             `json:"type"`
	Repository string             `json:"repository"`
	Require    map[string]Package `json:"require"`
}

type Package struct {
	Version   string            `json:"version"`
	Target    string            `json:"target"`
	Namespace string            `json:"namespace"`
	BaseUrl   string            `json:"baseUrl"`
	Mapping   map[string]string `json:"mapping"`
}

type ExportRequest struct {
	Format  string `json:"format"`
	Version string `json:"version"`
}

type ExportResponse struct {
	Href string `json:"href"`
}
