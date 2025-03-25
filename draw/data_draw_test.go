package draw_test

import "github.com/flowdev/ea-flow-doc/draw"

func buildBigTestFlowData() *draw.Flow {
	return buildBigTestFlowData1()
}

func buildBigTestFlowData5() *draw.Flow {
	cl1 := draw.NewCluster()
	flow := draw.NewFlow("bigTestFlow", draw.FlowModeNoLinks, 1500, false).AddCluster(
		cl1.AddStartComp(
			draw.NewStartPort("in").AddOutput(
				draw.NewArrow("", "").AddDataType(
					"data", "Data", "https://google.com?q=Data").AddDestination(
					draw.NewComp("Xa", "MiSo", "https://google.com?q=Data", cl1).AddOutput( // 1. split
						draw.NewArrow("special", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw.NewComp("", "To", "https://google.com?q=To", cl1).AddPluginGroup(
								draw.NewPluginGroup("semantics").AddPlugin(
									draw.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
								),
							).AddPluginGroup(
								draw.NewPluginGroup("subParser").AddPlugin(
									draw.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
								).AddPlugin(
									draw.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
								),
							).AddOutput(
								draw.NewArrow("out", "in1").AddDataType(
									"bigData", "BigDataType", "https://google.com?q=BigDataType").AddDestination(
									draw.NewComp("", "bigMerge", "https://google.com?q=bigMerge", cl1).AddOutput(
										draw.NewArrow("", "").AddDataType(
											"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
											draw.NewComp("postMerge", "PostMerge", "https://google.com?q=PostMerge", cl1).AddOutput( // 2. split
												draw.NewArrow("", "").AddDataType(
													"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
													draw.NewComp("", "Split1", "https://google.com?q=Split1", cl1).AddOutput(
														draw.NewArrow("", "").AddDataType(
															"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
															draw.NewComp("", "lastMerge", "https://google.com?q=lastMerge", cl1),
														),
													),
												),
											).AddOutput( // 2. split again
												draw.NewArrow("longNamedOutputPort", "inputPort").AddDataType(
													"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
													draw.NewComp("", "Split2", "https://google.com?q=Split2", cl1).AddOutput(
														draw.NewArrow("", "").AddDataType(
															"data", "MergedData", "https://google.com?q=MergedData").MustLinkComp("lastMerge", cl1),
													),
												),
											),
										),
									),
								),
							),
						),
					).AddOutput( // 1. split again
						draw.NewArrow("out", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw.NewComp("Mla", "Blue", "https://google.com?q=Blue", cl1).AddOutput(
								draw.NewArrow("", "in").AddDataType(
									"data2", "Data2", "https://google.com?q=Data2").AddDestination(
									draw.NewComp("bla2", "Blue", "https://google.com?q=Blue", cl1).AddOutput(
										draw.NewArrow("out", "in2").AddDataType(
											"data", "Data", "https://google.com?q=Data").MustLinkComp("bigMerge", cl1),
									),
								),
							),
						),
					),
				),
			),
		).AddStartComp(
			draw.NewStartPort("in2").AddOutput(
				draw.NewArrow("", "in").AddDataType(
					"data3", "Data3", "https://google.com?q=Data3").AddDestination(
					draw.NewComp("megaParser", "MegaParser", "https://google.com?q=MegaParser", cl1).AddPluginGroup(
						draw.NewPluginGroup("semantics").AddPlugin(
							draw.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
						),
					).AddPluginGroup(
						draw.NewPluginGroup("subParser").AddPlugin(
							draw.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
						).AddPlugin(
							draw.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
						),
					).AddOutput(
						draw.NewArrow("out", "in3").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDataType(
							"data2", "Data2", "https://google.com?q=Data2").AddDataType(
							"data3", "Data3", "https://google.com?q=Data3").MustLinkComp("bigMerge", cl1),
					),
				),
			),
		),
	)

	return flow
}

func buildBigTestFlowData4() *draw.Flow {
	cl1 := draw.NewCluster()
	flow := draw.NewFlow("bigTestFlow", draw.FlowModeNoLinks, 1500, false).AddCluster(
		cl1.AddStartComp(
			draw.NewStartPort("in").AddOutput(
				draw.NewArrow("", "").AddDataType(
					"data", "Data", "https://google.com?q=Data").AddDestination(
					draw.NewComp("Xa", "MiSo", "https://google.com?q=Data", cl1).AddOutput( // 1. split
						draw.NewArrow("special", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw.NewComp("", "To", "https://google.com?q=To", cl1).AddPluginGroup(
								draw.NewPluginGroup("semantics").AddPlugin(
									draw.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
								),
							).AddPluginGroup(
								draw.NewPluginGroup("subParser").AddPlugin(
									draw.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
								).AddPlugin(
									draw.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
								),
							).AddOutput(
								draw.NewArrow("out", "in1").AddDataType(
									"bigData", "BigDataType", "https://google.com?q=BigDataType").AddDestination(
									draw.NewComp("", "bigMerge", "https://google.com?q=bigMerge", cl1),
								),
							),
						),
					).AddOutput( // 1. split again
						draw.NewArrow("out", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw.NewComp("Mla", "Blue", "https://google.com?q=Blue", cl1).AddOutput(
								draw.NewArrow("", "in").AddDataType(
									"data2", "Data2", "https://google.com?q=Data2").AddDestination(
									draw.NewComp("bla2", "Blue", "https://google.com?q=Blue", cl1).AddOutput(
										draw.NewArrow("out", "in2").AddDataType(
											"data", "Data", "https://google.com?q=Data").MustLinkComp("bigMerge", cl1),
									),
								),
							),
						),
					),
				),
			),
		).AddStartComp(
			draw.NewStartPort("in2").AddOutput(
				draw.NewArrow("", "in").AddDataType(
					"data3", "Data3", "https://google.com?q=Data3").AddDestination(
					draw.NewComp("megaParser", "MegaParser", "https://google.com?q=MegaParser", cl1).AddPluginGroup(
						draw.NewPluginGroup("semantics").AddPlugin(
							draw.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
						),
					).AddPluginGroup(
						draw.NewPluginGroup("subParser").AddPlugin(
							draw.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
						).AddPlugin(
							draw.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
						),
					).AddOutput(
						draw.NewArrow("out", "in3").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDataType(
							"data2", "Data2", "https://google.com?q=Data2").AddDataType(
							"data3", "Data3", "https://google.com?q=Data3").MustLinkComp("bigMerge", cl1),
					),
				),
			),
		),
	)

	return flow
}

func buildBigTestFlowData3() *draw.Flow {
	flow := draw.NewFlow("bigTestFlow", draw.FlowModeNoLinks, 1500, false)
	flow.AddCluster(
		draw.NewCluster().AddStartComp(
			draw.NewStartPort("in3").AddOutput(
				draw.NewArrow("", "").AddDataType(
					"data", "Data", "https://google.com?q=Data").AddDataType(
					"data2", "Data2", "https://google.com?q=Data2").AddDataType(
					"data3", "Data3", "https://google.com?q=Data3").AddDestination(
					draw.NewComp("", "recursive", "https://google.com?q=recursive", nil).AddOutput(
						draw.NewArrow("", "").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw.NewComp("", "secondOp", "https://google.com?q=secondOp", nil).AddOutput(
								draw.NewArrow("out", "").AddDataType(
									"data", "Data", "https://google.com?q=Data").AddDataType(
									"data2", "Data2", "https://google.com?q=Data2").AddDataType(
									"data3", "Data3", "https://google.com?q=Data3").AddDestination(
									draw.NewLoop("recursive", "in3", "https://google.com?q=recursive"),
								),
							),
						),
					),
				),
			),
		),
	)

	return flow
}

func buildBigTestFlowData2() *draw.Flow {
	cl1 := draw.NewCluster()
	flow := draw.NewFlow("bigTestFlow", draw.FlowModeNoLinks, 1500, false).AddCluster(
		cl1.AddStartComp(
			draw.NewStartPort("in").AddOutput(
				draw.NewArrow("", "").AddDataType(
					"data", "Data", "https://google.com?q=Data").AddDestination(
					draw.NewComp("Xa", "MiSo", "https://google.com?q=Data", cl1).AddOutput( // 1. split
						draw.NewArrow("special", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw.NewComp("", "To", "https://google.com?q=To", cl1).AddPluginGroup(
								draw.NewPluginGroup("semantics").AddPlugin(
									draw.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
								),
							).AddPluginGroup(
								draw.NewPluginGroup("subParser").AddPlugin(
									draw.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
								).AddPlugin(
									draw.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
								),
							),
						),
					),
				),
			),
		),
	)

	return flow
}

func buildBigTestFlowData1() *draw.Flow {
	cl1 := draw.NewCluster()
	flow := draw.NewFlow("bigTestFlow", draw.FlowModeNoLinks, 1500, false).AddCluster(
		cl1.AddStartComp(
			draw.NewStartPort("in").AddOutput(
				draw.NewArrow("", "").AddDataType(
					"data", "Data", "https://google.com?q=Data").AddDestination(
					draw.NewComp("Xa", "MiSo", "https://google.com?q=Data", cl1).AddOutput( // 1. split
						draw.NewArrow("special", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw.NewComp("", "To", "https://google.com?q=To", cl1).AddPluginGroup(
								draw.NewPluginGroup("semantics").AddPlugin(
									draw.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
								),
							).AddPluginGroup(
								draw.NewPluginGroup("subParser").AddPlugin(
									draw.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
								).AddPlugin(
									draw.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
								),
							).AddOutput(
								draw.NewArrow("out", "in1").AddDataType(
									"bigData", "BigDataType", "https://google.com?q=BigDataType").AddDestination(
									draw.NewComp("", "bigMerge", "https://google.com?q=bigMerge", cl1).AddOutput(
										draw.NewArrow("", "").AddDataType(
											"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
											draw.NewComp("postMerge", "PostMerge", "https://google.com?q=PostMerge", cl1).AddOutput( // 2. split
												draw.NewArrow("", "").AddDataType(
													"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
													draw.NewComp("", "Split1", "https://google.com?q=Split1", cl1).AddOutput(
														draw.NewArrow("", "").AddDataType(
															"md1", "MergedData", "https://google.com?q=MergedData").AddDestination(
															draw.NewComp("", "lastMerge", "https://google.com?q=lastMerge", cl1).AddOutput(
																draw.NewArrow("", "").AddDestination(draw.NewEndPort("error")),
															),
														),
													),
												),
											).AddOutput( // 2. split again
												draw.NewArrow("longNamedOutputPort", "inputPort").AddDataType(
													"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
													draw.NewComp("", "Split2", "https://google.com?q=Split2", cl1).AddOutput(
														draw.NewArrow("", "").AddDataType(
															"md2", "MergedData", "https://google.com?q=MergedData").MustLinkComp("lastMerge", cl1),
													),
												),
											),
										),
									),
								),
							),
						),
					).AddOutput( // 1. split again
						draw.NewArrow("out", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw.NewComp("Mla", "Blue", "https://google.com?q=Blue", cl1).AddOutput(
								draw.NewArrow("", "in").AddDataType(
									"data2", "Data2", "https://google.com?q=Data2").AddDestination(
									draw.NewComp("bla2", "Blue", "https://google.com?q=Blue", cl1).AddOutput(
										draw.NewArrow("out", "in2").AddDataType(
											"data", "Data", "https://google.com?q=Data").MustLinkComp("bigMerge", cl1),
									),
								),
							),
						),
					),
				),
			),
		).AddStartComp(
			draw.NewStartPort("in2").AddOutput(
				draw.NewArrow("", "in").AddDataType(
					"data3", "Data3", "https://google.com?q=Data3").AddDestination(
					draw.NewComp("megaParser", "MegaParser", "https://google.com?q=MegaParser", cl1).AddPluginGroup(
						draw.NewPluginGroup("semantics").AddPlugin(
							draw.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
						),
					).AddPluginGroup(
						draw.NewPluginGroup("subParser").AddPlugin(
							draw.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
						).AddPlugin(
							draw.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
						),
					).AddOutput(
						draw.NewArrow("out", "in3").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDataType(
							"data2", "Data2", "https://google.com?q=Data2").AddDataType(
							"data3", "Data3", "https://google.com?q=Data3").MustLinkComp("bigMerge", cl1),
					),
				),
			),
		),
	)
	flow.AddCluster(
		draw.NewCluster().AddStartComp(
			draw.NewStartPort("in3").AddOutput(
				draw.NewArrow("", "").AddDataType(
					"data", "Data", "https://google.com?q=Data").AddDataType(
					"data2", "data2", "https://google.com?q=data2").AddDataType(
					"data3", "Data3", "https://google.com?q=Data3").AddDestination(
					draw.NewComp("", "recursive", "https://google.com?q=recursive", nil).AddOutput(
						draw.NewArrow("", "").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw.NewComp("", "secondOp", "https://google.com?q=secondOp", nil).AddOutput(
								draw.NewArrow("out", "").AddDataType(
									"data", "Data", "https://google.com?q=Data").AddDataType(
									"data2", "data2", "https://google.com?q=data2").AddDataType(
									"data3", "Data3", "https://google.com?q=Data3").AddDestination(
									draw.NewLoop("recursive", "in3", "https://google.com?q=recursive"),
								),
							),
						),
					),
				),
			),
		),
	)

	return flow
}
