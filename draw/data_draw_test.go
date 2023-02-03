package draw_test

import "github.com/flowdev/ea-flow-doc/draw"

var bigMerge = &draw.Merge{
	ID:   "bigMerge",
	Size: 3,
}

var lastMerge = &draw.Merge{
	ID:   "lastMerge",
	Size: 2,
}

var BigTestFlowData = &draw.Flow{
	Name: "bigTestFlow",
	AllShapes: &draw.Split{
		Shapes: [][]draw.Shape{
			{
				&draw.ExtPort{
					Name: "in",
				},
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
					},
				},
				&draw.Comp{
					Main: &draw.DataType{
						Name: "Xa", Type: "MiSo", Link: "https://google.com?q=MiSo",
					},
				},
				&draw.Split{
					Shapes: [][]draw.Shape{
						{
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
								},
								SrcPort: "special",
								DstPort: "in",
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
								SrcPort: "out",
								DstPort: "in1",
							},
							bigMerge,
						}, {
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
								},
								SrcPort: "out",
								DstPort: "in",
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
					DstPort:   "in",
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
					SrcPort: "out",
					DstPort: "in2",
				},
				bigMerge,
			}, {
				&draw.ExtPort{
					Name: "in2",
				},
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data3", Type: "Data3", Link: "https://google.com?q=Data3"},
					},
					DstPort: "in",
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
					SrcPort: "out",
					DstPort: "in3",
				},
				bigMerge,
			}, {
				&draw.Comp{
					Main: &draw.DataType{
						Type: "bigMerge", Link: "https://google.com?q=bigMerge",
					},
				},
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data", Type: "MergedData", Link: "https://google.com?q=MergedData"},
					},
				},
				&draw.Comp{
					Main: &draw.DataType{
						Name: "postMerge", Type: "PostMerge", Link: "https://google.com?q=PostMerge",
					},
				},
				&draw.Split{
					Shapes: [][]draw.Shape{
						{
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "data", Type: "MergedData", Link: "https://google.com?q=MergedData"},
								},
							},
							&draw.Comp{
								Main: &draw.DataType{
									Type: "Split1", Link: "https://google.com?q=Split1",
								},
							},
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "data", Type: "MergedData", Link: "https://google.com?q=MergedData"},
								},
							},
							lastMerge,
						}, {
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "data", Type: "MergedData", Link: "https://google.com?q=MergedData"},
								},
								SrcPort: "longNamedOutputPort",
								DstPort: "inputPort",
							},
							&draw.Comp{
								Main: &draw.DataType{
									Type: "Split2", Link: "https://google.com?q=Split2",
								},
							},
							&draw.Arrow{
								DataTypes: []*draw.DataType{
									{Name: "data", Type: "MergedData", Link: "https://google.com?q=MergedData"},
								},
							},
							lastMerge,
						},
					},
				},
			}, {
				&draw.Comp{
					Main: &draw.DataType{
						Type: "lastMerge", Link: "https://google.com?q=lastMerge",
					},
				},
			}, { // empty to force more space
			}, {
				&draw.ExtPort{
					Name: "in3",
				},
				&draw.Arrow{
					DataTypes: []*draw.DataType{
						{Name: "data", Type: "Data", Link: "https://google.com?q=Data"},
						{Name: "data2", Type: "data2", Link: "https://google.com?q=data2"},
						{Name: "data3", Type: "Data3", Link: "https://google.com?q=Data3"},
					},
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
					SrcPort: "out",
				},
				&draw.Loop{
					Name: "recursive", Port: "in3", Link: "https://google.com?q=recursive",
				},
			},
		},
	},
}
