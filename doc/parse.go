package doc

import (
	"fmt"
	"github.com/dave/dst"
	"go/token"
	"slices"
	"strings"
)

// Parse parses a [*dst.File] and transforms it into a Doc that can later be rendered.
func Parse(file *dst.File) (*Doc, error) {
	root := Concat{
		Textf("package %s", file.Name.Name),
		HardLine{},
	}

	for _, decl := range file.Decls {
		node, err := parseDecl(decl)
		if err != nil {
			return nil, err
		}

		root = append(root, node, HardLine{})
	}

	return &Doc{
		Node: root,
	}, nil
}

func parseDecl(decl dst.Decl) (Node, error) {
	switch d := decl.(type) {
	case *dst.GenDecl:
		return parseGenDecl(d)
	case *dst.FuncDecl:
		return parseFuncDecl(d)
	default:
		return nil, fmt.Errorf("unknown declaration type: %T", decl)
	}
}

func parseGenDecl(decl *dst.GenDecl) (Node, error) {
	switch decl.Tok {
	case token.IMPORT:
		return parseImportDecl(decl)
	case token.TYPE:
		return parseTypeDecl(decl)
	case token.CONST:
		return parseConstDecl(decl)
	case token.VAR:
		return parseVarDecl(decl)
	default:
		return nil, fmt.Errorf("unknown generic declaration token: %s", decl.Tok)
	}
}

func parseImportDecl(decl *dst.GenDecl) Node {
	if len(decl.Specs) == 1 {
		return Concat{
			Text("import "),
			parseImportSpec(decl.Specs[0].(*dst.ImportSpec)),
		}
	}

	std, ext := sortImportSpecs(decl)
	specs := make(Concat, 0, 2*len(decl.Specs))
	for _, spec := range std {
		specs = append(specs, HardLine{}, parseImportSpec(spec))
	}

	for i, spec := range ext {
		var prefix Node = HardLine{}
		if i == 0 && len(std) > 0 {
			prefix = DoubleLine{}
		}

		specs = append(specs, prefix, parseImportSpec(spec))
	}

	return Concat{
		Text("import ("),
		Indent{Node: specs},
		HardLine{},
		Text(")"),
	}
}

func sortImportSpecs(decl *dst.GenDecl) (std []*dst.ImportSpec, ext []*dst.ImportSpec) {
	for _, spec := range decl.Specs {
		importSpec := spec.(*dst.ImportSpec)
		if strings.ContainsRune(importSpec.Path.Value, '.') {
			ext = append(ext, importSpec)
		} else {
			std = append(std, importSpec)
		}
	}

	slices.SortFunc(std, func(a, b *dst.ImportSpec) int {
		return strings.Compare(a.Path.Value, b.Path.Value)
	})

	slices.SortFunc(ext, func(a, b *dst.ImportSpec) int {
		return strings.Compare(a.Path.Value, b.Path.Value)
	})

	return std, ext
}

func parseImportSpec(spec *dst.ImportSpec) Node {
	if spec.Name != nil {
		return Textf("%s %s", spec.Name.Name, spec.Path.Value)
	}

	return Text(spec.Path.Value)
}

func parseTypeDecl(decl *dst.GenDecl) (Node, error) {

}

func parseConstDecl(decl *dst.GenDecl) (Node, error) {

}

func parseVarDecl(decl *dst.GenDecl) (Node, error) {

}

func parseFuncDecl(decl *dst.FuncDecl) (Node, error) {

}
