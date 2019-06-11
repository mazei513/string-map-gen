package generator

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

type templateData struct {
	Package  string
	TypeName string
	Name     string
	Items    []Item
}

type Item struct {
	Key   string
	Value string
}

func Generate(typeName, dir, filename string) error {
	if typeName == "" || dir == "" {
		return errors.New("invalid arguments")
	}

	pkg, err := getPackage(dir)
	if err != nil {
		return err
	}

	fileBytes, err := ioutil.ReadFile(path.Join(dir, filename))
	if err != nil {
		return err
	}
	constNames, err := getPrefixedNames(string(fileBytes), typeName)
	if err != nil {
		return err
	}
	items := getItemsFromNames(constNames, typeName)

	d := templateData{
		Package:  pkg.Name,
		TypeName: typeName,
		Name:     getNameFromObjType(typeName),
		Items:    items,
	}

	path := getFilepath(dir, d.Name)

	b, err := generateFileBytes(path, d)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, b, 0644); err != nil {
		return err
	}

	return nil
}

func getPackage(dir string) (*packages.Package, error) {
	p, err := packages.Load(&packages.Config{
		Dir: dir,
	}, ".")
	if err != nil {
		return nil, errors.Wrap(err, "failed to load package info")
	}

	if len(p) != 1 {
		return nil, errors.Errorf("invalid directory given (%s)", dir)
	}

	return p[0], nil
}

func getNameFromObjType(s string) string {
	s = strings.TrimLeft(s, "*")
	return upperFirst(s)
}

func getFilepath(dir, name string) string {
	fn := strings.ToLower(name) + "_stringmap.go"
	return filepath.Join(dir, fn)
}

func generateFileBytes(filename string, data templateData) ([]byte, error) {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}

	src, err := imports.Process(filename, buf.Bytes(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to format file")
	}
	return src, nil
}

func getPrefixedNames(file string, prefix string) ([]string, error) {
	f, err := parseFile(file)
	if err != nil {
		return nil, err
	}

	a := getNames(f)
	var res []string

	for _, s := range a {
		if s == prefix {
			continue
		}

		if strings.HasPrefix(s, prefix) {
			res = append(res, s)
		}
	}

	return res, nil
}

func parseFile(file string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, "stringmap.go", file, 0)
}

func getNames(f ast.Node) []string {
	var names []string

	ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.Ident:
			s = x.Name
		}
		if s != "" {
			names = append(names, s)
		}
		return true
	})

	return names
}

func getItemsFromNames(names []string, prefix string) []Item {
	items := make([]Item, len(names))
	for i, n := range names {
		trimmed := strings.TrimLeft(strings.TrimPrefix(n, prefix), "_")
		items[i] = Item{
			Key:   upperFirst(trimmed),
			Value: n,
		}
	}
	return items
}

func upperFirst(s string) string {
	if len(s) == 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
