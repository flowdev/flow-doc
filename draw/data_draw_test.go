package draw_test

import "github.com/flowdev/ea-flow-doc/draw"

var bigMerge = &draw.Merge{
	ID:   "bigMerge",
	Size: 3,
}

var BigTestFlowData = &draw.Flow{
	Name: "bigTestFlow",
	AllShapes: &draw.Split{
		Shapes: [][]any{
			{
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
					},
					HasSrcComp: false, SrcPort: "in",
					HasDstComp: true, DstPort: "",
				},
				&draw.Comp{
					Main: &draw.DataType{
						Name: "Xa", Type: "MiSo", Link: "https://google.com?q=MiSo",
					},
				},
				&draw.Split{
					Shapes: [][]any{
						{
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
								},
								HasSrcComp: true, SrcPort: "special",
								HasDstComp: true, DstPort: "in",
							},
							&draw.Comp{
								Main: &draw.DataType{
									Type: "To", Link: "https://google.com?q=To",
								},
								Plugins: []*draw.PluginGroup{
									{
										Title: "semantics",
										Types: []*draw.Plugin{
											{Type: "TextSemantics", Link: "https://google.com?q=TextSemantics"},
										},
									},
									{
										Title: "subParser",
										Types: []*draw.Plugin{
											{Type: "LiteralParser", Link: "https://google.com?q=LiteralParser", GoLink: true},
											{Type: "NaturalParser", Link: "https://google.com?q=NaturalParser"},
										},
									},
								},
							},
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "bigData", Type: "BigDataType", Link: "https://google.com?q=BigDataType"},
								},
								HasSrcComp: true, SrcPort: "out",
								HasDstComp: true, DstPort: "in1",
							},
							bigMerge,
						}, {
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
								},
								HasSrcComp: true, SrcPort: "out",
								HasDstComp: true, DstPort: "in",
							},
							&draw.Comp{
								Main: &draw.DataType{
									Name: "Mla", Type: "Blue", Link: "https://google.com?q=Blue",
								},
							},
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "data2", Type: "Data2", Link: "https://google.com?q=Data2"},
								},
								HasSrcComp: true, SrcPort: "",
								HasDstComp: false, DstPort: "",
							},
							&draw.Sequel{
								Number: 1,
							},
						},
					},
				},
			}, {
				&draw.Sequel{
					Number: 1,
				},
				&draw.Arrow{
					DataTypes: []*draw.DataType{},
					HasSrcComp:  false, SrcPort: "",
					HasDstComp: true, DstPort: "in",
				},
				&draw.Comp{
					Main: &draw.DataType{
						Name: "bla2", Type: "Blue", Link: "https://google.com?q=Blue",
					},
				},
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
					},
					HasSrcComp: true, SrcPort: "out",
					HasDstComp: true, DstPort: "in2",
				},
				bigMerge,
			}, {
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data3", Type: "Data3", Link: "https://google.com?q=Data3"},
					},
					HasSrcComp: false, SrcPort: "in2",
					HasDstComp: true, DstPort: "in",
				},
				&draw.Comp{
					Main: &draw.DataType{
						Name: "megaParser", Type: "MegaParser", Link: "https://google.com?q=MegaParser",
					},
					Plugins: []*draw.PluginGroup{
						{
							Title: "semantics",
							Types: []*draw.Plugin{
								{Type: "TextSemantics", Link: "https://google.com?q=TextSemantics"},
							},
						},
						{
							Title: "subParser",
							Types: []*draw.Plugin{
								{Type: "LiteralParser", Link: "https://google.com?q=LiteralParser", GoLink: true},
								{Type: "NaturalParser", Link: "https://google.com?q=NaturalParser"},
							},
						},
					},
				},
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
						{Name: "data2", Type: "data2", Link: "https://google.com?q=data2"},
						{Name: "data3", Type: "Data3", Link: "https://google.com?q=Data3"},
					},
					HasSrcComp: true, SrcPort: "out",
					HasDstComp: true, DstPort: "in3",
				},
				bigMerge,
			}, {
				&draw.Comp{
					Main: &draw.DataType{
						Type: "bigMerge", Link: "https://google.com?q=bigMerge",
					},
				},
			}, { // empty to force more space
			}, {
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
						{Name: "data2", Type: "data2", Link: "https://google.com?q=data2"},
						{Name: "data3", Type: "Data3", Link: "https://google.com?q=Data3"},
					},
					HasSrcComp: false, SrcPort: "in3",
					HasDstComp: true, DstPort: "",
				},
				&draw.Comp{
					Main: &draw.DataType{
						Type: "recursive", Link: "https://google.com?q=recursive",
					},
				},
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
					},
					HasSrcComp: true, SrcPort: "",
					HasDstComp: true, DstPort: "",
				},
				&draw.Comp{
					Main: &draw.DataType{
						Type: "secondOp", Link: "https://google.com?q=secondOp",
					},
				},
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
						{Name: "data2", Type: "data2", Link: "https://google.com?q=data2"},
						{Name: "data3", Type: "Data3", Link: "https://google.com?q=Data3"},
					},
					HasSrcComp: true, SrcPort: "out",
					HasDstComp: false, DstPort: "",
				},
				&draw.Loop{
					Name: "recursive", Port: "in3", Link: "https://google.com?q=recursive",
				},
			},
		},
	},
}
