package helpers

import (
	"fmt"
	"strings"
)

// AveragePrice computes average of top 5 prices (used when scraping; prices are string slices)
func AveragePrice(prices []string) float64 {
	sum := 0.0
	count := 0
	max := 5
	if len(prices) < max {
		max = len(prices)
	}
	for i := 0; i < max; i++ {
		p, ok := parsePrice(prices[i])
		if !ok || p > 10000 {
			continue
		}
		count++
		sum += p
	}
	if count == 0 {
		return -1
	}
	return sum / float64(count)
}

func parsePrice(s string) (float64, bool) {
	if len(s) == 0 {
		return 0, false
	}
	if s[0] == '$' {
		s = s[1:]
	}
	s = strings.ReplaceAll(s, ",", "")
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err == nil
}
