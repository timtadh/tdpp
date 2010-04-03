package grammar

import "parser/set"
import "fmt"

type Table [][]int

func (self *Grammar) First(i int) (s *set.Set) {
    s = set.New()
    if self.ALL[i].Terminal {
        s.Add(i)
        return
    }
    for _, p := range self.NTP[i] {
        all_e := true
        for _, sym := range self.P[p] {

            syms := self.First(sym)
            s.Union(syms)
            if !syms.Has(e) {
                all_e = false
                break
            }
        }
        if all_e {
            s.Add(e)
        }
    }
    return
}

func (self *Grammar) Follow(sym int) (s *set.Set) {
    has := func(slice []int, i int) (bool, int) {
        for j,v := range slice {
            if v == i { return true, j }
        }
        return false, 0
    }


    s = set.New()
    if sym == self.ORDER[0] {
        s.Add(end)
    }
    for _, nt := range self.ORDER {
        for _, p := range self.NTP[nt] {
            if h, i := has(self.P[p], sym); h {
                if i+1 < len(self.P[p]) {
                    f := self.First(self.P[p][i+1])
                    if f.Has(e) {
                        f.Remove(e)
                        s.Union(self.Follow(nt))
                    }
                    s.Union(f)
                } else if i+1 == len(self.P[p]) && sym != nt {
                    s.Union(self.Follow(nt))
                }
            }
        }
    }
    return s
}

func (self *Grammar) MakeM() Table {
    M := make(Table, len(self.ALL))
    for i,_ := range M {
        M[i] = make([]int, len(self.ALL))
        for j, _ := range M[i] {
            M[i][j] = -1
        }
    }
    for _, nt := range self.ORDER {
        for _, p := range self.NTP[nt] {
            first := self.First(self.P[p][0])
            for _, sym := range first.Slice() {
                if sym == e { continue }
                if !self.ALL[sym].Terminal { continue }
                M[nt][sym] = p
            }
            if first.Has(e) {
                follow := self.Follow(nt)
                for _, sym := range follow.Slice() {
                    M[nt][sym] = p
                }
            }
        }
    }
    return M
}

func (self Table) String() string {
    s := "     "
    for j, _ := range self[0] {
        s += fmt.Sprintf("%3v  ", j)
    }
    s += "\n---"
    for j := 0; j < len(self[0]); j++ {
        s += "-----"
    }
    s += "\n"
    for i, row := range self {
        s += fmt.Sprintf("%2v|  ", i)
        for j, item := range row {
            if item == -1 {
                s += fmt.Sprintf("%3v", "-")
            } else {
                s += fmt.Sprintf("%3v", item)
            }
            if j+1 != len(row){
                s += fmt.Sprint(", ")
            }
        }
        s += fmt.Sprintln()
    }
    return s
}
