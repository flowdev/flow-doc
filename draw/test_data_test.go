package draw_test

import "github.com/flowdev/ea-flow-doc/draw"

var BigTestFlowData = draw.Flow{
	Shapes: [][]interface{}{
		{
			&draw.Arrow{
				DataType: []string{"Data"},
				HasSrcOp: false, SrcPort: "in",
				HasDstOp: true, DstPort: "",
			},
			&draw.Op{
				Main: &draw.Rect{
					Text: []string{"ra", "(MiSo)"},
				},
			},
			&draw.Split{
				Shapes: [][]interface{}{
					{
						&draw.Arrow{
							DataType: []string{"Data"},
							HasSrcOp: true, SrcPort: "special",
							HasDstOp: true, DstPort: "in",
						},
						&draw.Op{
							Main: &draw.Rect{
								Text: []string{"do"},
							},
							Plugins: []*draw.Plugin{
								{
									Title: "semantics",
									Rects: []*draw.Rect{
										{Text: []string{"TextSemantics"}},
									},
								},
								{
									Title: "subParser",
									Rects: []*draw.Rect{
										{Text: []string{"LitralParser"}},
										{Text: []string{"NaturalParser"}},
									},
								},
							},
						},
						&draw.Arrow{
							DataType: []string{"BigDataType"},
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: true, DstPort: "in1",
						},
						&draw.Merge{
							ID:   "BigMerge",
							Size: 3,
						},
					}, {
						&draw.Arrow{
							DataType: []string{"Data"},
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: true, DstPort: "in",
						},
						&draw.Op{
							Main: &draw.Rect{
								Text: []string{"bla", "(Blue)"},
							},
						},
						&draw.Arrow{
							DataType: []string{"Data2"},
							HasSrcOp: true, SrcPort: "",
							HasDstOp: false, DstPort: "...",
						},
					},
				},
			},
		}, {
			&draw.Split{
				Shapes: [][]interface{}{
					{
						&draw.Arrow{
							DataType: []string{},
							HasSrcOp: false, SrcPort: "...",
							HasDstOp: true, DstPort: "in",
						},
						&draw.Op{
							Main: &draw.Rect{
								Text: []string{"bla", "(Blue)"},
							},
						},
						&draw.Arrow{
							DataType: []string{"Data"},
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: true, DstPort: "in2",
						},
						&draw.Merge{
							ID:   "BigMerge",
							Size: 3,
						},
					}, {
						&draw.Arrow{
							DataType: []string{"Data3"},
							HasSrcOp: false, SrcPort: "in2",
							HasDstOp: true, DstPort: "in",
						},
						&draw.Op{
							Main: &draw.Rect{
								Text: []string{"megaParser", "(MegaParser)"},
							},
							Plugins: []*draw.Plugin{
								{
									Title: "semantics",
									Rects: []*draw.Rect{
										{Text: []string{"TextSemantics"}},
									},
								},
								{
									Title: "subParser",
									Rects: []*draw.Rect{
										{Text: []string{"LitralParser"}},
										{Text: []string{"NaturalParser"}},
									},
								},
							},
						},
						&draw.Arrow{
							DataType: []string{"(Data,", " data2,", " Data3)"},
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: true, DstPort: "in3",
						},
						&draw.Merge{
							ID:   "BigMerge",
							Size: 3,
						},
					},
				},
			},
		}, {
			&draw.Op{
				Main: &draw.Rect{
					Text: []string{"BigMerge"},
				},
			},
		}, { // empty to force more space
		}, {
			&draw.Arrow{
				DataType: []string{"(Data,", " data2,", " Data3)"},
				HasSrcOp: false, SrcPort: "in3",
				HasDstOp: true, DstPort: "",
			},
			&draw.Op{
				Main: &draw.Rect{
					Text: []string{"recursive"},
				},
			},
			&draw.Arrow{
				DataType: []string{"(Data)"},
				HasSrcOp: true, SrcPort: "",
				HasDstOp: true, DstPort: "",
			},
			&draw.Op{
				Main: &draw.Rect{
					Text: []string{"secondOp"},
				},
			},
			&draw.Arrow{
				DataType: []string{"(Data,", " data2,", " Data3)"},
				HasSrcOp: true, SrcPort: "out",
				HasDstOp: true, DstPort: "",
			},
			&draw.Rect{
				Text: []string{"recursive"},
			},
		},
	},
}
