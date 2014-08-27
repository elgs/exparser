package exparser

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

var parentheses = "()"

type Parser struct {
	Opps map[string]int
}

func (this *Parser) Evaluate(expression string) (string, error) {
	tokens := this.Tokenize(expression)
	_, rpn, err := this.ParseRPN(tokens)
	if err != nil {
		return "", err
	}
	return this.Calculate(rpn, true)
}

func (this *Parser) eval(op string, left string, right string) (string, error) {
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

func (this *Parser) Calculate(ts *Lifo, postfix bool) (string, error) {
	newTs := &Lifo{}
	for ti := ts.Pop(); ti != nil; ti = ts.Pop() {
		t := ti.(string)
		switch {
		case this.Opps[t] > 0:
			// operators
			if postfix {
				right := newTs.Pop()
				left := newTs.Pop()
				r, err := this.eval(t, left.(string), right.(string))
				if left == nil || right == nil || err != nil {
					return "", errors.New(fmt.Sprint("Failed to evaluate:", left, t, right))
				}
				newTs.Push(r)
			} else {
				right := ts.Pop()
				left := ts.Pop()
				r, err := this.eval(t, left.(string), right.(string))
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
		this.Calculate(newTs, !postfix)
	}
	return "", errors.New("Error")
}

// false o1 in first, true o2 out first
func (this *Parser) shunt(o1, o2 string) (bool, error) {
	op1 := this.Opps[o1]
	op2 := this.Opps[o2]
	if op1 == 0 || op2 == 0 {
		return false, errors.New(fmt.Sprint("Invalid operators:", o1, o2))
	}
	if op1 < op2 || op1 == op2 && op1%2 == 1 {
		return true, nil
	}
	return false, nil
}

func (this *Parser) ParseRPN(tokens []string) (isDec bool, output *Lifo, err error) {
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
		case this.Opps[token] > 0:
			// operator
			for o2 := opStack.Peep(); o2 != nil; o2 = opStack.Peep() {
				stackToken := o2.(string)
				if this.Opps[stackToken] == 0 {
					break
				}
				o2First, err := this.shunt(token, stackToken)
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

func (this *Parser) Tokenize(exp string) (tokens []string) {
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
		case this.Opps[s] > 0 || strings.ContainsRune(parentheses, v):
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
