package main

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type Foo2 struct {
	A time.Time
}

func ExampleFunction2() {
	spew.Print(Foo2{A: time.Now()})
	//output:

}
func ExampleFunction() {
	f := Function{
		Name: "Foobar",
		Body: `a, err := foo()
if err != nil {
	panic(err)
}
return a`,
		Recv: Reciever{
			Type: "",
			Name: "",
		},
		Params: Params{{
			Name: "a",
			Type: "string",
		}},
		Returns: Returns{{
			Name: "blah",
			Type: "string",
		}},
	}
	fmt.Println(f)
	//output:

}
