package main

import (
    "fmt"
    "os/exec"
    "os/signal"
    "os"
    "regexp"
    "math/rand"
)

func readTest() {
    var b []byte = make([]byte,1)
    os.Stdin.Read(b)
    fmt.Println(string(b))
}

func play() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs,syscall.SIGINT)

}

func main() {
    fmt.Println("Hello World")
    // highlightFile("test.go","cpp")
    words := generateWords(5,"./words/english1000")
    for i, v := range words {
        fmt.Println(i,v)
    }
    // disable input buffering
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    // do not display entered characters on the screen
    exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
    readTest()
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

