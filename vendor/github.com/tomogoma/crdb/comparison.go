package crdb

import (
	"fmt"
)

const (

	// Greater than.
	OpGT = ">"
	// Greater than or equal to.
	OpGTOrET = ">="
	// Equal to.
	OpET = "="
	// Less than or equal to.
	OpLTOrET = "<="
	// Less than.
	OpLT        = "<"
	OpIsNull    = "is_null"
	OpIsNotNull = "is_not_null"
)

type Comparison struct {
	Op  string
	Val interface{}
}

// NewComparisonString returns a pointer to Comparison if val is not empty
// otherwise nil.
func NewComparisonString(op string, val string) *Comparison {
	if val == "" {
		return nil
	}
	return &Comparison{Op: op, Val: val}
}

// WhereClause generates an SQL statement part of the form
//     column_one >= $4
// ...or aborts if c is nil.
// Note that the 'WHERE' keyword is NOT included in the result.
// It uses args to determine the index (4) in the statement and returns the same
// args with c.Val appended.
// Mapping in the example is done as follows, col => 'column_one', c.Op => '>='
// and $4 is generated based on len(args).
// c.Op uses the provided c.Op value if no standard operator is found.
func WhereClause(c *Comparison, col string, args []interface{}) (string, []interface{}) {

	if c == nil {
		return "", args
	}

	i := len(args) + 1
	var where string
	switch c.Op {
	case OpGT:
		where = fmt.Sprintf("%s > $%d", col, i)
	case OpGTOrET:
		where = fmt.Sprintf("%s >= $%d", col, i)
	case OpLTOrET:
		where = fmt.Sprintf("%s <= $%d", col, i)
	case OpLT:
		where = fmt.Sprintf("%s < $%d", col, i)
	case OpIsNull:
		where = fmt.Sprintf("%s IS NULL", col)
	case OpIsNotNull:
		where = fmt.Sprintf("%s IS NOT NULL", col)
	case OpET:
		where = fmt.Sprintf("%s = $%d", col, i)
	default:
		where = fmt.Sprintf("%s %s $%d", col, c.Op, i)
	}

	args = append(args, c.Val)

	return where, args
}

// ConcatWhereClause calls WhereClause on c, col, args and concatenates it to a previous
// where clause to have an SQL statement part of the form
//     col_one >= 4 AND col_two IS NULL
// 'col_one >= 4' is picked from where while 'AND' is picked from whereOp and
// 'col_two IS NULL' is generated using WhereClause. See WhereClause for details.
func ConcatWhereClause(c *Comparison, col, where, whereOp string, args []interface{}) (string, []interface{}) {

	if c == nil {
		return where, args
	}

	w, args := WhereClause(c, col, args)
	if where != "" {
		w = fmt.Sprintf("%s %s %s", where, whereOp, w)
	}
	return w, args
}
