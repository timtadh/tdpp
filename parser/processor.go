package parser

import "log"
// import "os"
import "fmt"
// import "strconv"
import . "stack"
import . "parser/grammar"
import . "parser/token"

type Handler func(interface{}, *Grammar, *Token, []interface{}) interface{}
type Handlers map[string]Handler

func p_to_str(p Production, gram *Grammar) string {
    s := "{"
    for i, sym := range p {
        s += fmt.Sprint(gram.ALL[sym].Name)
        if len(p) == i+1 { continue }
        s += fmt.Sprint(", ")
    }
    s += "}"
    return s
}

func process(statestore interface{}, gram *Grammar, p *Value, handlers Handlers) interface{} {
    hname := gram.HANDLER[p.Production]
//     fmt.Println("proc -------->", gram.ALL[p.Token.Id()].Name, p.Args, hname)s
    return handlers[hname](statestore, gram, p.Token, p.Args)
}

func collapse(statestore interface{}, gram *Grammar, stack *Stack, handlers Handlers) {
    if stack.Len() <= 1 { return }
    top := stack.Peek()
    acc := NewStack()
    for top != nil && len(top.Args) == cap(top.Args) {
        p := stack.Pop()
        acc.Queue(p)
        top = stack.Peek()
    }
    for !acc.Empty() {
        arg := process(statestore, gram, acc.Pop(), handlers)
        top.AddArg(arg)
    }
    if len(top.Args) == cap(top.Args) {
        collapse(statestore, gram, stack, handlers)
    }
}

func Process(statestore interface{}, gram *Grammar, symbols chan *Token, ack chan bool, parserr chan string, handlers Handlers) (interface{}, bool) {
    noerrors := make(chan bool)
    go func() {
//         log.Stderr("error starting")
        noerr := true
        for !closed(parserr) {
//             log.Stderr("error 1")
            r := <-parserr
//             log.Stderr("error 2")
            if closed(parserr) { break }
//             log.Stderr("error 3")
            log.Stderr(r)
//             log.Stderr("error 4")
            noerr = false
        }
//         log.Stderr("error forloop done")
        if !noerr { noerrors <- false
        } else { noerrors <- true }
        close(noerrors)
//         log.Stderr("error done")
    }()

    success := make(chan bool)
    var final *Value
    go func() {
        stack := NewStack()
        for r := range symbols {
            var top *Value
            if stack.Empty() { top = nil } else { top = stack.Peek() }
            if gram.ALL[r.Id()].Terminal {
//                 fmt.Printf("%-4v %15v %15v\n", gram.ALL[r.Id()].Name, r.Attr(), "terminal")
                top.AddArg(r.Attr())
            }
            collapse(statestore, gram, stack, handlers)
            if !gram.ALL[r.Id()].Terminal {
//                 i, _ := strconv.Atoi(r.Attr())
//                 fmt.Printf("%-4v %15v %15v\n", gram.ALL[r.Id()].Name, p_to_str(gram.P[i], gram), "nonterminal")
                stack.Push(NewValue(gram, r))
            }
            collapse(statestore, gram, stack, handlers)
            ack<-true
        }
        if stack.Len() != 1 {
//             log.Stderr("fatal processing error in parser, symbols still on stack", stack)
            success<-false
        } else {
            final = stack.Pop()
            success<-true
        }
        close(success)
    }()
//     fmt.Println("proc -------->", gram.ALL[p.Token.Id()].Name, p.Args)
    r1 := <-success
    r2 := <-noerrors
    <-success
    <-noerrors
//     log.Stderr("main 1")
    ret := r1 && r2
    if !ret {
        return nil, false
    }
    return process(statestore, gram, final, handlers), ret
}
