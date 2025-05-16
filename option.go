package main

import (
	"fmt"
)

/*
Represents an optional type, similar to Rust or OCaml.

The idiom is, if isSome is false, then value will be the zero value.

If isSome is true, then the value will be a Some(value).
*/
type optional[V any] struct {
	value  V
	isSome bool
}

func (ot optional[V]) String() string {
	if ot.isSome {
		return fmt.Sprintf("Some %v", ot.value)
	} else {
		return fmt.Sprintf("None")
	}
}
