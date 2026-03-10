package service

import "strconv"

func GeneratePageOSSKey(index int) string {
	return "page_" + strconv.Itoa(index)
}
