package php

import (
	"bytes"
	_ "embed"
	"fmt"
	"os/exec"
)

//go:embed extract.php
var symfonyExtractorScript string

// GenerateFromSymfony executes PHP using the embedded script.
// phpPath can be "php" or a full path to a local PHP binary.
// projectPath is the root directory of the Symfony application.
func GenerateFromSymfony(phpPath string, projectPath string) ([]byte, error) {
	if phpPath == "" {
		phpPath = "php"
	}

	// Run `php -r "<embedded script>" /path/to/symfony/project`
	cmd := exec.Command(phpPath, "-r", symfonyExtractorScript, projectPath)

	var stdoutBytes bytes.Buffer
	var stderrBytes bytes.Buffer
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to extract TypeAPI spec from Symfony: %w\nStderr: %s", err, stderrBytes.String())
	}

	return stdoutBytes.Bytes(), nil
}
