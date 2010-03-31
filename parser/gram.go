package parser

import "fmt"
import . "strings"

type Symbol struct {
    Terminal bool
    Name string
    Index int
}

type Production  []*Symbol
type Productions []Production

type Grammar struct {
    T   []string
    NT  []string
    NTP map[string][]int
    P   Productions
}

func (self *Grammar) String() string {
    s := ""
    s += fmt.Sprintln(self.NTP)
    s += fmt.Sprintln(self.NT)
    s += fmt.Sprintln(self.T)
    s += fmt.Sprintln(self.P)
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
    return fmt.Sprintf("(%v %v %v)", self.Terminal, self.Name, self.Index)
}

func NewSymbol(terminal bool, name string, index int) *Symbol {
    self := new(Symbol)
    self.Terminal = terminal
    self.Name = name
    self.Index = index
    return self
}

func MakeGrammar(grammar string) *Grammar {
    self := new(Grammar)
    lines := Split(grammar, "\n", 0)
    self.NTP = make(map[string][]int)
    self.P = make(Productions, 0, len(lines)-1)
    for i, line := range lines {
        if line == "" { continue }
        split := Split(line, "::=", 0)

        nt := TrimSpace(split[0])
        {
            self.NT, _ = IdempotentAppendString(self.NT, nt)
            if _, ok := self.NTP[nt]; !ok {
                self.NTP[nt] = make([]int, 0, 1)
            }
            self.NTP[nt] = AppendInt(self.NTP[nt], i)
        }

        fields := Fields(TrimSpace(split[1]))
        production := make([]*Symbol, len(fields))
        for j, p := range fields {
            var sym *Symbol
            if p[0] == '\'' && p[len(p)-1] == '\'' {
                var k int
                self.T, k = IdempotentAppendString(self.T, p[1:len(p)-1])
                sym = NewSymbol(true, p[1:len(p)-1], k)
            } else if p == "e" {
                sym = NewSymbol(true, p, -1)
            } else {
                var k int
                self.NT, k = IdempotentAppendString(self.NT, p)
                sym = NewSymbol(false, p, k)
            }
            production[j] = sym
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
