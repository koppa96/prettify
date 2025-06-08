package main

import (
	"os"

	"github.com/koppa96/prettify/config"
	"github.com/koppa96/prettify/doc"
)

func main() {
	config := config.Config{
		PrintWidth: 80,
		TabWidth:   4,
	}

	document := doc.Doc{
		Node: doc.Group{
			Node: doc.Concat{
				doc.Text("func foo("),
				doc.Group{
					Node: doc.Concat{
						doc.Indent{
							Node: doc.Concat{
								doc.SoftLine{},
								doc.Join{
									Sep: doc.Line{},
									Nodes: []doc.Node{
										doc.Text("arg1 string,"),
										doc.Text("arg2 int,"),
										doc.Concat{
											doc.Text("arg3 float64"),
											doc.SoftComma{},
										},
									},
								},
							},
						},
						doc.SoftLine{},
					},
				},
				doc.Text(") error"),
			},
		},
	}

	err := document.Render(config, os.Stdout)
	if err != nil {
		panic(err)
	}
}
