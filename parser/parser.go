package parser

import "log"
import . "stack"

func (self *Grammar) Parse(M Table, tokens <-chan *Token) (chan *Token, chan bool) {
    yield := make(chan *Token)
    ack := make(chan bool)
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
        stack.Push(self.ORDER[0])
        X := stack.Peek()
        a := next()
        for X != end {
            if X == a.id {
                yield <- a
                <-ack
                stack.Pop()
                a = next()
            } else if _, has := self.T[X]; has {
                log.Exit("error 1", self.ALL[X], self.ALL[a.id])
            } else if M[X][a.id] == -1 {
                log.Exit("error 2", self.ALL[X], self.ALL[a.id])
            } else {
                yield <- NewToken(X, "")
                <-ack
                stack.Pop()
                p := self.P[M[X][a.id]]
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
