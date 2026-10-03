package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/spf13/cobra"
)

var (
	patchTarget string
	patchData   string
	patchFile   string
	patchDryRun bool
)

func init() {
	patchCmd.Flags().StringVarP(&patchTarget, "target", "t", "", "Target entity (e.g. 'operations.getUsers' or 'types.User') [REQUIRED]")
	patchCmd.Flags().StringVarP(&patchData, "data", "d", "", "JSON Patch operations string or replacement payload")
	patchCmd.Flags().StringVarP(&patchFile, "file", "f", "", "File containing JSON Patch operations or replacement payload")
	patchCmd.Flags().BoolVar(&patchDryRun, "dry-run", false, "Preview patch result without writing changes to disk")

	_ = patchCmd.MarkFlagRequired("target")
	rootCmd.AddCommand(patchCmd)
}

var patchCmd = &cobra.Command{
	Use:   "patch SCHEMA_FILE",
	Short: "Applies RFC 6902 JSON Patches strictly within a targeted TypeAPI entity boundary",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		schemaFile := args[0]

		schemaFilePath, err := filepath.Abs(schemaFile)
		if err != nil {
			log.Fatal(err)
		}

		fileBytes, err := os.ReadFile(schemaFilePath)
		if err != nil {
			log.Fatalf("Error reading schema file %s: %v\n", schemaFile, err)
		}

		// 1. Resolve Target Base JSON Pointer
		targetPointer, err := buildTargetPointer(patchTarget)
		if err != nil {
			log.Fatalf("Invalid --target: %v\n", err)
		}

		// Read raw input payload
		payloadStr := patchData
		if patchFile != "" {
			b, err := os.ReadFile(patchFile)
			if err != nil {
				log.Fatalf("Error reading patch file: %v\n", err)
			}
			payloadStr = string(b)
		}

		if strings.TrimSpace(payloadStr) == "" {
			log.Fatal("Either --data (-d) or --file (-f) must be provided")
		}

		// 2. Normalize and Scope JSON Patch Operations to Target
		var finalPatch []map[string]interface{}

		if isJSONArray([]byte(payloadStr)) {
			var rawOps []map[string]interface{}
			if err := json.Unmarshal([]byte(payloadStr), &rawOps); err != nil {
				log.Fatalf("Invalid JSON patch array: %v\n", err)
			}

			for i, op := range rawOps {
				path, ok := op["path"].(string)
				if !ok {
					log.Fatalf("Patch operation at index %d missing 'path'", i)
				}

				// Resolve path relative to target base pointer
				scopedPath := qualifyPointer(targetPointer, path)
				op["path"] = scopedPath

				// Validate 'from' path for 'move' / 'copy' operations
				if fromPath, ok := op["from"].(string); ok {
					op["from"] = qualifyPointer(targetPointer, fromPath)
				}

				finalPatch = append(finalPatch, op)
			}
		} else {
			// Single replacement payload: Replace the entire target entity
			finalPatch = []map[string]interface{}{
				{
					"op":    "add",
					"path":  targetPointer,
					"value": json.RawMessage(payloadStr),
				},
			}
		}

		// 3. Compile and Apply RFC 6902 Patch
		patchBytes, err := json.Marshal(finalPatch)
		if err != nil {
			log.Fatalf("Failed to marshal patch operations: %v\n", err)
		}

		patchObj, err := jsonpatch.DecodePatch(patchBytes)
		if err != nil {
			log.Fatalf("Failed to parse JSON patch: %v\n", err)
		}

		patchedBytes, err := patchObj.Apply(fileBytes)
		if err != nil {
			log.Fatalf("Patch application failed on target '%s': %v\n", patchTarget, err)
		}

		// 4. Format Output and Validate
		var doc interface{}
		if err := json.Unmarshal(patchedBytes, &doc); err != nil {
			log.Fatalf("Patched result is invalid JSON: %v\n", err)
		}

		formattedBytes, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			log.Fatalf("Failed to format JSON output: %v\n", err)
		}

		if patchDryRun {
			fmt.Printf("🔍 [DRY-RUN] Successfully applied patch to target '%s':\n", patchTarget)
			fmt.Println(string(formattedBytes))
			return
		}

		if err := os.WriteFile(schemaFilePath, formattedBytes, 0644); err != nil {
			log.Fatalf("Failed to write patched file: %v\n", err)
		}

		fmt.Printf("✅ Successfully patched target '%s' in %s\n", patchTarget, schemaFile)
	},
}

// Builds the root JSON pointer for target selector (e.g. 'types.User' -> '/definitions/User')
func buildTargetPointer(target string) (string, error) {
	parts := strings.Split(target, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("format must be 'operations.<name>' or 'types.<name>'")
	}

	domain := parts[0]
	entityName := strings.Join(parts[1:], ".")

	if domain == "types" {
		domain = "definitions"
	}

	if domain != "operations" && domain != "definitions" {
		return "", fmt.Errorf("domain must be 'operations' or 'types'")
	}

	return fmt.Sprintf("/%s/%s", domain, entityName), nil
}

// Qualifies relative or absolute patch paths to guarantee they remain inside target boundary
func qualifyPointer(targetPointer string, path string) string {
	if path == "" || path == "/" {
		return targetPointer
	}

	// Prevent agent from escaping target boundary using root pointer
	if strings.HasPrefix(path, targetPointer) {
		return path
	}

	// Normalize relative paths (e.g., '/properties/age' or 'properties/age')
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return targetPointer + path
}

func isJSONArray(b []byte) bool {
	var a []interface{}
	return json.Unmarshal(b, &a) == nil
}
