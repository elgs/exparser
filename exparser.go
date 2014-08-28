package exparser

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type Parser struct {
	Operators map[string]int
	maxOpLen  int
	Keywords  []string
}

func (this *Parser) Init() {
	for k, _ := range this.Operators {
		if len(k) > this.maxOpLen {
			this.maxOpLen = len(k)
		}
	}
}

func (this *Parser) Evaluate(expression string) (string, error) {
	tokens := this.Tokenize(expression)
	//fmt.Println(expression, tokens)
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
		case this.Operators[t] > 0:
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
	op1 := this.Operators[o1]
	op2 := this.Operators[o2]
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
		case this.Operators[token] > 0:
			// operator
			for o2 := opStack.Peep(); o2 != nil; o2 = opStack.Peep() {
				stackToken := o2.(string)
				if this.Operators[stackToken] == 0 {
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
		case token == "(":
			opStack.Push(token)
		case token == ")":
			for o2 := opStack.Pop(); o2 != nil && o2.(string) != "("; o2 = opStack.Pop() {
				outputQueue = append(outputQueue, o2.(string))
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
	sq, dq := false, false
	var tmp string
	expRunes := []rune(exp)
	for i := 0; i < len(expRunes); i++ {
		v := expRunes[i]
		s := string(v)
		switch {
		case unicode.IsSpace(v):
			if sq || dq {
				tmp += s
			} else {
				if len(tmp) > 0 && !sq && !dq {
					tokens = append(tokens, tmp)
					tmp = ""
				}
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
				}
			}
		case s == "+" || s == "-" || s == "(" || s == ")":
			if sq || dq {
				tmp += s
			} else {
				if len(tmp) > 0 {
					tokens = append(tokens, tmp)
					tmp = ""
				}
				lastToken := ""
				if len(tokens) > 0 {
					lastToken = tokens[len(tokens)-1]
				}
				if (s == "+" || s == "-") && (len(tokens) == 0 || lastToken == "(" || this.Operators[lastToken] > 0) {
					// sign
					tmp += s
				} else {
					// operator
					tokens = append(tokens, s)
				}
			}
		default:
			if sq || dq {
				tmp += s
			} else {
				// until the max length of operators(n), check if next 1..n runes are operator, greedily
				opCandidateTmp := ""
				opCandidate := ""
				for j := 0; j < this.maxOpLen && i < len(expRunes)-1; j++ {
					next := string(expRunes[i+j])
					opCandidateTmp += next
					if this.Operators[opCandidateTmp] > 0 {
						opCandidate = opCandidateTmp
					}
				}
				if len(opCandidate) > 0 {
					if len(tmp) > 0 {
						tokens = append(tokens, tmp)
						tmp = ""
					}
					tokens = append(tokens, opCandidate)
					i += len(opCandidate) - 1
				} else {
					tmp += s
				}
			}
		}
	}
	if len(tmp) > 0 {
		tokens = append(tokens, tmp)
		tmp = ""
	}
	return
}
