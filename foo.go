package main

import (
	"fmt"
	foobar "github.com/googleapis/gax-go/v2"
)

type Foo struct {
}

func (f Foo) BlahWithRecv(s int, a *Foo) (*int, error) {
	b := 0
	return &b, fmt.Errorf("")
}

func Blah(s int, a Foo) (int, error) {
	return 0, fmt.Errorf("")
}

func Blah3(s int, a Foo) (Foo, error) {
	return Foo{}, fmt.Errorf("")
}

func MustSpannerBlah3(s foobar.Backoff, a Foo) (Foo, error) {
	return Foo{}, fmt.Errorf("")
}



