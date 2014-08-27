package exparser

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type opp struct {
	op            string
	precedence    uint8
	associativity bool `false: left, true: right`
}

var operators = "+-*/%^"
var parentheses = "()"
var operatorPrecedence = []opp{
	{"+", 2, false},
	{"-", 2, false},
	{"*", 3, false},
	{"/", 3, false},
	{"%", 3, false},
	{"^", 4, true},
}

func Evaluate(expression string) (string, error) {
	tokens := Tokenize(expression)
	_, rpn, err := ParseRPN(tokens)
	if err != nil {
		return "", err
	}
	return Calculate(rpn, true)
}

func eval(op string, left string, right string) (string, error) {
	isDec := strings.Contains(left, ".") || strings.Contains(right, ".") || op == "/"
	switch op {
	case "+":
		if isDec {
			l, err := strconv.ParseFloat(left, 64)
			r, err := strconv.ParseFloat(right, 64)
			return fmt.Sprint(l + r), err
		} else {
			l, err := strconv.ParseInt(left, 10, 64)
			r, err := strconv.ParseInt(right, 10, 64)
			return fmt.Sprint(l + r), err
		}
	case "-":
		if isDec {
			l, err := strconv.ParseFloat(left, 64)
			r, err := strconv.ParseFloat(right, 64)
			return fmt.Sprint(l - r), err
		} else {
			l, err := strconv.ParseInt(left, 10, 64)
			r, err := strconv.ParseInt(right, 10, 64)
			return fmt.Sprint(l - r), err
		}
	case "*":
		if isDec {
			l, err := strconv.ParseFloat(left, 64)
			r, err := strconv.ParseFloat(right, 64)
			return fmt.Sprint(l * r), err
		} else {
			l, err := strconv.ParseInt(left, 10, 64)
			r, err := strconv.ParseInt(right, 10, 64)
			return fmt.Sprint(l * r), err
		}
	case "/":
		if isDec {
			l, err := strconv.ParseFloat(left, 64)
			r, err := strconv.ParseFloat(right, 64)
			if r == 0 {
				return "", errors.New(fmt.Sprint("Failed to evaluate:", left, op, right))
			}
			return fmt.Sprint(l / r), err
		} else {
			l, err := strconv.ParseInt(left, 10, 64)
			r, err := strconv.ParseInt(right, 10, 64)
			if r == 0 {
				return "", errors.New(fmt.Sprint("Failed to evaluate:", left, op, right))
			}
			return fmt.Sprint(l / r), err
		}
	case "%":
		if isDec {
			return "", errors.New(fmt.Sprint("Failed to evaluate:", left, op, right))
		} else {
			l, err := strconv.ParseInt(left, 10, 64)
			r, err := strconv.ParseInt(right, 10, 64)
			if r == 0 {
				return "", errors.New(fmt.Sprint("Failed to evaluate:", left, op, right))
			}
			return fmt.Sprint(l % r), err
		}
	case "^":
		l, err := strconv.ParseFloat(left, 64)
		r, err := strconv.ParseFloat(right, 64)
		if isDec {
			return fmt.Sprint(math.Pow(l, r)), err
		} else {
			return fmt.Sprint(int64(math.Pow(l, r))), err
		}
	}
	return "", errors.New(fmt.Sprint("Failed to evaluate:", left, op, right))
}

func Calculate(ts *Lifo, postfix bool) (string, error) {
	newTs := &Lifo{}
	for ti := ts.Pop(); ti != nil; ti = ts.Pop() {
		t := ti.(string)
		switch {
		case strings.Contains(operators, t):
			// operators
			if postfix {
				right := newTs.Pop()
				left := newTs.Pop()
				r, err := eval(t, left.(string), right.(string))
				if left == nil || right == nil || err != nil {
					return "", errors.New(fmt.Sprint("Failed to evaluate:", left, t, right))
				}
				newTs.Push(r)
			} else {
				right := ts.Pop()
				left := ts.Pop()
				r, err := eval(t, left.(string), right.(string))
				if left == nil || right == nil || err != nil {
					return "", errors.New(fmt.Sprint("Failed to evaluate:", left, t, right))
				}
				newTs.Push(r)
			}
		default:
			// operands
			newTs.Push(t)
		}
		//newTs.Print()
	}
	if newTs.Len() == 1 {
		return newTs.Pop().(string), nil
	} else {
		Calculate(newTs, !postfix)
	}
	return "", errors.New("Error")
}

// false o1 in first, true o2 out first
func shunt(o1, o2 string) (bool, error) {
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
		return false, errors.New(fmt.Sprint("Invalid operators:", o1, o2))
	}
	if op1.precedence < op2.precedence || op1.precedence == op2.precedence && !op1.associativity {
		return true, nil
	}
	return false, nil
}

func ParseRPN(tokens []string) (isDec bool, output *Lifo, err error) {
	opStack := &Lifo{}
	outputQueue := []string{}
	for _, token := range tokens {
		_, err := strconv.ParseFloat(token, 64)
		isNum := err == nil
		switch {
		case isNum:
			if strings.Contains(token, ".") {
				isDec = true
			}
			outputQueue = append(outputQueue, token)
		case strings.Contains(operators, token):
			// operator
			for o2 := opStack.Peep(); o2 != nil; o2 = opStack.Peep() {
				stackToken := o2.(string)
				if !strings.Contains(operators, stackToken) {
					break
				}
				o2First, err := shunt(token, stackToken)
				if err != nil {
					return isDec, output, err
				}
				if o2First {
					outputQueue = append(outputQueue, opStack.Pop().(string))
				} else {
					break
				}
			}
			opStack.Push(token)
		case strings.Contains(parentheses, token):
			// parentheses
			if token == "(" {
				opStack.Push(token)
			} else if token == ")" {
				for o2 := opStack.Pop(); o2 != nil && o2.(string) != "("; o2 = opStack.Pop() {
					outputQueue = append(outputQueue, o2.(string))
				}
			}
		}
	}
	for o2 := opStack.Pop(); o2 != nil; o2 = opStack.Pop() {
		outputQueue = append(outputQueue, o2.(string))
	}
	//fmt.Println(outputQueue)
	output = &Lifo{}
	for i := 0; i < len(outputQueue); i++ {
		(*output).Push(outputQueue[len(outputQueue)-i-1])
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
				if (s == "+" || s == "-") && !n && (len(tokens) == 0 || tokens[len(tokens)-1] != ")") {
					l = false
					n = true
					tmp += s
				} else {
					l = false
					n = false
					tokens = append(tokens, s)
				}

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
