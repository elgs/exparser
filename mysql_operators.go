// mysql_operators
package exparser

import (
	"errors"
	"fmt"
	"strings"
)

var MysqlOperators = map[string]*Operator{
	":or:": &Operator{
		Precedence: 1,
		Eval:       evalMysql,
	},
	":nd:": &Operator{
		Precedence: 3,
		Eval:       evalMysql,
	},
	":eq:": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	":ne:": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	":gt:": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	":lt:": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	":ge:": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	":le:": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	":li:": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	":nl:": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	":nu:": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	":nn:": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	":rl:": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},

	"::eq::": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	"::ne::": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	"::gt::": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	"::lt::": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	"::ge::": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	"::le::": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	"::li::": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	"::nl::": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
	"::rl::": &Operator{
		Precedence: 5,
		Eval:       evalMysql,
	},
}

func evalMysql(op string, left string, right string) (string, error) {
	left = strings.ToUpper(strings.Replace(left, "-- ", "", -1))

	if op != ":or:" && op != ":nd:" {
		left = strings.ToUpper(strings.Replace(left, "'", "''", -1))
		right = strings.ToUpper(strings.Replace(right, "'", "''", -1))
	}
	switch op {
	case ":or:":
		return fmt.Sprint("(", left, " OR ", right, ")"), nil
	case ":nd:":
		return fmt.Sprint("(", left, " AND ", right, ")"), nil
	case ":eq:":
		return fmt.Sprint("(", left, "='", right, "')"), nil
	case ":ne:":
		return fmt.Sprint("(", left, "!='", right, "')"), nil
	case ":gt:":
		return fmt.Sprint("(", left, ">'", right, "')"), nil
	case ":lt:":
		return fmt.Sprint("(", left, "<'", right, "')"), nil
	case ":ge:":
		return fmt.Sprint("(", left, ">='", right, "')"), nil
	case ":le:":
		return fmt.Sprint("(", left, "<='", right, "')"), nil
	case ":li:":
		return fmt.Sprint("(", left, " LIKE '", right, "')"), nil
	case ":nl:":
		return fmt.Sprint("(", left, " NOT LIKE '", right, "')"), nil
	case ":nu:":
		return fmt.Sprint("(", left, " IS NULL)"), nil
	case ":nn:":
		return fmt.Sprint("(", left, " IS NOT NULL)"), nil
	case ":rl:":
		return fmt.Sprint("(", left, " RLIKE '", right, "')"), nil

	case "::eq::":
		return fmt.Sprint("(", left, "=", right, ")"), nil
	case "::ne::":
		return fmt.Sprint("(", left, "!=", right, ")"), nil
	case "::gt::":
		return fmt.Sprint("(", left, ">", right, ")"), nil
	case "::lt::":
		return fmt.Sprint("(", left, "<", right, ")"), nil
	case "::ge::":
		return fmt.Sprint("(", left, ">=", right, ")"), nil
	case "::le::":
		return fmt.Sprint("(", left, "<=", right, ")"), nil
	case ":li::":
		return fmt.Sprint("(", left, " LIKE ", right, ")"), nil
	case "::nl::":
		return fmt.Sprint("(", left, " NOT LIKE ", right, ")"), nil
	case "::rl::":
		return fmt.Sprint("(", left, " RLIKE ", right, ")"), nil
	}
	return "", errors.New(fmt.Sprint("Failed to evaluate:", left, op, right))
}
