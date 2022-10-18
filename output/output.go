package output

import (
	"bufio"
	"fmt"
	"github.com/briandowns/spinner"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var red = color.New(color.FgRed).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()

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
	fmt.Print(OkString(content))
}

func OkString(content string) string {
	prefix := green("✓")
	return fmt.Sprintf("\t%s %s\n", prefix, content)
}

func Warning(content string) {
	fmt.Print(WarningString(content))
}

func WarningString(content string) string {
	prefix := yellow("✖")
	return fmt.Sprintf("\t%s %s\n", prefix, content)
}

func Error(content string) {
	fmt.Print(ErrorString(content))
}

func ErrorString(content string) string {
	prefix := red("✖")
	return fmt.Sprintf("\t%s %s\n", prefix, content)
}

func Plain(content string) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		fmt.Print(PlainString(scanner.Text()))
	}
}

func PlainString(content string) string {
	return fmt.Sprintf("\t%s\n", content)
}

func Spin(prefix, completeMsg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	s.Prefix = "       " + prefix + " "
	s.FinalMSG = "       " + completeMsg + "\n\n"
	s.Start()

	return s
}

// StateSpin presents and progress spinner with a possible failure state and message
func StateSpin(prefix, completeMsg, failMsg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	s.Prefix = "       " + prefix + " "
	s.FinalMSG = "       " + completeMsg + "\n\n"
	s.Start()

	return s
}

func readInput(input chan rune) {
	var reader = bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	input <- char
}

func InputCancelFunc(waitFunc func(chan<- bool), timeout time.Duration, status chan<- bool) bool {
	input := make(chan rune, 1)
	go readInput(input)

	funcSuccess := make(chan bool, 1)
	go waitFunc(funcSuccess)

	for {
		select {
		case finishState := <-funcSuccess:
			status <- true
			return finishState

		case <-input:
			status <- false
			return false

		case <-time.After(timeout):
			status <- false
			return false
		}
	}
}
