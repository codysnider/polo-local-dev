package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/poloniex/polo-local-dev/cmd"
)

func main() {
	cmd.Execute()
}
