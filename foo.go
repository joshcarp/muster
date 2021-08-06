package main

import (
	"fmt"
)

type Foo struct {
}

func (f Foo) BlahWithRecv(s int, a Foo) (int, error) {
	return 0, fmt.Errorf("")
}

func Blah(s int, a Foo) (int, error) {
	return 0, fmt.Errorf("")
}

func Blah3(s int, a Foo) (Foo, error) {
	return Foo{}, fmt.Errorf("")
}

