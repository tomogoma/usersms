package crdb

import "strings"

// ColDesc returns a string containing cols in the given order separated by ",".
func ColDesc(cols ...string) string {
	desc := ""
	for _, col := range cols {
		if col == "" {
			continue
		}
		desc = desc + col + ", "
	}
	return strings.TrimSuffix(desc, ", ")
}

// ColDescTbl returns a string containing cols for tbl in the given order separated by ",".
func ColDescTbl(tbl string, cols ...string) string {
	desc := ""
	for _, col := range cols {
		if col == "" {
			continue
		}
		desc = desc + tbl + "." + col + ", "
	}
	return strings.TrimSuffix(desc, ", ")
}
