package javascript

import (
	"bytes"
	_ "embed"
	"fmt"
	"os/exec"
)

//go:embed extract.js
var nestjsExtractorScript string

// GenerateFromNestJS executes Node.js using the embedded script.
// nodePath can be "node" or a path to a specific Node binary.
// projectPath is the root directory of the NestJS application (or main module entrypoint).
func GenerateFromNestJS(nodePath string, projectPath string) ([]byte, error) {
	if nodePath == "" {
		nodePath = "node"
	}

	// Execute `node -e "<embedded script>" /path/to/nestjs/project`
	cmd := exec.Command(nodePath, "-e", nestjsExtractorScript, projectPath)

	var stdoutBytes bytes.Buffer
	var stderrBytes bytes.Buffer
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to extract TypeAPI spec from NestJS: %w\nStderr: %s", err, stderrBytes.String())
	}

	return stdoutBytes.Bytes(), nil
}
