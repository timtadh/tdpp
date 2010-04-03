package token

import "fmt"

type Token struct {
    id int
    attr string
}

func NewToken(id int, attr string) *Token {
    self := new(Token)
    self.id = id
    self.attr = attr
    return self
}

func (self *Token) Id() int { return self.id }
func (self *Token) Attr() string { return self.attr }

func (self *Token) String() string {
    return fmt.Sprintf("(%2v, %8v)", self.id, self.attr)
}
