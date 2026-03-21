package doc

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"strings"
)

// Parse parses a [*ast.File] and transforms it into a Doc that can later be rendered.
func Parse(file *ast.File) (*Doc, error) {
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

func parseDecl(decl ast.Decl) (Node, error) {
	switch d := decl.(type) {
	case *ast.GenDecl:
		return parseGenDecl(d)
	case *ast.FuncDecl:
		return parseFuncDecl(d), nil
	default:
		return nil, fmt.Errorf("unknown declaration type: %T", decl)
	}
}

func parseGenDecl(decl *ast.GenDecl) (Node, error) {
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

func parseImportDecl(decl *ast.GenDecl) Node {
	if len(decl.Specs) == 1 {
		return Concat{
			Text("import "),
			parseImportSpec(decl.Specs[0].(*ast.ImportSpec)),
		}
	}

	std, ext := sortImportSpecs(decl)

	stdNodes := make([]Node, 0, len(std))
	for _, spec := range std {
		stdNodes = append(stdNodes, parseImportSpec(spec))
	}

	extNodes := make([]Node, 0, len(ext))
	for _, spec := range ext {
		extNodes = append(extNodes, parseImportSpec(spec))
	}

	var blocks []Node
	if len(stdNodes) > 0 {
		blocks = append(blocks, Join{
			Sep:   HardLine{},
			Nodes: stdNodes,
		})
	}

	if len(extNodes) > 0 {
		blocks = append(blocks, Join{
			Sep:   HardLine{},
			Nodes: extNodes,
		})
	}

	return Concat{
		Text("import ("),
		Indent{
			Concat{
				HardLine{},
				Join{
					Sep:   DoubleLine{},
					Nodes: blocks,
				},
			},
		},
		HardLine{},
		Text(")"),
	}
}

func sortImportSpecs(decl *ast.GenDecl) (std []*ast.ImportSpec, ext []*ast.ImportSpec) {
	for _, spec := range decl.Specs {
		importSpec := spec.(*ast.ImportSpec)
		if strings.ContainsRune(importSpec.Path.Value, '.') {
			ext = append(ext, importSpec)
		} else {
			std = append(std, importSpec)
		}
	}

	slices.SortFunc(std, func(a, b *ast.ImportSpec) int {
		return strings.Compare(a.Path.Value, b.Path.Value)
	})

	slices.SortFunc(ext, func(a, b *ast.ImportSpec) int {
		return strings.Compare(a.Path.Value, b.Path.Value)
	})

	return std, ext
}

func parseImportSpec(spec *ast.ImportSpec) Node {
	if spec.Name != nil {
		return Textf("%s %s", spec.Name.Name, spec.Path.Value)
	}

	return Text(spec.Path.Value)
}

func parseTypeDecl(decl *ast.GenDecl) Node {
	if len(decl.Specs) == 1 {
		return Group{
			Concat{
				Text("type "),
				parseTypeSpec(decl.Specs[0].(*ast.TypeSpec)),
			},
		}
	}

	specs := make(Concat, 0, len(decl.Specs)*2)
	for _, spec := range decl.Specs {
		specs = append(specs, HardLine{}, Group{parseTypeSpec(spec.(*ast.TypeSpec))})
	}

	return Concat{
		Text("type ("),
		Indent{specs},
		HardLine{},
		Text(")"),
	}
}

func parseTypeSpec(spec *ast.TypeSpec) Node {
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

func parseParamList(list []*ast.Field) Node {
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

func parseParam(param *ast.Field) Node {
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

func parseExpr(expr ast.Expr) Node {
	switch e := expr.(type) {
	case *ast.Ident:
		return Text(e.Name)
	case *ast.InterfaceType:
		return parseInterfaceType(e)
	case *ast.StructType:
		return parseStructType(e)
	case *ast.FuncType:
		return parseFuncType(e)
	case *ast.StarExpr:
		return Concat{Text("*"), parseExpr(e.X)}
	case *ast.SelectorExpr:
		return parseSelectorExpr(e)
	}

	return nil
}

func parseSelectorExpr(s *ast.SelectorExpr) Node {
	return Concat{parseExpr(s.X), Text("." + s.Sel.Name)}
}

func parseInterfaceType(i *ast.InterfaceType) Node {
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

func parseInterfaceMethod(method *ast.Field) Node {
	t := method.Type.(*ast.FuncType)

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

func parseStructType(s *ast.StructType) Node {
	return nil
}

func parseFuncType(f *ast.FuncType) Node {
	return parseInterfaceMethod(&ast.Field{
		Names: []*ast.Ident{
			{Name: "func"},
		},
		Type: f,
	})
}

func parseConstDecl(decl *ast.GenDecl) Node {
	return nil
}

func parseVarDecl(decl *ast.GenDecl) Node {
	return nil
}

func parseFuncDecl(decl *ast.FuncDecl) Node {
	return nil
}
