package app

import "strings"

func sanitizeStringPointer(str **string) {
	if str != nil && *str != nil {
		**str = strings.TrimSpace(**str)
		if **str == "" {
			*str = nil
		}
	}
}
