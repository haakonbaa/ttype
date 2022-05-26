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

const CLEARSCREEN = "\033[1;1H\033[2J"

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
				fmt.Printf("\033[?25h")
				fmt.Printf("\033[28m")
				fmt.Printf("\033[0m")
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
		writer.WriteString(CLEARSCREEN)
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

		char, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		}
		if char == (*words)[wordIndex].split {
			wordIndex += 1
			if wordIndex >= len(*words) {
				break
			}
		} else if char == rune('\x7f') {
			// got backspace
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
		for i := 0; i < max(0, len(input)-len(target)); i++ {
			target += string(" ")
		}
		// if input is shorter than target, pad input with correct letters of target
		if len(input) < len(target) {
			input += target[len(input):len(target)]
		}
		// now input and target can easily be compared
		//writer.WriteString(fmt.Sprintf("[%s]\n[%s]\n",target,input))
		isCorrect := true
		for i := 0; i < len(target); i++ {
			match := target[i] == input[i]
			if match && !isCorrect {
				// change back to default color
				writer.WriteString("\033[0m")
				isCorrect = true
			} else if !match && isCorrect {
				writer.WriteString("\033[31m")
				isCorrect = false
			}
			writer.WriteString(string(target[i]))
		}
		if !isCorrect {
			writer.WriteString("\033[0m")
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
