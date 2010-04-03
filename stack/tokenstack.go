package stack

import "fmt"
import "strconv"
import "container/list"
import "parser/token"
import "parser/grammar"

type Value struct {
    Terminal bool
    Token *token.Token
    Args  []interface{}
}

func NewValue(gram *grammar.Grammar, token *token.Token) *Value {
    self := new(Value)
    self.Token = token
    self.Terminal = gram.ALL[token.Id()].Terminal
    if !self.Terminal {
        i, _ := strconv.Atoi(token.Attr())
        if gram.P[i][0] != 0 {
            self.Args = make([]interface{}, 0, len(gram.P[i]))
        }
    }
    return self
}

func (self *Value) AddArg(arg interface{}) {
    length := len(self.Args)
    self.Args = self.Args[0:length+1]
    self.Args[length] = arg
}

type Stack struct {
    list *list.List
}

func NewStack() *Stack {
    self := new(Stack)
    self.list = list.New()
    return self
}

func (self *Stack) Len() int { return self.list.Len() }

func (self *Stack) Empty() bool { return self.list.Len() <= 0 }

func (self *Stack) Peek() *Value {
    e := self.list.Front()
    v, _ := e.Value.(*Value)
    return v
}

func (self *Stack) Push(v *Value) {
    self.list.PushFront(v)
}

func (self *Stack) Queue(v *Value) {
    self.list.PushBack(v)
}

func (self *Stack) Pop() *Value {
    e := self.list.Front()
    v, _ := e.Value.(*Value)
    self.list.Remove(e)
    return v
}

func (self *Value) String() string {
    return fmt.Sprintf("<%-5v %20v %2v %v>", self.Terminal, self.Token, cap(self.Args), self.Args)
}

func (self *Stack) String() string {
    s := "stack:\n"
    for e := range self.list.Iter() {
        v, _ := e.(*Value)
        s += fmt.Sprintln("    ", v)
    }
    return s
}
