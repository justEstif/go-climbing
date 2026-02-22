package components

import "strconv"

func gradeValue(n int) string {
	return strconv.Itoa(n)
}

func gradeName(n int) string {
	return "V" + strconv.Itoa(n)
}
