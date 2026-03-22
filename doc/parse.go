package doc

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"slices"
	"strings"
)

// Parse parses a [*ast.File] and transforms it into a Doc that can later be rendered.
func Parse(file *ast.File) (*Doc, error) {
	root := []Node{
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

	return &Doc{Concat(root...)}, nil
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
		return Concat(
			Text("import "),
			parseImportSpec(decl.Specs[0].(*ast.ImportSpec)),
		)
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
		blocks = append(blocks, Join(stdNodes, HardLine{}))
	}

	if len(extNodes) > 0 {
		blocks = append(blocks, Join(extNodes, HardLine{}))
	}

	return Concat(
		Text("import ("),
		Indent{
			Concat(
				HardLine{},
				Join(blocks, DoubleLine{}),
			),
		},
		HardLine{},
		Text(")"),
	)
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
			Concat(
				Text("type "),
				parseTypeSpec(decl.Specs[0].(*ast.TypeSpec)),
			),
		}
	}

	specs := make([]Node, 0, len(decl.Specs)*2)
	for _, spec := range decl.Specs {
		specs = append(specs, HardLine{}, Group{parseTypeSpec(spec.(*ast.TypeSpec))})
	}

	return Concat(
		Text("type ("),
		Indent{Concat(specs...)},
		HardLine{},
		Text(")"),
	)
}

func parseTypeSpec(spec *ast.TypeSpec) Node {
	nodes := []Node{Text(spec.Name.Name)}
	if spec.TypeParams != nil {
		nodes = append(nodes, Group{
			Concat(
				Text("["),
				parseParamList(spec.TypeParams.List),
				Text("]"),
			),
		})
	}

	nodes = append(nodes, Space{}, parseExpr(spec.Type))
	return Concat(nodes...)
}

func parseParamList(list []*ast.Field) Node {
	params := make([]Node, 0, len(list))
	for _, param := range list {
		params = append(params, parseParam(param))
	}

	return Concat(
		Indent{
			Concat(
				SoftLine{},
				Join(params, Concat(Comma{}, Line{})),
				SoftComma{},
			),
		},
		SoftLine{},
	)
}

func parseParam(param *ast.Field) Node {
	if len(param.Names) == 0 {
		return parseExpr(param.Type)
	}

	names := make([]Node, 0, len(param.Names))
	for _, name := range param.Names {
		names = append(names, Text(name.Name))
	}

	return Concat(
		Join(names, Concat(Comma{}, Space{})),
		Space{},
		parseExpr(param.Type),
	)
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
		return Concat(Text("*"), parseExpr(e.X))
	case *ast.SelectorExpr:
		return parseSelectorExpr(e)
	case *ast.BasicLit:
		return Text(e.Value)
	case *ast.ArrayType:
		return parseArrayType(e)
	case *ast.CompositeLit:
		return parseCompositeLit(e)
	}

	panic("unknown expression type: " + reflect.TypeOf(expr).Elem().Name())
}

func parseCompositeLit(lit *ast.CompositeLit) Node {
	if len(lit.Elts) == 0 {
		return Concat(parseExpr(lit.Type), Text("{}"))
	}

	elts := make([]Node, 0, len(lit.Elts))
	for _, el := range lit.Elts {
		elts = append(elts, parseExpr(el))
	}

	return Concat(
		parseExpr(lit.Type),
		Group{
			Concat(
				Text("{"),
				Indent{
					Concat(
						SoftLine{},
						Join(elts, Concat(Comma{}, Line{})),
						SoftComma{},
					),
				},
				SoftLine{},
				Text("}"),
			),
		},
	)
}

func parseArrayType(a *ast.ArrayType) Node {
	if a.Len == nil {
		return Concat(Text("[]"), parseExpr(a.Elt))
	}

	return Concat(
		Text("["),
		parseExpr(a.Len),
		Text("]"),
		parseExpr(a.Elt),
	)
}

func parseSelectorExpr(s *ast.SelectorExpr) Node {
	return Concat(parseExpr(s.X), Text("."+s.Sel.Name))
}

func parseInterfaceType(i *ast.InterfaceType) Node {
	if i.Methods == nil {
		return Text("interface{}")
	}

	methods := make([]Node, 0, len(i.Methods.List))
	for _, method := range i.Methods.List {
		methods = append(methods, parseInterfaceMethod(method))
	}

	return Concat(
		Text("interface {"),
		Indent{
			Concat(
				HardLine{},
				Join(methods, HardLine{}),
			),
		},
		HardLine{},
		Text("}"),
	)
}

func parseInterfaceMethod(method *ast.Field) Node {
	return Group{
		Node: Concat(
			Text(method.Names[0].Name),
			parseSignature(method.Type.(*ast.FuncType)),
		),
	}
}

func parseSignature(t *ast.FuncType) Node {
	nodes := []Node{Text("(")}
	if t.Params != nil {
		nodes = append(nodes, parseParamList(t.Params.List))
	}
	nodes = append(nodes, Text(")"))
	if t.Results != nil {
		if len(t.Results.List) == 1 && len(t.Results.List[0].Names) == 0 {
			nodes = append(nodes, Space{}, parseExpr(t.Results.List[0].Type))
		} else {
			nodes = append(nodes, Group{
				Concat(
					Text(" ("),
					parseParamList(t.Results.List),
					Text(")"),
				),
			})
		}
	}

	return Concat(nodes...)
}

func parseStructType(s *ast.StructType) Node {
	return Text("struct{}")
}

func parseFuncType(f *ast.FuncType) Node {
	return Group{
		Node: Concat(
			Text("func"),
			parseSignature(f),
		),
	}
}

func parseConstDecl(decl *ast.GenDecl) Node {
	return nil
}

func parseVarDecl(decl *ast.GenDecl) Node {
	if len(decl.Specs) == 1 {
		return Concat(
			Text("var "),
			parseValueSpec(decl.Specs[0].(*ast.ValueSpec)),
		)
	}

	nodes := make([]Node, 0, len(decl.Specs))
	for _, spec := range decl.Specs {
		nodes = append(nodes, parseValueSpec(spec.(*ast.ValueSpec)))
	}

	return Concat(
		Text("var ("),
		Indent{
			Concat(
				HardLine{},
				Join(nodes, HardLine{}),
			),
		},
		HardLine{},
		Text(")"),
	)
}

func parseValueSpec(spec *ast.ValueSpec) Node {
	var names Node
	if len(spec.Names) == 0 {
		names = Text(spec.Names[0].Name)
	} else {
		nameNodes := make([]Node, 0, len(spec.Names))
		for _, name := range spec.Names {
			nameNodes = append(nameNodes, Text(name.Name))
		}

		names = Group{
			Indent{
				Join(nameNodes, Concat(Comma{}, Line{})),
			},
		}
	}

	nodes := []Node{names}
	if (spec.Type) != nil {
		nodes = append(nodes, Space{}, parseExpr(spec.Type))
	}

	if len(spec.Values) == 0 {
		return Concat(nodes...)
	}

	nodes = append(nodes, Text(" = "))

	var values Node
	if len(spec.Values) == 1 {
		values = parseExpr(spec.Values[0])
	} else {
		valueNodes := make([]Node, 0, len(spec.Values))
		for _, value := range spec.Values {
			valueNodes = append(valueNodes, parseExpr(value))
		}

		values = Group{
			Indent{
				Join(valueNodes, Concat(Comma{}, Line{})),
			},
		}
	}

	nodes = append(nodes, values)
	return Concat(nodes...)
}

func parseFuncDecl(decl *ast.FuncDecl) Node {
	nodes := []Node{Text("func ")}
	if decl.Recv != nil {
		nodes = append(nodes, Group{
			Concat(
				Text("("),
				Indent{
					Concat(
						SoftLine{},
						Textf("%s ", decl.Recv.List[0].Names[0].Name),
						parseExpr(decl.Recv.List[0].Type),
					),
				},
				SoftLine{},
				Text(") "),
			),
		})
	}

	nodes = append(
		nodes,
		Text(decl.Name.Name),
		parseSignature(decl.Type),
	)

	if decl.Body == nil || len(decl.Body.List) == 0 {
		nodes = append(nodes, Text(" {}"))
		return Concat(nodes...)
	}

	nodes = append(nodes, parseBlockStmt(decl.Body))

	return Concat(nodes...)
}

func parseBlockStmt(s *ast.BlockStmt) Node {
	stmts := make([]Node, 0, len(s.List))
	for _, stmt := range s.List {
		stmts = append(stmts, parseStmt(stmt))
	}

	return Concat(
		Text("{"),
		Indent{
			Concat(
				HardLine{},
				Join(stmts, HardLine{}),
			),
		},
		HardLine{},
		Text("}"),
	)
}

func parseStmt(stmt ast.Stmt) Node {
	switch s := stmt.(type) {
	case *ast.ReturnStmt:
		return parseReturnStmt(s)
	}

	panic("unknown statement type: " + reflect.TypeOf(stmt).Elem().Name())
}

func parseReturnStmt(s *ast.ReturnStmt) Node {
	if len(s.Results) == 0 {
		return Text("return")
	}

	if len(s.Results) == 1 {
		return Concat(
			Text("return "),
			parseExpr(s.Results[0]),
		)
	}

	results := make([]Node, 0, len(s.Results))
	for _, r := range s.Results {
		results = append(results, parseExpr(r))
	}

	return Concat(
		Text("return "),
		Group{
			Indent{
				Join(results, Concat(Comma{}, Line{})),
			},
		},
	)
}
