package main

import "os"
import "log"
import "fmt"
import "parser"

func getstdin() string {
    bytes := make([]byte, 1000)
    n, err := os.Stdin.Read(bytes)
    if err != nil {
        log.Exit("Could not read from stdin", err)
    }
    return string(bytes[0:n])
}

func main() {
    grammar := getstdin()
    fmt.Println(grammar)
    gram := parser.MakeGrammar(grammar)
    fmt.Println(gram)
    for _, k := range gram.ORDER {
        fmt.Println(gram.ALL[k])
        fmt.Println("   FIRST:", gram.First(k))
        fmt.Println("  FOLLOW:", gram.Follow(k))
    }
    M := gram.MakeM()
    fmt.Println(M)
    results, ack := gram.Parse(M, []int{15, 16, 15, 11, 15, 16, 15})
    for r := range results {
        fmt.Println(gram.ALL[r])
        ack<-true
    }
}
