package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"syscall"
)

func readTest() {
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	fmt.Println(string(b))
}

type Word struct {
	wformat string // formated string to show on screen
	word    string // string to match with input
	split   rune   // split between this and next word
}

const ANSI_CLEARSCREEN = "\033[1;1H\033[2J"
const ANSI_SHOWCURSOR  = "\033[?25h"
const ANSI_RESET       = "\033[0m"

const ANSI_RED         = "\033[31m"

func play(words *[]Word) {
	// buffer output
	writer := bufio.NewWriter(os.Stdout)
	writer.Flush()
	defer writer.Flush()
	// create reader
	reader := bufio.NewReader(os.Stdin)
	// handle interrupts
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGWINCH, syscall.SIGINT)
	go func() {
		var sig os.Signal
		for true {
			sig = <-sigs
			if sig == syscall.SIGINT {
				fmt.Printf(ANSI_SHOWCURSOR)
				fmt.Printf(ANSI_RESET)
				exec.Command("stty", "-F", "/dev/tty", "echo").Run()
				os.Exit(1)
			}
			if sig == syscall.SIGWINCH {
				fmt.Println("Screen resize")
			}
			fmt.Println("recieved signal")
			fmt.Println(sig)
		}
	}()
	wordIndex := 0
	insertedWords := make([]string, len(*words))
	fmt.Println(insertedWords, wordIndex)
	for wordIndex < len(*words) {
		// begin writing to screen
		writer.WriteString(ANSI_CLEARSCREEN)
		// Nice debug information:
		writer.WriteString(fmt.Sprintf("%v\n", insertedWords))
		for i := 0; i < len(*words); i++ {
			if i <= wordIndex {
				formatWordErrors((*words)[i].wformat, insertedWords[i], writer)
				writer.WriteString(string((*words)[i].split))
			} else {
				writer.WriteString((*words)[i].wformat)
				writer.WriteString(string((*words)[i].split))
			}
		}
		writer.Flush()

		// Handle input logic
		char, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		}
		if char == (*words)[wordIndex].split {
			if insertedWords[wordIndex] == "" {

			} else {
				wordIndex += 1
				if wordIndex >= len(*words) {
					break
				}
			}
		} else if char == rune('\x7f') {
			// got backspace
			// jump one word back if current word is empty
			if len(insertedWords[wordIndex]) == 0 {
				wordIndex = max(0, wordIndex-1)
			}
			// erase last letter of word
			if len(insertedWords[wordIndex]) > 0 {
				insertedWords[wordIndex] = insertedWords[wordIndex][:len(insertedWords[wordIndex])-1]
			}
		} else {
			insertedWords[wordIndex] = insertedWords[wordIndex] + string(char)
		}
	}
}

func max[T int | uint](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Takes a word target and the user input and formats it writing to writer
func formatWordErrors(target, input string, writer *bufio.Writer) {
	if target == input {
		writer.WriteString(target)
	} else {
		// if input is longer than target pad target with spacec to length of input
		for i := len(target); i < len(input); i++ {
			target += string(" ")
		}
		// if input is shorter than target, pad input with correct letters of target
		if len(input) < len(target) {
			input += target[len(input):]
		}
		// now input and target can easily be compared
		// writer.WriteString(fmt.Sprintf("target: [%s]\ninput:  [%s]\n",target,input))
		isCorrect := true
		for i := 0; i < len(target); i++ {
			match := target[i] == input[i]
			// change color in transitions between correct and incorrect letters
			if match && !isCorrect {
				// change back to default color
				writer.WriteString(ANSI_RESET)
				isCorrect = true
			} else if !match && isCorrect {
				// change to red color
				writer.WriteString(ANSI_RED)
				isCorrect = false
			}
			writer.WriteString(string(input[i]))
		}
		if !isCorrect {
			// make sure color is normal at EOT
			writer.WriteString(ANSI_RESET)
		}
	}
	// writer.WriteString(input)
}

func highlightFile(filename, language string) {
	out, err := exec.Command(fmt.Sprintf("source-highlight -s '%s' -f esc256 -i '%s' -o test.txt", language, filename)).Output()
	out2, _ := exec.Command("ls").Output()
	fmt.Println(string(out2))
	fmt.Println(out, err)
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func generateWords(n int, filename string) *[]Word {
	reWord := regexp.MustCompile(`[a-zA-Z]+`)
	words := make([]Word, n)
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	allWordsInFile := reWord.FindAll(data, -1)
	for i := 0; i < n; i++ {
		w := string(allWordsInFile[rand.Intn(len(allWordsInFile)-1)])
		words[i] = Word{
			wformat: w,
			word:    w,
			split:   rune(' '),
		}
	}
	return &words
}
