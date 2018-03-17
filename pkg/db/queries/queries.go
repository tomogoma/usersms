package queries

import (
	"fmt"
	"errors"
	"strings"
)

const (
	OrderAsc  = "ASC"
	OrderDesc = "DESC"

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

type ColOrders struct {
	Orders []string
	Cols   []interface{}
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

// JoinWhereClause calls WhereClause on c, col, args and joins it to a previous
// where clause to have an SQL statement part of the form
//     col_one >= 4 AND col_two IS NULL
// 'col_one >= 4' is picked from where while 'AND' is picked from whereOp and
// 'col_two IS NULL' is generated using WhereClause. See WhereClause for details.
func JoinWhereClause(c *Comparison, col, where, whereOp string, args []interface{}) (string, []interface{}) {

	if c == nil {
		return where, args
	}

	w, args := WhereClause(c, col, args)
	if where != "" {
		w = fmt.Sprintf("%s %s %s", where, whereOp, w)
	}
	return w, args
}

// OrderBy generates an SQL statement part of the form
//     ORDER BY col_one ASC, col_two DESC
// It determines this using co. mapFunc is used to map the col in each Col
// to a real column e.g. col_one of tbl_a.col_one. The error returned by
// mapFunc is returned as is and OrderBy stops generating.
func OrderBy(co *ColOrders, mapFunc func(modelCol interface{}) (string, error)) (string, error) {

	if co == nil {
		return "", nil
	}

	numCols := len(co.Cols)
	if numCols != len(co.Orders) {
		return "", errors.New("number of columns and orders not equal in ColOrders")
	}

	if numCols == 0 {
		return "", nil
	}

	var orderBy, currOrder, currCol string
	var err error
	for i, modelCol := range co.Cols {

		switch co.Orders[i] {
		case OrderAsc:
			currOrder = "ASC"
		case OrderDesc:
			currOrder = "DESC"
		default:
			return "", fmt.Errorf("unknown order in ColOrders: %s", co.Orders[i])
		}

		currCol, err = mapFunc(modelCol)
		if err != nil {
			return "", err
		}

		orderBy = fmt.Sprintf("%s%s %s, ", orderBy, currCol, currOrder)
	}

	return strings.TrimSuffix("ORDER BY "+orderBy, ", "), nil
}

// Pagination generates an SQL statement part of the form
//     LIMIT $3 OFFSET $4
// It appends offset and count to the args and uses the length of the args
// to determine the index of the offset and count args (3 and 4) before
// returning the new set of args.
func Pagination(offset, count int64, args []interface{}) (string, []interface{}) {
	args = append(args, count)
	p := fmt.Sprintf("LIMIT $%d", len(args))
	args = append(args, offset)
	return fmt.Sprintf("%s OFFSET $%d", p, len(args)), args
}
