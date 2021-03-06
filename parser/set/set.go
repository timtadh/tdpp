package set

import "fmt"

type Set struct {
    set map[int]bool
}

func New() *Set {
    self := new(Set)
    self.set = make(map[int]bool)
    return self
}

func (self *Set) Add(i int) {
    self.set[i] = true
}

func (self *Set) Remove(i int) {
    self.set[i] = false, false
}

func (self *Set) Union(b *Set) {
    for k,_ := range b.set {
        self.set[k] = true
    }
}

func (self *Set) Has(i int) bool {
    _, has := self.set[i]
    return has
}

func (self *Set) Slice() (slice []int) {
    slice = make([]int, len(self.set))
    i := 0
    for k,_ := range self.set {
        slice[i] = k
        i++
    }
    return
}

func (self *Set) String() string {
    s := "<"
    for k, _ := range self.set {
        s += fmt.Sprint(k) + ", "
    }
    s += ">"
    return s
}
