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
}
