package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"log"
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

	sort.SliceStable(root.Decls, sortDecl(root.Decls))

	var buf bytes.Buffer
	err = decorator.Fprint(&buf, root)
	if err != nil {
		return "", errors.Wrap(err, "could not format output")
	}

	return buf.String(), nil
}

func funcName(f *dst.FuncDecl) string {
	if f.Recv == nil {
		return f.Name.Name
	}
	recvType, ok := f.Recv.List[0].Type.(*dst.StarExpr)
	if !ok {
		log.Print(f.Recv.List[0].Type)
		return f.Name.Name
	}
	ident, ok := recvType.X.(*dst.Ident)
	if !ok {
		log.Print(recvType.X)
		return f.Name.Name
	}
	return ident.Name + f.Name.Name
}

func genName(g *dst.GenDecl) string {
	switch g.Tok {
	case token.CONST, token.VAR:
		spec, ok := g.Specs[0].(*dst.ValueSpec)
		if !ok {
			log.Print("genName", g.Specs[0])
			return ""
		}
		return spec.Names[0].Name
	case token.STRUCT, token.TYPE:
		spec, ok := g.Specs[0].(*dst.TypeSpec)
		if !ok {
			log.Print("genName", g.Specs[0])
			return ""
		}
		return spec.Name.Name
	default:
		log.Print("genName ", g.Tok, " not supported")
		return ""
	}
}

func genScore(g *dst.GenDecl) int {
	switch g.Tok {
	case token.CONST:
		return 2
	case token.IMPORT:
		return 1
	case token.VAR:
		return 3
	case token.STRUCT, token.FUNC, token.TYPE:
		return 5
	default:
		return 4
	}
}

func lessFuncFunc(this, other *dst.FuncDecl) bool {
	return lessName(funcName(this), funcName(other))
}

func lessFuncGen(this *dst.FuncDecl, other *dst.GenDecl) bool {
	// only structs/types are mixed with functions
	if other.Tok != token.STRUCT && other.Tok != token.TYPE {
		return false
	}
	// free functions are always later/bigger
	if this.Recv == nil || len(this.Recv.List) == 0 {
		return false
	}

	return lessName(funcName(this), genName(other))
}

func lessGen(this, other *dst.GenDecl) bool {
	if genScore(this) != genScore(other) {
		return genScore(this) < genScore(other)
	}
	return lessName(genName(this), genName(other))
}

func lessName(this, other string) bool {
	// public (exported) nodes first
	if ast.IsExported(this) != ast.IsExported(other) {
		return ast.IsExported(this)
	}

	return this < other
}

/*
imports

types and methods

functions
*/
func sortDecl(a []dst.Decl) func(i, j int) bool {
	DeclPriority := map[reflect.Type]int{
		reflect.TypeOf(&dst.FuncDecl{}): 1,
		reflect.TypeOf(&dst.GenDecl{}):  1,
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
			switch otherV := other.(type) {
			case *dst.FuncDecl:
				return lessFuncFunc(thisV, otherV)
			case *dst.GenDecl:
				return lessFuncGen(thisV, otherV)
			default:
				fmt.Printf("%T %t other not implemented\n", other, other)
			}

		case *dst.GenDecl:
			switch otherV := other.(type) {
			case *dst.FuncDecl:
				return !lessFuncGen(otherV, thisV)
			case *dst.GenDecl:
				return lessGen(thisV, otherV)
			default:
				fmt.Printf("%T %t other not implemented\n", other, other)
			}
		default:
			fmt.Printf("%T %t this not implemented\n", other, other)
		}

		return true
	}
}
