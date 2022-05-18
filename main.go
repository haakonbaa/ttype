package main

import (
    "fmt"
    "os/exec"
)

func main() {
    fmt.Println("Hello World")
    highlightFile("test.go","cpp")
}

func highlightFile(filename,language string) {
    out, err := exec.Command(fmt.Sprintf("source-highlight -s '%s' -f esc256 -i '%s' -o test.txt",language,filename)).Output()
    out2 , _ := exec.Command("ls").Output()
    fmt.Println(string(out2))
    fmt.Println(out,err)
}
