package generator

import (
	"bytes"
	"go/scanner"
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

	src, err := ioutil.ReadFile(path.Join(dir, filename))
	if err != nil {
		return err
	}
	constNames, err := getPrefixedNames(src, typeName)
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

func getPrefixedNames(src []byte, prefix string) ([]string, error) {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	constStmt := false
	inParen := false
	global := true
	var names []string
	for {
		pos, tok, lit := s.Scan()
		if s.ErrorCount > 0 {
			return nil, errors.Errorf("invalid go file given: invalid token found at %v", pos)
		}

		if tok == token.EOF {
			break
		}

		switch tok {
		case token.CONST:
			if global {
				constStmt = true
			}
		case token.LPAREN:
			if constStmt {
				inParen = true
			}
		case token.RPAREN:
			if constStmt && inParen {
				constStmt = false
				inParen = false
			}
		case token.IDENT:
			if constStmt {
				if strings.HasPrefix(lit, prefix) {
					names = append(names, lit)
				}
				if !inParen {
					constStmt = false
				}
			}
		case token.FUNC:
			global = false
		case token.RBRACE:
			if !global {
				global = true
			}
		}
	}
	return names, nil
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
