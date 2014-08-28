exparser
========

An expression parser.

```go
package main

import (
	"exparser"
	"fmt"
	"os"
)

func main() {
	input := args()
	operators := map[string]int{
		"+": 1,
		"-": 1,
		"*": 3,
		"/": 3,
		"%": 3,
		"^": 4,
	}
	parser := &exparser.Parser{
		Operators: operators,
	}

	r, err := parser.Calculate(input[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(r)
}

func args() []string {
	ret := []string{}
	if len(os.Args) <= 1 {
		fmt.Println("Usage: calc expression")
		os.Exit(0)
	} else {
		ret = os.Args[1:]
	}
	return ret
}

```