package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/poloniex/polo-local-dev/cmd"
	_ "github.com/poloniex/polo-local-dev/config"
)

func main() {
	fmt.Println("")
	cmd.Execute()
	fmt.Println("")
}
