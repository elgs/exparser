// exparser_test
package exparser

import (
	"fmt"
	"testing"
)

func xTestTokenize(t *testing.T) {

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
		Operators: CalcOperators,
	}
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

func xTestEvaluate(t *testing.T) {
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
		Operators: CalcOperators,
	}
	for _, v := range pass {
		r, err := parser.Calculate(v.in)
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
	var pass = []struct {
		in string
		ex []string
	}{
		{"A:ne:'A (B)'", []string{"A", ":ne:", "A (B)"}},
	}
	var fail = []struct {
		in string
		ex []string
	}{}
	parser := &Parser{
		Operators: MysqlOperators,
	}
	for _, v := range pass {
		tokens := parser.Tokenize(v.in)
		if !CompareSlices(tokens, v.ex) {
			t.Error("Expected:", v.ex, "actual:", tokens)
		}
	}
	for _, _ = range fail {

	}
}

func TestCalculateMySQLFilters(t *testing.T) {
	var pass = []struct {
		in string
		ex string
	}{
		{`f1:gt:'A (B)':or:f2:lt:4:nd:f3:nn:''`, "((F1>'A (B)') OR ((F2<'4') AND (F3 IS NOT NULL)))"},
	}
	var fail = []struct {
		in string
		ex string
	}{}
	parser := &Parser{
		Operators: MysqlOperators,
	}
	for _, v := range pass {
		r, err := parser.Calculate(v.in)
		if err != nil {
			fmt.Println(err)
		}
		if r != v.ex {
			t.Error("Expected:", v.ex, "actual:", r)
		}
	}
	for _, _ = range fail {

	}
}
