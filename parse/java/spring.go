package java

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed Extract.java
var springExtractorSource string

// GenerateFromSpring executes Java against the compiled output directory of a Spring Boot project.
// classpath typically points to "target/classes:target/dependency/*" (Maven) or "build/classes/java/main" (Gradle).
func GenerateFromSpring(javaPath string, projectPath string, classPath string) ([]byte, error) {
	if javaPath == "" {
		javaPath = "java"
	}

	// Create a temp file for the Java source code
	tmpDir, err := os.MkdirTemp("", "sdkgen_spring")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	javaFilePath := filepath.Join(tmpDir, "Extract.java")
	if err := os.WriteFile(javaFilePath, []byte(springExtractorSource), 0644); err != nil {
		return nil, fmt.Errorf("failed to write java extractor file: %w", err)
	}

	// Fallback classpath to Maven target directory if not specified
	if classPath == "" {
		classPath = fmt.Sprintf("%s/target/classes:%s/target/dependency/*", projectPath, projectPath)
	}

	// Execute Java 11+ single-file source runner:
	// java -cp <classpath> ExtractSpring.java <projectRootPackageOrClassesDir>
	cmd := exec.Command(javaPath, "-cp", classPath, javaFilePath, projectPath)

	var stdoutBytes bytes.Buffer
	var stderrBytes bytes.Buffer
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to extract TypeAPI spec from Spring: %w\nStderr: %s", err, stderrBytes.String())
	}

	return stdoutBytes.Bytes(), nil
}
