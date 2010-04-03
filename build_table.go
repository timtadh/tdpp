package main

import "os"
import "log"
import "fmt"
import "parser/grammar"

func getstdin() string {
    bytes := make([]byte, 1000)
    n, err := os.Stdin.Read(bytes)
    if err != nil {
        log.Exit("Could not read from stdin", err)
    }
    return string(bytes[0:n])
}

func main() {
    g := getstdin()
    fmt.Println(g)
    gram := grammar.MakeGrammar(g)
    fmt.Println(gram)
    for _, k := range gram.ORDER {
        fmt.Println(gram.ALL[k])
        fmt.Println("   FIRST:", gram.First(k))
        fmt.Println("  FOLLOW:", gram.Follow(k))
    }
}
