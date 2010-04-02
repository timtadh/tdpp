package parser

import "log"
import . "stack"

func (self *Grammar) Parse(M Table, tokens chan int) (chan int, chan bool) {
    yield := make(chan int)
    ack := make(chan bool)
    go func() {
        next := func() int {
            r := <-tokens
            if closed(tokens) {
                return end
            }
            return r
        }
        stack := NewStack()
        stack.Push(end)
        stack.Push(self.ORDER[0])
        X := stack.Peek()
        a := next()
        for X != end {
            if X == a {
                yield <- a
                <-ack
                stack.Pop()
                a = next()
            } else if _, has := self.T[X]; has {
                log.Exit("error 1")
            } else if M[X][a] == -1 {
                log.Exit("error 2", self.ALL[X], self.ALL[a])
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
