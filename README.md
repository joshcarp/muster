# Muster

Generate Must functions for functions with two return types where the last type is an error:

```go
type Foo struct {

}

func Blah(s int, a Foo) (int, error) {
    return 0, fmt.Errorf("")
}

```

then run `muster .`
```go
func MustBlah(s int, a Foo) int {
    val, err := Blah(s, a)
    if err != nil {
        panic(err)
    }
    return val
}
```
## Install

```bash
go get github.com/joshcarp/muster
```

## Features
- [ ] Specify more than one file
- [ ] Input from stdin