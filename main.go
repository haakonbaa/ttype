package main

import (
    "fmt"
    "os/exec"
    "os"
    "regexp"
    "math/rand"
)

func main() {
    fmt.Println("Hello World")
    // highlightFile("test.go","cpp")
    words := generateWords(5,"./words/english1000")
    for i, v := range words {
        fmt.Println(i,v)
    }
}

func highlightFile(filename,language string) {
    out, err := exec.Command(fmt.Sprintf("source-highlight -s '%s' -f esc256 -i '%s' -o test.txt",language,filename)).Output()
    out2 , _ := exec.Command("ls").Output()
    fmt.Println(string(out2))
    fmt.Println(out,err)
}

func random(min, max int) int {
    return rand.Intn(max-min)+min
}

func generateWords(n int, filename string) []string {
    reWord := regexp.MustCompile(`[a-zA-Z]+`)
    strings := make([]string,n)
    data, err := os.ReadFile(filename)
    if err != nil {
        panic(err)
    }
    allWordsInFile := reWord.FindAll(data,-1)
    fmt.Println(len(allWordsInFile))
    for i := 0; i < n; i++ {
        strings[i] = string(allWordsInFile[rand.Intn(len(allWordsInFile)-1)])
    }
    return strings
}

