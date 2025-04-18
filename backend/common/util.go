package common

import (
	"fmt"
	"slices"
)

// delete by reference
func RemoveIndexSlice[K any](s []K, i int) []K {
	return slices.Delete(s, i, i+1)
}

// make a copy and delete the number from that copy
func RemoveIndexSliceCopied[K any](s []K, i int) []K {
	ret := make([]K, len(s)-1)
	ret = append(ret, s[:i]...)
	return append(ret, s[i+1:]...)
}

func MapToString[K comparable, V any](items map[K]V) string {
	if len(items) == 0 {
		return "[]"
	}

	s := "{\n"
	suf := ",\n"
	for key, value := range items {
		s = fmt.Sprintf("%s%v: %v%s", s, key, value, suf)
	}
	s = fmt.Sprintf("%s}", s)

	return s
}

func SliceToString[K any](items []K) string {
	if len(items) == 0 {
		return "[]"
	}

	s := "[\n"
	suf := ",\n"
	for index, item := range items {
		if index == len(items)-1 {
			suf = "\n]"
		}
		s = fmt.Sprintf("%s%v%s", s, item, suf)
	}

	return s
}
