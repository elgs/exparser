package exparser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type opp struct {
	op            string
	precedence    uint8
	associativity bool `false: left, true: right`
}

var operators = "+-*/^"
var parentheses = "()"
var operatorPrecedence = []opp{
	{"+", 2, false},
	{"-", 2, false},
	{"*", 3, false},
	{"/", 3, false},
	{"^", 4, true},
}

func shunt(o1, o2 string) (string, error) {
	op1Valid, op2Valid := false, false
	var op1 opp
	var op2 opp
	for _, v := range operatorPrecedence {
		if v.op == o1 {
			op1 = v
			op1Valid = true
		}
		if v.op == o2 {
			op2 = v
			op2Valid = true
		}
		if op1Valid && op2Valid {
			break
		}
	}
	if !op1Valid || !op2Valid {
		return "", errors.New(fmt.Sprint("Invalid operators:", o1, o2))
	}
	if op1.precedence < op2.precedence || op1.precedence == op2.precedence && !op1.associativity {
		return op2.op, nil
	}
	return op1.op, nil
}

func ParseRPN(tokens []string) (isDec bool, output Lifo, err error) {
	opStack := &Lifo{}
	outputQueu := []string{}
	for _, token := range tokens {
		_, err := strconv.ParseFloat(token, 64)
		isNum := err != nil
		switch {
		case isNum:
			if strings.Contains(token, ".") {
				isDec = true
			}
			outputQueu = append(outputQueu, token)
		case strings.Contains(operators, token):
			// operator
			o2 := opStack.Peep()
			for o2 != nil {
				stackToken := o2.(string)
				op, err := shunt(token, stackToken)
				if err != nil {
					return isDec, output, err
				}
				if op == stackToken {
					outputQueu = append(outputQueu, opStack.Pop().(string))
				}
				o2 = opStack.Peep()
			}
			opStack.Push(token)
		case strings.Contains(parentheses, token):
			// parentheses
		}
	}
	return
}

func Tokenize(exp string) (tokens []string) {
	sq, dq, l, n := false, false, false, false
	var tmp string
	for _, v := range exp {
		s := string(v)
		switch {
		case unicode.IsNumber(v) || s == ".":
			if !sq && !dq {
				if !n && len(tmp) > 0 {
					tokens = append(tokens, tmp)
					tmp = ""
				}
				l = false
				n = true
			}
			tmp += s
		case unicode.IsLetter(v):
			if !sq && !dq {
				if !l && len(tmp) > 0 {
					tokens = append(tokens, tmp)
					tmp = ""
				}
				l = true
				n = false
			}
			tmp += s
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
			if dq {
				tmp += s
			} else {
				sq = !sq
				if !sq {
					if len(tmp) > 0 {
						tokens = append(tokens, tmp)
						tmp = ""
					}
					l = false
					n = false
				}
			}
		case s == "\"":
			if sq {
				tmp += s
			} else {
				dq = !dq
				if !dq {
					if len(tmp) > 0 {
						tokens = append(tokens, tmp)
						tmp = ""
					}
					l = false
					n = false
				}
			}
		case strings.ContainsRune(operators, v) || strings.ContainsRune(parentheses, v):
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
