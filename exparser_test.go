// exparser_test
package exparser

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	var pass = []struct {
		in string
		ex []string
	}{
		{"1 + 2", []string{"1", "+", "2"}},
		{"1+2", []string{"1", "+", "2"}},
		{"1+2+(3*4)", []string{"1", "+", "2", "+", "(", "3", "*", "4", ")"}},
		{"1+2+(3*4)^34", []string{"1", "+", "2", "+", "(", "3", "*", "4", ")", "^", "34"}},
		{"'123  456' 789", []string{"123  456", "789"}},
		{`123 "456  789"`, []string{"123", "456  789"}},
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
