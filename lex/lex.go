package lex

// import "fmt"
import re "regexp"
import . "parser/token"
import . "parser/grammar"

const skip = 0xff

type MatchHandler func(T map[string]int, str string) *Token
type Handlers map[string]MatchHandler
type Regexs map[string]*re.Regexp

func Lex(gram *Grammar, input <-chan byte, regexs Regexs, handlers Handlers) (<-chan *Token) {
    tokens := make(chan *Token)
    go func() {
        next := func() byte {
            r := <-input
            if closed(input) { return skip }
            return r
        }
        found := func(k string, buf []byte, regex *re.Regexp) []byte {
            var chr byte
            for regex.Match(buf) {
                chr = next()
                if chr == skip { break }
                buf = append(buf, chr)
            }
            var r *Token
            if chr == skip {
                r = handlers[k](gram.TOKENS, string(buf))
                if r != nil { tokens<-r }
                return buf[0:0]
            } else {
                r = handlers[k](gram.TOKENS, string(buf[0:len(buf)-1]))
                if r != nil { tokens<-r }
                buf[0] = buf[len(buf)-1]
                return buf[0:1]
            }
            return nil
        }

        buf := make([]byte, 0, 0xf)
        chr := next()
        if chr == skip { close(tokens); return }
        buf = append(buf, chr)
        for {
            nomatch := true
            for k, regex := range regexs {
                if regex.Match(buf) {
                    buf = found(k, buf, regex)
                    nomatch = false
                    break
                }
            }
            if nomatch {
                chr := next()
                if chr == skip { break }
                buf = append(buf, chr)
            }
        }
        close(tokens)
    }()
    return tokens
}

func append(slice []byte, val byte) []byte {
    length := len(slice)
    if cap(slice) == length {
        // we need to expand
        newsl := make([]byte, length, 2*(length+1))
        for i,v := range slice {
            newsl[i] = v
        }
        slice = newsl
    }
    slice = slice[0:length+1]
    slice[length] = val
    return slice
}
