package main

import (
	"fmt"
	"strings"
	"unicode"
)

var operators = "+-*/^"

func Tokenize(exp string) (tokens []string) {
	l, n := false, false
	var tmp string
	for _, v := range exp {
		s := string(v)
		switch {
		case unicode.IsNumber(v) || s == ".":
			if !n && len(tmp) > 0 {
				tokens = append(tokens, tmp)
				tmp = ""
			}
			tmp += s
			l = false
			n = true
		case unicode.IsLetter(v):
			if !l && len(tmp) > 0 {
				tokens = append(tokens, tmp)
				tmp = ""
			}
			tmp += s
			l = true
			n = false
		case unicode.IsSpace(v):
			if len(tmp) > 0 {
				tokens = append(tokens, tmp)
				tmp = ""
			}
			l = false
			n = false
		case s == "'":
			if len(tmp) > 0 {
				tokens = append(tokens, tmp)
				tmp = ""
			}
			l = false
			n = false
		case s == "\"":
			if len(tmp) > 0 {
				tokens = append(tokens, tmp)
				tmp = ""
			}
			l = false
			n = false
		case string(v) == "(":
			if len(tmp) > 0 {
				tokens = append(tokens, tmp)
				tmp = ""
			}
			l = false
			n = false
			tokens = append(tokens, s)
		case s == ")":
			if len(tmp) > 0 {
				tokens = append(tokens, tmp)
				tmp = ""
			}
			l = false
			n = false
			tokens = append(tokens, s)
		case strings.ContainsRune(operators, v):
			if len(tmp) > 0 {
				tokens = append(tokens, tmp)
				tmp = ""
			}
			l = false
			n = false
			tokens = append(tokens, s)
		default:
			if len(tmp) > 0 {
				tokens = append(tokens, tmp)
				tmp = ""
			}
			l = false
			n = false
			fmt.Println("Oops: ", s)
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
