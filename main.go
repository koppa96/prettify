package main

import (
	"github.com/dave/dst/decorator"
	"os"

	"github.com/koppa96/prettify/config"
	"github.com/koppa96/prettify/doc"
)

func main() {
	cfg := config.Config{
		PrintWidth: 80,
		TabWidth:   4,
	}

	f, err := decorator.Parse(`package foo

import (
	"os"
	"fmt"
	_ "github.com/lib/pq"
	"io"
)

type Foo[T any] interface {
	Bar(param1 string, param2 T) (string, error)
}`)
	if err != nil {
		panic(err)
	}

	document, err := doc.Parse(f)
	if err != nil {
		panic(err)
	}

	//document := doc.Doc{
	//	Node: doc.Concat{
	//		doc.Text("package test"),
	//		doc.DoubleLine{},
	//		doc.Text("import ("),
	//		doc.Indent{
	//			Node: doc.Concat{
	//				doc.HardLine{},
	//				doc.Text("\"fmt\""),
	//				doc.HardLine{},
	//				doc.Text("\"os\""),
	//			},
	//		},
	//		doc.HardLine{},
	//		doc.Text(")"),
	//		doc.DoubleLine{},
	//		doc.Group{
	//			Node: doc.Concat{
	//				doc.Text("func foo("),
	//				doc.Indent{
	//					Node: doc.Concat{
	//						doc.SoftLine{},
	//						doc.Join{
	//							Sep: doc.Line{},
	//							Nodes: []doc.Node{
	//								doc.Text("arg1 string,"),
	//								doc.Text("arg2 int,"),
	//								doc.Text("arg3 float64"),
	//							},
	//						},
	//						doc.SoftComma{},
	//					},
	//				},
	//				doc.SoftLine{},
	//				doc.Text(")"),
	//				doc.Group{
	//					Node: doc.Concat{
	//						doc.Text(" ("),
	//						doc.Indent{
	//							Node: doc.Concat{
	//								doc.SoftLine{},
	//								doc.Join{
	//									Sep: doc.Line{},
	//									Nodes: []doc.Node{
	//										doc.Text("string,"),
	//										doc.Concat{
	//											doc.Text("error"),
	//											doc.SoftComma{},
	//										},
	//									},
	//								},
	//							},
	//						},
	//						doc.SoftLine{},
	//						doc.Text(") {"),
	//					},
	//				},
	//			},
	//		},
	//		doc.Indent{
	//			Node: doc.Concat{
	//				doc.HardLine{},
	//				doc.Text("err := fmt.Fprintf("),
	//				doc.Group{
	//					Node: doc.Concat{
	//						doc.Indent{
	//							Node: doc.Concat{
	//								doc.SoftLine{},
	//								doc.Join{
	//									Sep: doc.Line{},
	//									Nodes: []doc.Node{
	//										doc.Text("\"%s %d %f\","),
	//										doc.Text("arg1,"),
	//										doc.Text("arg2,"),
	//										doc.Concat{
	//											doc.Text("arg3"),
	//											doc.SoftComma{},
	//										},
	//									},
	//								},
	//							},
	//						},
	//						doc.SoftLine{},
	//					},
	//				},
	//				doc.Text(")"),
	//			},
	//		},
	//		doc.HardLine{},
	//		doc.Text("}"),
	//	},
	//}

	err = document.Render(cfg, os.Stdout)
	if err != nil {
		panic(err)
	}
}
