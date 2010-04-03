package parser

import "fmt"
// import "strconv"
import . "stack"
import . "parser/grammar"
import . "parser/token"

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

func collapse(gram *Grammar, stack *Stack) {
    if stack.Len() <= 1 { return }
    top := stack.Peek()
    acc := NewStack()
    for top != nil && len(top.Args) == cap(top.Args) {
        p := stack.Pop()
        fmt.Println("proc -------->", gram.ALL[p.Token.Id()].Name, p.Args)
        acc.Queue(p)
        top = stack.Peek()
    }
    for !acc.Empty() {
        top.AddArg(gram.ALL[acc.Pop().Token.Id()].Name)
    }
    if len(top.Args) == cap(top.Args) {
        collapse(gram, stack)
    }
}

func Process(gram *Grammar, symbols chan *Token, ack chan bool) {
    stack := NewStack()
    for r := range symbols {
        var top *Value
        if stack.Empty() { top = nil } else { top = stack.Peek() }
        if gram.ALL[r.Id()].Terminal {
//             fmt.Printf("%-4v %15v %15v\n", gram.ALL[r.Id()].Name, r.Attr(), "terminal")
            top.AddArg(gram.ALL[r.Id()].Name)
        }
        collapse(gram, stack)
        if !gram.ALL[r.Id()].Terminal {
//             i, _ := strconv.Atoi(r.Attr())
//             fmt.Printf("%-4v %15v %15v\n", gram.ALL[r.Id()].Name, p_to_str(gram.P[i], gram), "nonterminal")
            stack.Push(NewValue(gram, r))
        }
        collapse(gram, stack)
        ack<-true
    }
    if stack.Len() != 1 {
//         collapse(gram, stack)
        fmt.Println(stack)
        fmt.Println("errrrrrrrrrr")
        return
    }
    p := stack.Pop()
    fmt.Println("proc -------->", gram.ALL[p.Token.Id()].Name, p.Args)
}
