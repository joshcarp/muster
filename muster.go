package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/tools/go/packages"
)

var (
	stream = flag.Bool("stream", false, "use stdin/stdout instead of writing to file")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("stringer: ")
	flag.Usage = Usage
	flag.Parse()
	var tags []string

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}
	if *stream {
		methodStream(os.Stdin, os.Stdout)
	} else {
		methodFiles(tags, args)
	}
}

func methodFiles(tags []string, args []string) {
	pkgs := parseFiles(tags, args)
	// Run generate for each file.
	for i, file := range pkgs[0].Syntax {
		filename := path.Base(pkgs[0].GoFiles[i])
		if strings.Contains(filename, "_must.go") {
			continue
		}
		contents, err := generate(pkgs[0].Fset, file)
		if err != nil {
			log.Print(err)
			continue
		}
		os.WriteFile(strings.ReplaceAll(filename, ".go", "_must.go"), []byte(contents), 0644)
	}
}

func methodStream(reader io.Reader, writer io.Writer) {
	stdin, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	fset, file, err := parseStrings(string(stdin))
	if err != nil {
		log.Fatal(err)
	}
	contents, err := generate(fset, file)
	writer.Write([]byte(contents))
}

func parseStrings(content string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(token.NewFileSet(), "", []byte(content), 0)
	if err != nil {
		return nil, nil, err
	}
	return fset, file, nil

}

func parseFiles(tags []string, args []string) []*packages.Package {
	cfg := &packages.Config{
		Mode:       packages.LoadSyntax,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, args...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	return pkgs
}

// generate produces the String method for the named type.
func generate(fset *token.FileSet, file *ast.File) (string, error) {
	tmpl := `// Code generated by muster. DO NOT EDIT.
// see https://github.com/joshcarp/muster for more details.
package {{file.Name.Name}}
import (
{{range $j, $import := file.Imports}}
{{NodeAsString $import}}{{end}})

{{Generate fset file }}
`
	result, err := WithTemplate(tmpl, nil,
		"fset", func() interface{} { return fset },
		"file", func() interface{} { return file },
		"Generate", Generate,
		"NodeAsString", NodeAsStringFunc(fset),
	)
	if err != nil {
		return "", err
	}
	return importCode(result)
}

func Generate(fset *token.FileSet, file *ast.File) string {
	var b bytes.Buffer
	ast.Inspect(file, genDecl(&b, fset))
	return b.String()
}

// genDecl processes one declaration clause.
func genDecl(writer io.Writer, fset *token.FileSet) func(ast.Node) bool {
	return func(node ast.Node) bool {
		decl, ok := node.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if decl.Type == nil {
			return true
		}
		if decl.Type.Results == nil {
			return true
		}
		if len(decl.Type.Results.List) != 2 {
			return true
		}
		if NodeAsString(fset, decl.Type.Results.List[1].Type) != "error" {
			return true
		}
		call := FunctionFromDecl(decl, fset)
		must := call
		var err error
		must.Body, err = WithTemplate(
			`val, err := {{if ne .Recv.Name ""}}{{.Recv.Name}}.{{end}}{{.Name}}({{ range $i, $e := .Params }}{{$e.Name}}, {{end}})
	if err != nil {
		panic(err)
	}
	return val`, must)
		if err != nil {
			log.Printf("error printing function %s", must.Name)
		}
		switch must.IsExported() {
		case true:
			must.Name = "Must" + must.Name
		case false:
			must.Name = "must" + strcase.ToCamel(must.Name)
		}
		must.Returns = must.Returns[:len(must.Returns)-1]
		must.Comment = fmt.Sprintf("// %s calls %s and panics if err is not nil.", must.Name, call.Name)
		writer.Write([]byte(must.String()))
		return false
	}
}
