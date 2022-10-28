package output

import (
	"bufio"
	"container/list"
	"fmt"
	"golang.org/x/term"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/pborman/ansi"
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

func FifoOutput(title string, lines int, content <-chan string, closeSignal <-chan bool, finished chan<- bool) {

	// Setup prefix and wrappers
	prefix := "┃ "
	topWrapper := "┏" + strings.Repeat("━", len(title)-1)
	bottomWrapper := "┗" + strings.Repeat("━", len(title)-1)

	// Allocate queue list
	writeQueue := list.New()

	// Calc current terminal width
	// TODO: Test this in headless execution
	terminalWidth, _, err := term.GetSize(0)
	if err != nil {
		return
	}

	// Determine max length of output strings
	maxLineLength := float64(terminalWidth - 10)

	// Write upper wrapper
	Plain("\n" + title + "\n")
	Plain(topWrapper + "\n")

	// Prefill newlines
	for i := 0; i < lines-1; i++ {
		Plain(prefix + "\n")
		writeQueue.PushBack(" ")
	}
	Plain(bottomWrapper + "\n")

	// Allocate previous line write count
	lastWrite := lines

	for {
		select {
		case newContent := <-content:

			// Strip ANSI control codes
			cleanContent, _ := ansi.Strip([]byte(newContent))

			// Split lines
			newContentLines := strings.Split(strings.ReplaceAll(string(cleanContent), "\r\n", "\n"), "\n")

			for _, line := range newContentLines {

				// FIFO the queue
				writeQueue.PushBack(line)
				if writeQueue.Len() > lines {
					writeQueue.Remove(writeQueue.Front())
				}

				// Clear previous writes
				for i := 0; i < lastWrite; i++ {

					// Render ANSI codes to move cursor up one and clear line
					fmt.Print("\033[1A\033[K")
				}

				// Set lastWrite to number of lines we are about to write
				lastWrite = writeQueue.Len()
				for e := writeQueue.Front(); e != nil; e = e.Next() {

					// Check string assertion on content
					lineContent, lineContentOk := e.Value.(string)
					if !lineContentOk {

						// Deduct one line from written count and continue loop
						lastWrite--
						continue
					}

					// Replace tabs with spaces
					lineContent = strings.ReplaceAll(lineContent, "\t", strings.Repeat(" ", 4))

					// Trim to prevent wrapping
					lineContent = lineContent[:int(math.Min(maxLineLength, float64(len(lineContent))))]

					// Render with plain text formatting
					Plain(prefix + lineContent + "\n")
				}

				lastWrite++
				Plain(bottomWrapper + "\n")
			}

		case <-closeSignal:

			// Clear previous writes along with wrappers
			for i := 0; i < lastWrite+3; i++ {

				// Render ANSI codes to move cursor up one and clear line
				fmt.Print("\033[1A\033[K")
			}

			// Give terminal time to catch up. This shouldn't be necessary but mac is gunna mac.
			time.Sleep(time.Millisecond * 100)

			// Write to the finished channel to unblock the parent/calling process and return
			finished <- true
			return
		}
	}
}

func Spin(prefix, completeMsg string) *spinner.Spinner {
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
