// Code generated by muster. DO NOT EDIT.
package main

func MustBlah(param0 int, param1 Foo) int {
	val, err := Blah(param0, param1)
	if err != nil {
		panic(err)
	}
	return val
}
