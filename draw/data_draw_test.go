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
					HasSrcOp: false, SrcPort: "in",
					HasDstOp: true, DstPort: "",
				},
				&draw.Op{
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
								HasSrcOp: true, SrcPort: "special",
								HasDstOp: true, DstPort: "in",
							},
							&draw.Op{
								Main: &draw.DataType{
									Type: "To", Link: "https://google.com?q=MiSo",
								},
								Plugins: []*draw.Plugin{
									{
										Title: "semantics",
										Types: []*draw.PluginType{
											{Type: "TextSemantics", Link: "https://google.com?q=TextSemantics"},
										},
									},
									{
										Title: "subParser",
										Types: []*draw.PluginType{
											{Type: "LiteralParser", Link: "https://google.com?q=LiteralParser"},
											{Type: "NaturalParser", Link: "https://google.com?q=NaturalParser"},
										},
									},
								},
							},
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "bigData", Type: "BigDataType", Link: "https://google.com?q=MiSo"},
								},
								HasSrcOp: true, SrcPort: "out",
								HasDstOp: true, DstPort: "in1",
							},
							bigMerge,
						}, {
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
								},
								HasSrcOp: true, SrcPort: "out",
								HasDstOp: true, DstPort: "in",
							},
							&draw.Op{
								Main: &draw.DataType{
									Name: "Mla", Type: "Blue", Link: "https://google.com?q=Blue",
								},
							},
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "data2", Type: "Data2", Link: "https://google.com?q=Data2"},
								},
								HasSrcOp: true, SrcPort: "",
								HasDstOp: false, DstPort: "...",
							},
						},
					},
				},
			}, {
				&draw.Arrow{
					DataTypes: []*draw.DataType{},
					HasSrcOp:  false, SrcPort: "...",
					HasDstOp: true, DstPort: "in",
				},
				&draw.Op{
					Main: &draw.DataType{
						Name: "bla2", Type: "Blue", Link: "https://google.com?q=Blue",
					},
				},
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
					},
					HasSrcOp: true, SrcPort: "out",
					HasDstOp: true, DstPort: "in2",
				},
				bigMerge,
			}, {
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data3", Type: "Data3", Link: "https://google.com?q=Data3"},
					},
					HasSrcOp: false, SrcPort: "in2",
					HasDstOp: true, DstPort: "in",
				},
				&draw.Op{
					Main: &draw.DataType{
						Name: "megaParser", Type: "MegaParser", Link: "https://google.com?q=MegaParser",
					},
					Plugins: []*draw.Plugin{
						{
							Title: "semantics",
							Types: []*draw.PluginType{
								{Type: "TextSemantics", Link: "https://google.com?q=TextSemantics"},
							},
						},
						{
							Title: "subParser",
							Types: []*draw.PluginType{
								{Type: "LiteralParser", Link: "https://google.com?q=LiteralParser"},
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
					HasSrcOp: true, SrcPort: "out",
					HasDstOp: true, DstPort: "in3",
				},
				bigMerge,
			}, {
				&draw.Op{
					Main: &draw.DataType{
						Type: "bigMerge", Link: "https://google.com?q=Blue",
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
					HasSrcOp: false, SrcPort: "in3",
					HasDstOp: true, DstPort: "",
				},
				&draw.Op{
					Main: &draw.DataType{
						Type: "recursive", Link: "https://google.com?q=recursive",
					},
				},
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
					},
					HasSrcOp: true, SrcPort: "",
					HasDstOp: true, DstPort: "",
				},
				&draw.Op{
					Main: &draw.DataType{
						Type: "secondOp", Link: "https://google.com?q=recursive",
					},
				},
				&draw.Arrow{
					DataTypes: []*draw.DataType{{Name: "data", Type: "Data"}, {Name: "data2", Type: "data2"}, {Name: "data3", Type: "Data3"}},
					HasSrcOp:  true, SrcPort: "out",
					HasDstOp: false, DstPort: "...back to: recursive",
				},
			},
		},
	},
}
