package parser

import "log"
import . "stack"

func (self *Grammar) Parse(M Table, tokens []int) (chan int, chan bool) {
    yield := make(chan int)
    ack := make(chan bool)
    go func() {
        next := func(i int) (token, index int) {
            if i == len(tokens) {
                return end, i+1
            }
            return tokens[i], i+1
        }
        var a int
        stack := NewStack()
        index := 0
        stack.Push(end)
        stack.Push(self.ORDER[0])
        X := stack.Peek()
        a, index = next(index)
        for X != end {
            if X == a {
                yield <- a
                <-ack
                stack.Pop()
                a, index = next(index)
            } else if _, has := self.T[X]; has {
                log.Exit("error 1")
            } else if M[X][a] == -1 {
                log.Exit("error 2")
            } else {
                yield <- X
                <-ack
                stack.Pop()
                p := self.P[M[X][a]]
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
    }()
    return yield, ack
}
