// mysql_operators
package exparser

import (
	"errors"
	"fmt"
	//"strings"
)

var MysqlOperators = map[string]*Operator{
	":or:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":nd:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":eq:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":ne:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":gt:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":lt:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":ge:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":le:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":li:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":nl:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":nu:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":nn:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":rl:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
}

func evalMysql(op string, left string, right string) (string, error) {
	switch op {
	case ":or:":
	case ":nd:":
	case ":eq:":
	case ":ne:":
	case ":gt:":
	case ":lt:":
	case ":ge:":
	case ":le:":
	case ":li:":
	case ":nl:":
	case ":nu:":
	case ":nn:":
	case ":rl:":
	default:
	}
	return "", errors.New(fmt.Sprint("Failed to evaluate:", left, op, right))
}
