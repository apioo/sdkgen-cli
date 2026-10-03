package python

import (
	"bytes"
	_ "embed"
	"fmt"
	"os/exec"
)

//go:embed extract.py
var pythonExtractorScript string

// GenerateFromFastAPI executes python using the embedded script.
// pythonPath can be "python3", "python", or a path to a virtualenv python binary.
// entrypoint is the target module and app variable, e.g. "main:app".
func GenerateFromFastAPI(pythonPath string, entrypoint string) ([]byte, error) {
	if pythonPath == "" {
		pythonPath = "python3"
	}

	// Run `python -c "<embedded script>" main:app`
	cmd := exec.Command(pythonPath, "-c", pythonExtractorScript, entrypoint)

	var stdoutBytes bytes.Buffer
	var stderrBytes bytes.Buffer
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to extract TypeAPI spec from FastAPI: %w\nStderr: %s", err, stderrBytes.String())
	}

	return stdoutBytes.Bytes(), nil
}
