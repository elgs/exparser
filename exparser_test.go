// exparser_test
package exparser

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	operators := map[string]int{
		"+": 1,
		"-": 1,
		"*": 3,
		"/": 3,
		"%": 3,
		"^": 4,
	}
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
	parser := &Parser{
		Operators: operators,
	}
	parser.Init()
	for _, v := range pass {
		tokens := parser.Tokenize(v.in)
		if !CompareSlices(tokens, v.ex) {
			t.Error("Expected:", v.ex, len(v.ex), "actual:", tokens, len(tokens))
		}
	}
	for _, v := range fail {
		tokens := parser.Tokenize(v.in)
		if CompareSlices(tokens, v.ex) {
			t.Error("Expected:", v.ex, "actual:", tokens)
		}
	}
}

func TestEvaluate(t *testing.T) {
	operators := map[string]int{
		"+": 1,
		"-": 1,
		"*": 3,
		"/": 3,
		"%": 3,
		"^": 4,
	}
	var pass = []struct {
		in string
		ex string
	}{
		{"-1 + 2^3 + (-1 + 2^3) + (-1 + 2^3)", "21"},
		{"1.2+2", "3.2"},
		{"1+2+(3*4)*5", "63"},
		{"1+2+(3*4)^3", "1731"},
		{"2^3^3", "134217728"},
		{"3 ^4", "81"},
		{"3 + 4 * 2 / ( 1-5 ) ^ 2 ^ 3", "3.0001220703125"},
	}
	var fail = []struct {
		in string
		ex string
	}{}
	parser := &Parser{
		Operators: operators,
	}
	parser.Init()
	for _, v := range pass {
		r, err := parser.Evaluate(v.in)
		if err != nil {
			t.Error(err.Error())
		}
		if r != v.ex {
			t.Error("Expected:", v.ex, "actual:", r)
		}
	}
	for _, _ = range fail {

	}
}

func TestTokenizeFilters(t *testing.T) {
	operators := map[string]int{
		"=":  1,
		"!=": 1,
		">":  1,
		"<":  1,
		">=": 1,
		"<=": 1,
	}
	_ = []string{
		":gt:",
		":lt:",
		"like",
		"null",
		"not",
	}
	var pass = []struct {
		in string
		ex []string
	}{
		{"A!='ABC'", []string{"A", "!=", "ABC"}},
	}
	var fail = []struct {
		in string
		ex []string
	}{}
	parser := &Parser{
		Operators: operators,
	}
	parser.Init()
	for _, v := range pass {
		tokens := parser.Tokenize(v.in)
		if !CompareSlices(tokens, v.ex) {
			t.Error("Expected:", v.ex, "actual:", tokens)
		}
	}
	for _, _ = range fail {

	}
}
