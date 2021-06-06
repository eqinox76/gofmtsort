package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"reflect"
	"sort"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/pkg/errors"
)

func Sort(source []byte) (string, error) {
	root, err := decorator.Parse(source)

	if err != nil {
		return "", errors.Wrap(err, "could not parse")
	}

	// for _, n := range root.Decls {
	// 	fmt.Printf("[%T] %+v\n", n, n)
	// }

	// fmt.Println("===============")
	sort.Slice(root.Decls, sortDecl(root.Decls))

	// for _, n := range root.Decls {
	// 	fmt.Printf("[%T] %+v\n", n, n)
	// }

	var buf bytes.Buffer
	err = decorator.Fprint(&buf, root)
	if err != nil {
		return "", errors.Wrap(err, "could not format output")
	}

	return buf.String(), nil
}

func sortDecl(a []dst.Decl) func(i, j int) bool {
	DeclPriority := map[reflect.Type]int{
		reflect.TypeOf(&dst.FuncDecl{}): 1,
	}

	return func(i, j int) bool {
		this := a[i]
		other := a[j]

		thisPrio, ok := DeclPriority[reflect.TypeOf(this)]
		if !ok {
			fmt.Printf("%T %t not implemented\n", this, this)
			return true
		}
		otherPrio, ok := DeclPriority[reflect.TypeOf(other)]
		if !ok {
			fmt.Printf("%T %t not implemented\n", other, other)
			return true
		}

		if thisPrio != otherPrio {
			return thisPrio < otherPrio
		}

		switch thisV := this.(type) {
		case *dst.FuncDecl:
			otherV := other.(*dst.FuncDecl)

			// public (exported) functions first
			if ast.IsExported(thisV.Name.Name) != ast.IsExported(otherV.Name.Name) {
				return ast.IsExported(thisV.Name.Name)
			}

			return thisV.Name.Name < otherV.Name.Name
		default:
			fmt.Printf("%T %t not implemented\n", other, other)
		}

		return true
	}
}
