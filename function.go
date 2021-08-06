package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"text/template"
	"unicode"

	"github.com/Masterminds/sprig"
)

func FromDecl(decl *ast.FuncDecl, fset *token.FileSet) Function {
	a := Function{Name: decl.Name.Name}
	for _, e := range decl.Type.Params.List {
		a.Params = append(a.Params, Param{
			Name: fmt.Sprintf("%s", e.Names[0].Name),
			Type: NodeAsString(fset, e.Type),
		})
	}
	for _, e := range decl.Type.Results.List {
		a.Returns = append(a.Returns, Return{
			Type: NodeAsString(fset, e.Type),
		})
	}
	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		a.Recv.Type = NodeAsString(fset, decl.Recv.List[0].Type)
		a.Recv.Name = NodeAsString(fset, decl.Recv.List[0].Names[0])
	}
	return a
}

func NodeAsString(fset *token.FileSet, v interface{}) string {
	var b bytes.Buffer
	printer.Fprint(&b, fset, v)
	br := b.String()
	return br
}

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
	Comment string
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
	str, err := WithTemplate(`{{.Comment}}
func {{.Recv}}{{.Name}}({{.Params}})({{.Returns}}){
	{{.Body}}
}
`, f)
	if err != nil {
		panic(err)
	}
	return formatCode(str)
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

