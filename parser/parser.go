package parser

// import "log"
import "fmt"
import . "parser/stack"
import . "parser/grammar"
import . "parser/token"

const e = 0
const end = 1

func Parse(gram *Grammar, M Table, tokens <-chan *Token, lexerr <-chan string) (chan *Token, chan bool, chan string) {
    yield := make(chan *Token)
    ack := make(chan bool)
    errors := make(chan string)

    throw := func(err string) {
        if !closed(errors) {
            errors<-err
        }
//         close(ack)
//         close(yield)
//         close(errors)
    }

    go func() {
        for !closed(lexerr) {
            r := <-lexerr
            if closed(lexerr) { return }
            errors <- r
        }
    }()

    go func() {
        next := func() *Token {
            r := <-tokens
            if closed(tokens) {
                return NewToken(end, "")
            }
            return r
        }
        stack := NewStack()
        stack.Push(end)
        stack.Push(gram.ORDER[0])
        X := stack.Peek()
        a := next()
        for X != end {
            if X == a.Id() {
                yield <- a
                <-ack
                stack.Pop()
                a = next()
            } else if _, has := gram.T[X]; has {
                throw(fmt.Sprint("error 1", gram.ALL[X], gram.ALL[a.Id()]))
                break
            } else if M[X][a.Id()] == -1 {
                throw(fmt.Sprint("error 2", gram.ALL[X], gram.ALL[a.Id()]))
                break
            } else {
                yield <- NewToken(X, fmt.Sprint(M[X][a.Id()]))
                <-ack
                stack.Pop()
                p := gram.P[M[X][a.Id()]]
                for j := len(p)-1; j >= 0; j-- {
                    sym := p[j]
                    if sym == e { continue }
                    stack.Push(sym)
                }
            }
            X = stack.Peek()
        }
        close(ack)
        close(yield)
        close(errors)
    }()
    return yield, ack, errors
}
