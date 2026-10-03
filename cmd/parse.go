package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/apioo/sdkgen-cli/parse/csharp"
	"github.com/apioo/sdkgen-cli/parse/java"
	"github.com/apioo/sdkgen-cli/parse/javascript"
	"github.com/apioo/sdkgen-cli/parse/php"
	"github.com/apioo/sdkgen-cli/parse/python"
	"github.com/spf13/cobra"
)

var (
	pythonPath string
	phpPath    string
	nodePath   string
	javaPath   string
	dotnetPath string
	classPath  string
	outputFile string
)

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.Flags().StringVarP(&pythonPath, "python", "p", "python3", "Path to the Python binary or virtual environment")
	parseCmd.Flags().StringVar(&phpPath, "php", "php", "Path to the PHP binary")
	parseCmd.Flags().StringVar(&nodePath, "node", "node", "Path to the Node.js binary")
	parseCmd.Flags().StringVar(&javaPath, "java", "java", "Path to the Java binary")
	parseCmd.Flags().StringVar(&dotnetPath, "dotnet", "dotnet", "Path to the .NET CLI binary")
	parseCmd.Flags().StringVarP(&classPath, "classpath", "c", "", "Classpath for compiled Spring Boot classes/dependencies")
	parseCmd.Flags().StringVarP(&outputFile, "out", "o", "", "Output file path for generated TypeAPI spec (default prints to stdout)")
}

var parseCmd = &cobra.Command{
	Use:   "parse FRAMEWORK [ENTRYPOINT]",
	Short: "Parses an existing framework project and generates a TypeAPI specification",
	Long: `Scans a web framework project and extracts routes, operations, and data models to generate a TypeAPI specification.

Supported frameworks:
  - aspnet   (entrypoint format: project directory or .csproj path, defaults to ".")
  - fastapi  (entrypoint format: "module:app", e.g. "main:app")
  - nestjs   (entrypoint format: project directory or main entrypoint, defaults to ".")
  - spring   (entrypoint format: project directory path, defaults to ".")
  - symfony  (entrypoint format: project directory path, defaults to ".")

Examples:
  sdkgen parse aspnet . -o typeapi.json
  sdkgen parse fastapi main:app -o typeapi.json
  sdkgen parse nestjs . -o typeapi.json
  sdkgen parse spring . -o typeapi.json
  sdkgen parse symfony . -o typeapi.json`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		framework := args[0]

		var specBytes []byte
		var err error

		switch framework {
		case "aspnet":
			projectPath := "."
			if len(args) >= 2 {
				projectPath = args[1]
			}
			fmt.Printf("🔍 Extracting TypeAPI specification from ASP.NET Core project (%s)...\n", projectPath)
			specBytes, err = csharp.GenerateFromAspNet(dotnetPath, projectPath)
			if err != nil {
				log.Fatalf("❌ Parse failed: %v\n", err)
			}

		case "fastapi":
			if len(args) < 2 {
				log.Fatalf("❌ Missing entrypoint for FastAPI (e.g. 'sdkgen parse fastapi main:app')")
			}
			entrypoint := args[1]
			fmt.Printf("🔍 Extracting TypeAPI specification from FastAPI (%s)...\n", entrypoint)
			specBytes, err = python.GenerateFromFastAPI(pythonPath, entrypoint)
			if err != nil {
				log.Fatalf("❌ Parse failed: %v\n", err)
			}

		case "nestjs":
			projectPath := "."
			if len(args) >= 2 {
				projectPath = args[1]
			}
			fmt.Printf("🔍 Extracting TypeAPI specification from NestJS project (%s)...\n", projectPath)
			specBytes, err = javascript.GenerateFromNestJS(nodePath, projectPath)
			if err != nil {
				log.Fatalf("❌ Parse failed: %v\n", err)
			}

		case "spring":
			projectPath := "."
			if len(args) >= 2 {
				projectPath = args[1]
			}
			fmt.Printf("🔍 Extracting TypeAPI specification from Spring project (%s)...\n", projectPath)
			specBytes, err = java.GenerateFromSpring(javaPath, projectPath, classPath)
			if err != nil {
				log.Fatalf("❌ Parse failed: %v\n", err)
			}

		case "symfony":
			projectPath := "."
			if len(args) >= 2 {
				projectPath = args[1]
			}
			fmt.Printf("🔍 Extracting TypeAPI specification from Symfony project (%s)...\n", projectPath)
			specBytes, err = php.GenerateFromSymfony(phpPath, projectPath)
			if err != nil {
				log.Fatalf("❌ Parse failed: %v\n", err)
			}

		default:
			log.Fatalf("❌ Unsupported framework '%s'. Supported frameworks: aspnet, fastapi, nestjs, spring, symfony\n", framework)
		}

		// Write output to file if specified, otherwise print to stdout
		if outputFile != "" {
			if err := os.WriteFile(outputFile, specBytes, 0644); err != nil {
				log.Fatalf("❌ Failed to write output file %s: %v\n", outputFile, err)
			}
			fmt.Printf("✅ Successfully generated TypeAPI specification at %s\n", outputFile)
		} else {
			fmt.Println(string(specBytes))
		}
	},
}
