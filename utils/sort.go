package utils

import (
	"fmt"
	"strings"
)

// SortQuery for
func SortQuery(sort string) string {
	var sortQuery string

	if sort != "" {
		splitted := strings.Split(sort, ",")
		for _, item := range splitted {
			if string(item[0]) == "-" {
				sortQuery += fmt.Sprintf("%s %s, ", string(item[1:]), "DESC")
			} else {
				sortQuery += fmt.Sprintf("%s, ", item)
			}
		}
		sortQuery = sortQuery[:len(sortQuery)-2]
	}
	return sortQuery
}
