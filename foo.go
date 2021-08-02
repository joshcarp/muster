package main

import (
	"fmt"
)

type Foo struct {
}

func (f Foo) Blah(s int, a Foo) (int, error) {
	return 0, fmt.Errorf("")
}

func Blah(s int, a Foo) (int, error) {
	return 0, fmt.Errorf("")
}
