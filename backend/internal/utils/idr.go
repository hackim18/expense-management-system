package utils

import (
	"strconv"
	"strings"
)

func FormatIDR(amount int64) string {
	if amount == 0 {
		return "Rp 0"
	}

	negative := amount < 0
	if negative {
		amount = -amount
	}

	raw := strconv.FormatInt(amount, 10)
	var parts []string
	for len(raw) > 3 {
		parts = append([]string{raw[len(raw)-3:]}, parts...)
		raw = raw[:len(raw)-3]
	}
	if raw != "" {
		parts = append([]string{raw}, parts...)
	}

	formatted := strings.Join(parts, ".")
	if negative {
		formatted = "-" + formatted
	}
	return "Rp " + formatted
}
