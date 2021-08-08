# Muster

muster is for API developers to automatically generate Must functions for any function that returns a value and an error.

Sometimes it's annoying to create test structs when you need to populate a field with a function that returns (value, error).

## Install 📥

```bash
go get github.com/joshcarp/muster
```

## Problem 🔥
```go
func TestSomeComplicatedFunc(t *testing.T){
	value, _ := someFunction("blah") // This is annoying
	testCase := SomeComplicatedStruct{
                    Foo: "Bar"
                    ComplicatedField:value
}
//...
}
```

A better alternative to this is if there was a `mustSomethingFunction(string)interface{}` that panics if there is an error:

```go
func TestSomeComplicatedFunc(t *testing.T){
	testCase := SomeComplicatedStruct{
                    Foo: "Bar"
                    ComplicatedField: mustSomeFunction("blah") // Much cleaner
}
//...
}
```


## Example 🔧
```go
type Foo struct {

}

func Blah(s int, a Foo) (int, error) {
    return 0, fmt.Errorf("")
}

```

then run `muster .` or `muster <filename>.go` or `cat <filename> | muster --stream > output.txt` 
```go
func MustBlah(s int, a Foo) int {
    val, err := Blah(s, a)
    if err != nil {
        panic(err)
    }
    return val
}
```

## Features 💯
- [x] Specify more than one file
- [x] Input from stdin
- [x] Works with methods
- [x] Works with external types
- [x] Include regex
- [x] Exclude regex
- [x] Variadic functions
