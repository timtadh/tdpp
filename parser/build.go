package parser

import "set"
// import "fmt"

func (self *Grammar) First(i int) (s *set.Set) {
    s = set.New()
    if self.ALL[i].Terminal {
        s.Add(i)
        return
    }
    for _, p := range self.NT[i] {
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
        for _, p := range self.NT[nt] {
            if h, i := has(self.P[p], sym); h {
//                 fmt.Println(self.ALL[sym], self.NT[nt], p, self.P[p], i, len(self.P[p]))
                if i+1 < len(self.P[p]) {
                    f := self.First(self.P[p][i+1])
//                     fmt.Println(f)
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
//     if sym == ORDER[0]: symbols.add('$')
//     tsym = ('nt', sym)
//     for nt in ORDER:
//         for p in PRODUCTIONS[nt]:
//             if tsym in p:
//                 i = p.index(tsym)
//                 if i+1 < len(p):
//                     f = first(p[i+1][1])
//                     if 'e' in f:
//                         f.remove('e')
//                         symbols |= follow(nt)
//                     symbols |= f
//                 elif i+1 == len(p) and sym != nt:
//                     symbols |= follow(nt)
    return s
}
