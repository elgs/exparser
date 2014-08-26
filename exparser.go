package exparser

import (
	"fmt"
	"strings"
	"unicode"
)

var operators = "+-*/^()"

func Tokenize(exp string) (tokens []string) {
	sq, dq, l, n := false, false, false, false
	var tmp string
	for _, v := range exp {
		s := string(v)
		switch {
		case unicode.IsNumber(v) || s == ".":
			if !n && len(tmp) > 0 && !sq && !dq {
				tokens = append(tokens, tmp)
				tmp = ""
			}
			tmp += s
			l = false
			n = true
		case unicode.IsLetter(v):
			if !l && len(tmp) > 0 && !sq && !dq {
				tokens = append(tokens, tmp)
				tmp = ""
			}
			tmp += s
			l = true
			n = false
		case unicode.IsSpace(v):
			if sq || dq {
				tmp += s
			} else {
				if len(tmp) > 0 && !sq && !dq {
					tokens = append(tokens, tmp)
					tmp = ""
				}
				l = false
				n = false
			}
		case s == "'":
			if !dq {
				sq = !sq
			}
			if !sq && !dq {
				if len(tmp) > 0 {
					tokens = append(tokens, tmp)
					tmp = ""
				}
				l = false
				n = false
			}
		case s == "\"":
			if !sq {
				dq = !dq
			}
			if !sq && !dq {
				if len(tmp) > 0 {
					tokens = append(tokens, tmp)
					tmp = ""
				}
				l = false
				n = false
			}
		case strings.ContainsRune(operators, v):
			if sq || dq {
				tmp += s
			} else {
				if len(tmp) > 0 {
					tokens = append(tokens, tmp)
					tmp = ""
				}
				l = false
				n = false
				tokens = append(tokens, s)
			}
		default:
			if sq || dq {
				tmp += s
			} else {
				if len(tmp) > 0 {
					tokens = append(tokens, tmp)
					tmp = ""
				}
				l = false
				n = false
			}
		}
	}
	if len(tmp) > 0 {
		tokens = append(tokens, tmp)
		tmp = ""
	}
	return
}

func main() {
	s := "1.2 + 3 * (4-5)"
	tokens := Tokenize(s)
	for i, v := range tokens {
		fmt.Println(i, v)
	}

}
