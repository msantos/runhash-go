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

	for n, v := range nodes {
		if !re.MatchString(v) {
			weights[n] = 1.0
			values[n] = v

			continue
		}

		p := strings.SplitN(v, ":", 2)

		i, err := strconv.ParseFloat(p[0], 64)
		if err != nil {
			panic(err)
		}

		weights[n] = i
		values[n] = p[1]
	}

	max := maxWeight(weights)
	for i, w := range weights {
		weights[i] = w / max
	}

	return weights, values
}

func maxWeight(array []float64) float64 {
	max := 1
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}
