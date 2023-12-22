package hash

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/nspcc-dev/hrw"
)

var re = regexp.MustCompile("^[0-9]+:")

func Sort(key string, nodes []string) []string {
	if len(nodes) < 2 {
		return nodes
	}

	weights, values := toWeightValues(nodes)

	hrw.SortSliceByWeightValue(
		values,
		weights,
		hrw.Hash([]byte(key)),
	)

	return values
}

func toWeightValues(nodes []string) ([]float64, []string) {
	weights := make([]float64, len(nodes))
	values := make([]string, len(nodes))

	for i, v := range nodes {
		if !re.MatchString(v) {
			weights[i] = 1.0
			values[i] = v

			continue
		}

		p := strings.SplitN(v, ":", 2)

		n, err := strconv.ParseFloat(p[0], 64)
		if err != nil {
			panic(err)
		}

		weights[i] = n
		values[i] = p[1]
	}

	max := maxWeight(weights)
	for i, w := range weights {
		weights[i] = w / max
	}

	return weights, values
}

func maxWeight(array []float64) float64 {
	max := 1.0
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}
