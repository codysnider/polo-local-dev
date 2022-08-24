package output

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var red = color.New(color.FgRed).PrintfFunc()
var green = color.New(color.FgGreen).PrintfFunc()
var yellow = color.New(color.FgYellow).PrintfFunc()

func Title(content string) {
	fmt.Println("┏" + strings.Repeat("━", len(content)+2) + "┓")
	fmt.Printf("┃ %s ┃\n", content)
	fmt.Println("┗" + strings.Repeat("━", len(content)+2) + "┛")
}

func Section(content string) {
	fmt.Printf("\n   %s \n", content)
	fmt.Println("  " + strings.Repeat("━", len(content)+2))
	fmt.Println("")
}

func Ok(content string) {
	green("    ✓ ")
	fmt.Printf(" %s \n", content)
}

func Warning(content string) {
	yellow("    ✖ ")
	fmt.Printf(" %s \n", content)
}

func Error(content string) {
	red("    ✖ ")
	fmt.Printf(" %s \n", content)
}
