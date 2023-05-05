package main

import (
	"strings"
)

func IsEmptyMess(s string) bool {
	if strings.TrimSpace(s) == "" {
		return true
	}

	return false
}

func IsPrintable(s string) bool {
	for _, v := range s {
		if !(32 <= v && v <= 127) {
			return false
		}
	}

	return true
}
