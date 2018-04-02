package crdb

import (
	"errors"
	"fmt"
	"strings"
)

const (
	OrderAsc  = "ASC"
	OrderDesc = "DESC"
)

type ColOrders struct {
	Orders []string
	Cols   []interface{}
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
