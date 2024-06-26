package utils

import (
	"strings"
)

// Converts a string to CamelCase
func ToCamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := true
	prevIsCap := false
	for i, v := range []byte(s) {
		isCap := v >= 'A' && v <= 'Z'
		isLow := v >= 'a' && v <= 'z'

		if capNext || i == 0 {
			if isLow {
				v += 'A'
				v -= 'a'
			}
		} else if prevIsCap && isCap {
			v += 'a'
			v -= 'A'
		}

		prevIsCap = isCap

		if isCap || isLow {
			n.WriteByte(v)
			capNext = false
		} else if isNum := v >= '0' && v <= '9'; isNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}

// Check string is empty
func IsEmpty(s string) bool {
	return strings.Trim(s, " ") == ""
}
