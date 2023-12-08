package main

import (
	"strings"
)

func autocomplete[T interface{ GetName() string }](entities []T, val string) (match *T, ok bool) {
	v := strings.ToLower(val)
	for _, entity := range entities {
		c := strings.ToLower(entity.GetName())

		if c != v && strings.HasPrefix(c, v) {
			return &entity, true
		}
	}

	return nil, false
}
