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

func play() {
	// buffer output
	writer := bufio.NewWriter(os.Stdout)
	writer.Write([]byte("Hello"))
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
	for true {
		char, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		}
		if char == rune('\x7F') {
			writer.WriteString(fmt.Sprintf("[<-]"))
		} else {
			writer.WriteString(fmt.Sprintf("[%c]", char))
		}
		writer.Flush()
	}
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

func generateWords(n int, filename string) []string {
	reWord := regexp.MustCompile(`[a-zA-Z]+`)
	strings := make([]string, n)
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	allWordsInFile := reWord.FindAll(data, -1)
	fmt.Println(len(allWordsInFile))
	for i := 0; i < n; i++ {
		strings[i] = string(allWordsInFile[rand.Intn(len(allWordsInFile)-1)])
	}
	return strings
}
