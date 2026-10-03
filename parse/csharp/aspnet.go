package csharp

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed Extract.cs
var aspnetExtractorSource string

// GenerateFromAspNet executes the dotnet CLI against an ASP.NET Core project.
// dotnetPath can be "dotnet" or a full path to the dotnet executable.
// projectPath is the directory containing the .csproj file or compiled assembly.
func GenerateFromAspNet(dotnetPath string, projectPath string) ([]byte, error) {
	if dotnetPath == "" {
		dotnetPath = "dotnet"
	}

	// Create a temp directory for the C# extractor source
	tmpDir, err := os.MkdirTemp("", "sdkgen_aspnet")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	csFilePath := filepath.Join(tmpDir, "Extract.cs")
	if err := os.WriteFile(csFilePath, []byte(aspnetExtractorSource), 0644); err != nil {
		return nil, fmt.Errorf("failed to write C# extractor file: %w", err)
	}

	// Execute via dotnet script / single file runner
	// Passes projectPath to C# script to locate compiled assemblies
	cmd := exec.Command(dotnetPath, "run", "--project", projectPath, "--", csFilePath)

	var stdoutBytes bytes.Buffer
	var stderrBytes bytes.Buffer
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to extract TypeAPI spec from ASP.NET Core: %w\nStderr: %s", err, stderrBytes.String())
	}

	return stdoutBytes.Bytes(), nil
}
