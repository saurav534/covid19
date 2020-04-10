package external

import (
	"strconv"
	"strings"
)

func toInt(input string) int {
	output, _ := strconv.Atoi(strings.ReplaceAll(input, ",", ""))
	return output
}
