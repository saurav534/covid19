package external

import (
	"covid19/common"
	"sort"
	"strconv"
	"strings"
)

func toInt(input string) int {
	output, _ := strconv.Atoi(strings.ReplaceAll(input, ",", ""))
	return output
}

func sortStates(data []*common.StateData) {
	sort.Slice(data, func(i, j int) bool {
		i1, _ := strconv.Atoi(data[i].Total)
		i2, _ := strconv.Atoi(data[j].Total)
		return i1 > i2
	})
}
