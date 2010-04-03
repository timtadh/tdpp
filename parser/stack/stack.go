package stack

import "container/list"

type Stack struct {
    list *list.List
}

func NewStack() *Stack {
    self := new(Stack)
    self.list = list.New()
    return self
}

func (self *Stack) Empty() bool { return self.list.Len() <= 0 }

func (self *Stack) Peek() int {
    e := self.list.Front()
    t, _ := e.Value.(int)
    return t
}

func (self *Stack) Push(t int) {
    self.list.PushFront(t)
}

func (self *Stack) Pop() int {
    e := self.list.Front()
    t, _ := e.Value.(int)
    self.list.Remove(e)
    return t
}
