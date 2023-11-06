package main

import (
	"github.com/apioo/sdkgen-cli/cmd"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cmd.Execute()
}
