package parser

import "fmt"
import . "strings"

const e = 0
const end = 1

type Symbol struct {
    Terminal bool
    Name string
}

type Production  []int
type Productions []Production

type Grammar struct {
    TOKENS  map[string]int
    T       map[int]int
    NT      map[int]int
    NTP     map[int][]int
    P       Productions
    ALL     []*Symbol
    ORDER   []int
}

func (self *Grammar) String() string {
    s := ""
    s += fmt.Sprintln("TOKENS: ", self.TOKENS)
    s += fmt.Sprintln("T: ", self.T)
    s += fmt.Sprintln("NT: ", self.NT)
    s += fmt.Sprintln("NTP: ", self.NTP)
    s += fmt.Sprintln("P: ", self.P)
    s += "ALL:\n"
    for i, sym := range self.ALL {
        s += fmt.Sprintf("  %v = %v\n", i, sym)
    }
    s += fmt.Sprintln("\nORDER: ", self.ORDER)
    return s
}

func (self Productions) String() string {
    s := ""
    for i,p := range self {
        s += fmt.Sprintln(i, p)
    }
    return s
}

func (self *Symbol) String() string {
    return fmt.Sprintf("(%v %v)", self.Terminal, self.Name)
}

func NewSymbol(terminal bool, name string) *Symbol {
    self := new(Symbol)
    self.Terminal = terminal
    self.Name = name
    return self
}

func (self *Symbol) Eq(b *Symbol) bool {
    return self.Terminal == b.Terminal && self.Name == b.Name
}

func MakeGrammar(grammar string) *Grammar {
    self := new(Grammar)
    lines := Split(grammar, "\n", 0)
    self.TOKENS = make(map[string]int)
    self.T = make(map[int]int)
    self.NT = make(map[int]int)
    self.NTP = make(map[int][]int)
    self.P = make(Productions, 0, len(lines)-1)
    self.ALL, _ = IdempotentAppendSymbol(self.ALL, NewSymbol(true, "e"))
    self.ALL, _ = IdempotentAppendSymbol(self.ALL, NewSymbol(true, "$"))
    self.T[0] = 0
    self.T[1] = 1
    for i, line := range lines {
        if line == "" { continue }
        split := Split(line, "::=", 0)

        {
            var k int
            nt := TrimSpace(split[0])
            self.ALL, k = IdempotentAppendSymbol(self.ALL, NewSymbol(false, nt))
            if _, ok := self.NTP[k]; !ok {
                self.NT[k] = len(self.NT)
                self.NTP[k] = make([]int, 0, 1)
                self.ORDER = AppendInt(self.ORDER, k)
            }
            self.NTP[k] = AppendInt(self.NTP[k], i)
        }

        fields := Fields(TrimSpace(split[1]))
        production := make([]int, len(fields))
        for j, p := range fields {
            var sym *Symbol
            if p[0] == '\'' && p[len(p)-1] == '\'' {
                sym = NewSymbol(true, p[1:len(p)-1])
            } else if p == "e" {
                sym = NewSymbol(true, p)
            } else {
                sym = NewSymbol(false, p)
            }
            var k int
            self.ALL, k = IdempotentAppendSymbol(self.ALL, sym)
            production[j] = k
            if sym.Terminal {
                if _, has := self.T[k]; !has {
                    self.T[k] = len(self.T)
                    self.TOKENS[sym.Name] = k
                }
            }
        }
        self.P = AppendProduction(self.P, production)
    }
    return self
}

func AppendProduction(slice Productions, val Production) Productions {
    length := len(slice)
    if cap(slice) == length {
        // we need to expand
        newsl := make(Productions, length, 2*(length+1))
        for i,v := range slice {
            newsl[i] = v
        }
        slice = newsl
    }
    slice = slice[0:length+1]
    slice[length] = val
    return slice
}

func AppendInt(slice []int, val int) []int {
    length := len(slice)
    if cap(slice) == length {
        // we need to expand
        newsl := make([]int, length, 2*(length+1))
        for i,v := range slice {
            newsl[i] = v
        }
        slice = newsl
    }
    slice = slice[0:length+1]
    slice[length] = val
    return slice
}

func IdempotentAppendString(slice []string, val string) ([]string, int) {
    for i,v := range slice {
        if v == val {
            return slice, i
        }
    }
    length := len(slice)
    if cap(slice) == length {
        // we need to expand
        newsl := make([]string, length, 2*(length+1))
        for i,v := range slice {
            newsl[i] = v
        }
        slice = newsl
    }
    slice = slice[0:length+1]
    slice[length] = val
    return slice, length
}


func IdempotentAppendSymbol(slice []*Symbol, val *Symbol) ([]*Symbol, int) {
    for i,v := range slice {
        if v.Eq(val) {
            return slice, i
        }
    }
    length := len(slice)
    if cap(slice) == length {
        // we need to expand
        newsl := make([]*Symbol, length, 2*(length+1))
        for i,v := range slice {
            newsl[i] = v
        }
        slice = newsl
    }
    slice = slice[0:length+1]
    slice[length] = val
    return slice, length
}
