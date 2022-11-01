package draw_test

import "github.com/flowdev/ea-flow-doc/draw"

var bigMerge = draw.Merge{
	ID:   "bigMerge",
	Size: 3,
}

var BigTestFlowData = draw.Flow{
	Shapes: [][]interface{}{
		{
			draw.Arrow{
				DataTypes: []draw.DataType{{Name: "data", Type: "Data"}},
				HasSrcOp:  false, SrcPort: "in",
				HasDstOp: true, DstPort: "",
			},
			draw.Op{
				Main: draw.DataType{Name: "Xa", Type: "MiSo"},
			},
			draw.Split{
				Shapes: [][]interface{}{
					{
						draw.Arrow{
							DataTypes: []draw.DataType{{Name: "data", Type: "Data"}},
							HasSrcOp:  true, SrcPort: "special",
							HasDstOp: true, DstPort: "in",
						},
						draw.Op{
							Main: draw.DataType{Type: "To"},
							Plugins: []draw.Plugin{
								{
									Title: "semantics",
									Types: []draw.PluginType{
										{Type: "TextSemantics"},
									},
								},
								{
									Title: "subParser",
									Types: []draw.PluginType{
										{Type: "LitralParser"},
										{Type: "NaturalParser"},
									},
								},
							},
						},
						draw.Arrow{
							DataTypes: []draw.DataType{{Name: "bigData", Type: "BigDataType"}},
							HasSrcOp:  true, SrcPort: "out",
							HasDstOp: true, DstPort: "in1",
						},
						bigMerge,
					}, {
						draw.Arrow{
							DataTypes: []draw.DataType{{Name: "data", Type: "Data"}},
							HasSrcOp:  true, SrcPort: "out",
							HasDstOp: true, DstPort: "in",
						},
						draw.Op{
							Main: draw.DataType{Name: "Mla", Type: "Blue"},
						},
						draw.Arrow{
							DataTypes: []draw.DataType{{Name: "data2", Type: "Data2"}},
							HasSrcOp:  true, SrcPort: "",
							HasDstOp: false, DstPort: "...",
						},
					},
				},
			},
		}, {
			draw.Arrow{
				DataTypes: []draw.DataType{},
				HasSrcOp:  false, SrcPort: "...",
				HasDstOp: true, DstPort: "in",
			},
			draw.Op{
				Main: draw.DataType{Name: "bla2", Type: "Blue"},
			},
			draw.Arrow{
				DataTypes: []draw.DataType{{Name: "data", Type: "Data"}},
				HasSrcOp:  true, SrcPort: "out",
				HasDstOp: true, DstPort: "in2",
			},
			bigMerge,
		}, {
			draw.Arrow{
				DataTypes: []draw.DataType{{Name: "data3", Type: "Data3"}},
				HasSrcOp:  false, SrcPort: "in2",
				HasDstOp: true, DstPort: "in",
			},
			draw.Op{
				Main: draw.DataType{Name: "megaParser", Type: "MegaParser"},
				Plugins: []draw.Plugin{
					{
						Title: "semantics",
						Types: []draw.PluginType{
							{Type: "TextSemantics"},
						},
					},
					{
						Title: "subParser",
						Types: []draw.PluginType{
							{Type: "LitralParser"},
							{Type: "NaturalParser"},
						},
					},
				},
			},
			draw.Arrow{
				DataTypes: []draw.DataType{{Name: "data", Type: "Data"}, {Name: "data2", Type: "data2"}, {Name: "data3", Type: "Data3"}},
				HasSrcOp:  true, SrcPort: "out",
				HasDstOp: true, DstPort: "in3",
			},
			bigMerge,
		}, {
			draw.Op{
				Main: draw.DataType{Type: "bigMerge"},
			},
		}, { // empty to force more space
		}, {
			draw.Arrow{
				DataTypes: []draw.DataType{{Name: "data", Type: "Data"}, {Name: "data2", Type: "data2"}, {Name: "data3", Type: "Data3"}},
				HasSrcOp:  false, SrcPort: "in3",
				HasDstOp: true, DstPort: "",
			},
			draw.Op{
				Main: draw.DataType{Type: "recursive"},
			},
			draw.Arrow{
				DataTypes: []draw.DataType{{Name: "data", Type: "Data"}},
				HasSrcOp:  true, SrcPort: "",
				HasDstOp: true, DstPort: "",
			},
			draw.Op{
				Main: draw.DataType{Type: "secondOp"},
			},
			draw.Arrow{
				DataTypes: []draw.DataType{{Name: "data", Type: "Data"}, {Name: "data2", Type: "data2"}, {Name: "data3", Type: "Data3"}},
				HasSrcOp:  true, SrcPort: "out",
				HasDstOp: true, DstPort: "",
			},
			draw.Text{Text: "recursive"},
		},
	},
}
