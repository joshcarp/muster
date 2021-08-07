package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

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

func Test_methodStream(t *testing.T) {
	type args struct {
		reader string
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
	}{
		{
			name: "simple",
			args: args{
				reader: `package main

func Blah(s int, a Foo) (int, error) {
	return 0, fmt.Errorf("")
}
`,
			},
			wantWriter: `// Code generated by muster. DO NOT EDIT.
// see https://github.com/joshcarp/muster for more details.
package main

// MustBlah calls Blah and panics if err is not nil.
func MustBlah(s int, a Foo) int {
	val, err := Blah(s, a)
	if err != nil {
		panic(err)
	}
	return val
}
`,
		},
		{
			name: "pointers",
			args: args{
				reader: `package main

func Blah(s *int, a Foo) (*int, error) {
	b := 0
	return &b, fmt.Errorf("")
}
`,
			},
			wantWriter: `// Code generated by muster. DO NOT EDIT.
// see https://github.com/joshcarp/muster for more details.
package main

// MustBlah calls Blah and panics if err is not nil.
func MustBlah(s *int, a Foo) *int {
	val, err := Blah(s, a)
	if err != nil {
		panic(err)
	}
	return val
}
`,
		},
		{
			name: "receiver",
			args: args{
				reader: `package main

type Foo struct {
}

func (f Foo) BlahWithRecv(s int, a Foo) (int, error) {
	return s, fmt.Errorf("")
}
`,
			},
			wantWriter: `// Code generated by muster. DO NOT EDIT.
// see https://github.com/joshcarp/muster for more details.
package main

// MustBlahWithRecv calls BlahWithRecv and panics if err is not nil.
func (f Foo) MustBlahWithRecv(s int, a Foo) int {
	val, err := f.BlahWithRecv(s, a)
	if err != nil {
		panic(err)
	}
	return val
}
`,
		},
		{
			name: "pointer-receiver",
			args: args{
				reader: `package main

type Foo struct {
}

func (f *Foo) BlahWithRecv(s int, a Foo) (int, error) {
	return s, fmt.Errorf("")
}
`,
			},
			wantWriter: `// Code generated by muster. DO NOT EDIT.
// see https://github.com/joshcarp/muster for more details.
package main

// MustBlahWithRecv calls BlahWithRecv and panics if err is not nil.
func (f *Foo) MustBlahWithRecv(s int, a Foo) int {
	val, err := f.BlahWithRecv(s, a)
	if err != nil {
		panic(err)
	}
	return val
}
`,
		},
		{
			name: "external-types",
			args: args{
				reader: `package main

import (
	"fmt"
	foobar "github.com/googleapis/gax-go/v2"
)


func MustSpannerBlah3(s foobar.Backoff, a Foo) (Foo, error) {
	return Foo{}, fmt.Errorf("")
}
`,
			},
			wantWriter: `// Code generated by muster. DO NOT EDIT.
// see https://github.com/joshcarp/muster for more details.
package main

import (
	foobar "github.com/googleapis/gax-go/v2"
)

// MustMustSpannerBlah3 calls MustSpannerBlah3 and panics if err is not nil.
func MustMustSpannerBlah3(s foobar.Backoff, a Foo) Foo {
	val, err := MustSpannerBlah3(s, a)
	if err != nil {
		panic(err)
	}
	return val
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			methodStream(strings.NewReader(tt.args.reader), writer)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("methodStream() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
