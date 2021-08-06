package main

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"go/format"
	"log"
	"text/template"
	"unicode"
)

type Param struct {
	Name string
	Type string
}

type Params []Param

func (p Params) String() string {
	str, err := WithTemplate(`{{range $i, $e := .}}{{.Name}} {{.Type}},{{end}}`, p)
	if err != nil {
		return ""
	}
	return str
}

type Return struct {
	Name string
	Type string
}

type Returns []Return

func (p Returns) String() string {
	str, err := WithTemplate(`{{range $i, $e := .}}{{.Name}} {{.Type}},{{end}}`, p)
	if err != nil {
		return ""
	}
	return str
}

type Reciever struct {
	Type string
	Name string
}

func (f Reciever) String() string {
	if f.Type == "" {
		return ""
	}
	str, err := WithTemplate(`({{.Name}} {{.Type}})`, f)
	if err != nil {
		panic(err)
	}
	return str
}

type Function struct {
	Name    string
	Body    string
	Recv    Reciever
	Params  Params
	Returns Returns
}

func (f Function) IsExported() bool {
	if len(f.Name) == 0 {
		return false
	}
	return unicode.IsUpper(rune(f.Name[0]))

}
func (f Function) String() string {
	str, err := WithTemplate(`func {{.Recv}}{{.Name}}({{.Params}})({{.Returns}}){
	{{.Body}}
}
`, f)
	if err != nil {
		panic(err)
	}
	return FormatCode(str)
}

func WithTemplate(tmplstr string, data interface{}, funcs ...interface{}) (string, error) {
	funcmap := sprig.FuncMap()
	err := extraFuncs(funcmap, funcs...)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("anzdata").
		Funcs(map[string]interface{}(funcmap)).
		Parse(tmplstr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), err
}

func extraFuncs(m map[string]interface{}, funcs ...interface{}) error {
	if len(funcs)%2 != 0 {
		return fmt.Errorf("extra funcs should be even with form ['funcname', func...]")
	}
	for i := 0; i < len(funcs)-1; i += 2 {
		key, ok := funcs[i].(string)
		if !ok {
			return fmt.Errorf("key of wrong type, key should be string type")
		}
		val := funcs[i+1]
		m[key] = val
	}
	return nil
}

func FormatCode(buf string) string {
	src, err := format.Source([]byte(buf))
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return buf
	}
	return string(src)
}
