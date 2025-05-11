package main

func sliceGet[V any](s []V, compFunc func(V) bool) (V, bool) {
	for _, item := range s {
		if compFunc(item) {
			return item, true
		}
	}
	var v V
	return v, false
}
