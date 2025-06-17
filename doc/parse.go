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

		root = append(root, HardLine{}, node, HardLine{})
	}

	return &Doc{root}, nil
}

func parseDecl(decl dst.Decl) (Node, error) {
	switch d := decl.(type) {
	case *dst.GenDecl:
		return parseGenDecl(d)
	case *dst.FuncDecl:
		return parseFuncDecl(d), nil
	default:
		return nil, fmt.Errorf("unknown declaration type: %T", decl)
	}
}

func parseGenDecl(decl *dst.GenDecl) (Node, error) {
	switch decl.Tok {
	case token.IMPORT:
		return parseImportDecl(decl), nil
	case token.TYPE:
		return parseTypeDecl(decl), nil
	case token.CONST:
		return parseConstDecl(decl), nil
	case token.VAR:
		return parseVarDecl(decl), nil
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
		Indent{specs},
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

func parseTypeDecl(decl *dst.GenDecl) Node {
	if len(decl.Specs) == 1 {
		return Group{
			Concat{
				Text("type "),
				parseTypeSpec(decl.Specs[0].(*dst.TypeSpec)),
			},
		}
	}

	specs := make(Concat, 0, len(decl.Specs)*2)
	for _, spec := range decl.Specs {
		specs = append(specs, HardLine{}, Group{parseTypeSpec(spec.(*dst.TypeSpec))})
	}

	return Concat{
		Text("type ("),
		Indent{specs},
		HardLine{},
		Text(")"),
	}
}

func parseTypeSpec(spec *dst.TypeSpec) Node {
	node := Concat{Text(spec.Name.Name)}
	if spec.TypeParams != nil {
		node = append(node, Group{
			Concat{
				Text("["),
				parseParamList(spec.TypeParams.List),
				Text("]"),
			},
		})
	}

	return append(node, Space{}, parseExpr(spec.Type))
}

func parseParamList(list []*dst.Field) Node {
	params := make([]Node, 0, len(list))
	for _, param := range list {
		params = append(params, parseParam(param))
	}

	return Concat{
		Indent{
			Concat{
				SoftLine{},
				Join{
					Sep: Concat{
						Comma{},
						Line{},
					},
					Nodes: params,
				},
				SoftComma{},
			},
		},
		SoftLine{},
	}
}

func parseParam(param *dst.Field) Node {
	if len(param.Names) == 0 {
		return parseExpr(param.Type)
	}

	names := make([]Node, 0, len(param.Names))
	for _, name := range param.Names {
		names = append(names, Text(name.Name))
	}

	return Concat{
		Join{
			Sep: Concat{
				Comma{},
				Space{},
			},
			Nodes: names,
		},
		Space{},
		parseExpr(param.Type),
	}
}

func parseExpr(expr dst.Expr) Node {
	switch e := expr.(type) {
	case *dst.Ident:
		return parseIdent(e)
	case *dst.InterfaceType:
		return parseInterfaceType(e)
	}

	return nil
}

func parseIdent(i *dst.Ident) Node {
	if i.Path != "" {
		return Textf("%s.%s", i.Path, i.Name)
	}

	return Text(i.Name)
}

func parseInterfaceType(i *dst.InterfaceType) Node {
	if i.Methods == nil {
		return Text("interface{}")
	}

	methods := make([]Node, 0, len(i.Methods.List))
	for _, method := range i.Methods.List {
		methods = append(methods, parseInterfaceMethod(method))
	}

	return Expand{
		Concat{
			Text("interface {"),
			Indent{
				Concat{
					HardLine{},
					Join{
						Sep:   HardLine{},
						Nodes: methods,
					},
				},
			},
			HardLine{},
			Text("}"),
		},
	}
}

func parseInterfaceMethod(method *dst.Field) Node {
	t := method.Type.(*dst.FuncType)

	node := Concat{Textf("%s(", method.Names[0].Name)}
	if t.Params != nil {
		node = append(node, parseParamList(t.Params.List))
	}
	node = append(node, Text(")"))
	if t.Results != nil {
		if len(t.Results.List) == 1 && len(t.Results.List[0].Names) == 0 {
			node = append(node, Space{}, parseExpr(t.Results.List[0].Type))
		} else {
			node = append(node, Group{
				Concat{
					Text(" ("),
					parseParamList(t.Results.List),
					Text(")"),
				},
			})
		}
	}

	return Group{node}
}

func parseConstDecl(decl *dst.GenDecl) Node {
	return nil
}

func parseVarDecl(decl *dst.GenDecl) Node {
	return nil
}

func parseFuncDecl(decl *dst.FuncDecl) Node {
	return nil
}
