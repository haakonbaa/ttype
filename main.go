package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// hide cursor
	fmt.Printf("\033[?25l")
	// hide input
	//f mt.Printf("\033[8m")
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// highlightFile("test.go","cpp")
	words := generateWords(5, "./words/english1000")
	for i, v := range words {
		fmt.Println(i, v)
	}
	play()
	for true {

	}
	readTest()
}
