# Muster

Generate Must functions for functions with two return types:

```go
type Foo struct {

}

func Blah(s int, a Foo) (int, error) {
    return 0, fmt.Errorf("")
}

```

then run `muster .`
```go
func MustBlah(param0 int, param1 Foo) int {
    val, err := Blah(param0, param1)
    if err != nil {
        panic(err)
    }
return val
}
```