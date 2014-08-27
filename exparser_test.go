// exparser_test
package exparser

import (
	"fmt"
	"testing"
)

func TestTokenize(t *testing.T) {
	var pass = []struct {
		in string
		ex []string
	}{
		{"-1 + 2", []string{"-1", "+", "2"}},
		{"+1+2", []string{"+1", "+", "2"}},
		{"1+2+(3*4)", []string{"1", "+", "2", "+", "(", "3", "*", "4", ")"}},
		{"1+2+(3*4)^34", []string{"1", "+", "2", "+", "(", "3", "*", "4", ")", "^", "34"}},
		{"'123  456' 789", []string{"123  456", "789"}},
		{`123 "456  789"`, []string{"123", "456  789"}},
		{`123 "456  '''789"`, []string{"123", "456  '''789"}},
	}
	var fail = []struct {
		in string
		ex []string
	}{}
	for _, v := range pass {
		tokens := Tokenize(v.in)
		if !CompareSlices(tokens, v.ex) {
			t.Error("Expected:", v.ex, "actual:", tokens)
		}
	}
	for _, v := range fail {
		tokens := Tokenize(v.in)
		if CompareSlices(tokens, v.ex) {
			t.Error("Expected:", v.ex, "actual:", tokens)
		}
	}
}

func TestEvaluate(t *testing.T) {
	var pass = []struct {
		in string
		ex string
	}{
		//{"-1 + 2^3 + (-1 + 2^3) + (-1 + 2^3)", "3"},
		//{"1.2+2", "3.2"},
		//{"1+2+(3*4)*5", ""},
		//{"1+2+(3*4)^3", ""},
		{"2^3^3", ""},
		//{"3 + 4 * 2 / ( 1-5 ) ^ 2 ^ 3", ""},
	}
	var fail = []struct {
		in string
		ex string
	}{}
	for _, v := range pass {
		//tokens := Tokenize(v.in)
		//_, output, err := ParseRPN(tokens)
		//if err != nil {
		//	fmt.Println(err)
		//}
		////fmt.Println("isDec:", isDec)
		//for o := output.Pop(); o != nil; o = output.Pop() {
		//	//fmt.Print(o.(string), ", ")
		//}
		////fmt.Println()
		fmt.Println(Evaluate(v.in))
	}
	for _, _ = range fail {

	}
}
