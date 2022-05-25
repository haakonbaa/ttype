package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	fmt.Printf("\033[?25l") // hide cursor
	//f mt.Printf("\033[8m") // hide input
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// highlightFile("test.go","cpp")
	words := generateWords(5, "./words/english1000")
	for i, v := range *words {
		fmt.Println(i, v)
	}
	play(words)
    fmt.Printf("\033[?25h")
	fmt.Printf("\033[28m")
    fmt.Printf("\033[0m") // reset
    exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}
