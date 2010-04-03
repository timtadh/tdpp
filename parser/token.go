package parser

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
