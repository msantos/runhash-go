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

	values, weights := toWeightValues(nodes)

	hrw.SortSliceByWeightValue(
		values,
		weights,
		hrw.Hash([]byte(key)),
	)

	return values
}

func toWeightValues(nodes []string) ([]string, []float64) {
	values := make([]string, len(nodes))
	weights := make([]float64, len(nodes))

	max := 1.0

	for i, v := range nodes {
		if !re.MatchString(v) {
			values[i] = v
			weights[i] = 1.0

			continue
		}

		value := strings.SplitN(v, ":", 2)

		weight, err := strconv.ParseFloat(value[0], 64)
		if err != nil {
			panic(err)
		}

		weights[i] = weight
		values[i] = value[1]

		if weight > max {
			max = weight
		}
	}

	if max > 1.0 {
		for i, w := range weights {
			weights[i] = w / max
		}
	}

	return values, weights
}
