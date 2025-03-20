package draw2_test

import "github.com/flowdev/ea-flow-doc/draw2"

var BigTestFlowData = buildBigTestFlowData()

func buildBigTestFlowData5() *draw2.Flow {
	cl1 := draw2.NewCluster()
	brk1 := draw2.NewBreakStart(99)
	flow := draw2.NewFlow("bigTestFlow", draw2.FlowModeNoLinks, 1500, false).AddCluster(
		cl1.AddStartComp(
			draw2.NewStartPort("in").AddOutput(
				draw2.NewArrow("", "").AddDataType(
					"data", "Data", "https://google.com?q=Data").AddDestination(
					draw2.NewComp("Xa", "MiSo", "https://google.com?q=Data", cl1).AddOutput( // 1. split
						draw2.NewArrow("special", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw2.NewComp("", "To", "https://google.com?q=To", cl1).AddPluginGroup(
								draw2.NewPluginGroup("semantics").AddPlugin(
									draw2.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
								),
							).AddPluginGroup(
								draw2.NewPluginGroup("subParser").AddPlugin(
									draw2.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
								).AddPlugin(
									draw2.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
								),
							).AddOutput(
								draw2.NewArrow("out", "in1").AddDataType(
									"bigData", "BigDataType", "https://google.com?q=BigDataType").AddDestination(
									draw2.NewComp("", "bigMerge", "https://google.com?q=bigMerge", cl1).AddOutput(
										draw2.NewArrow("", "").AddDataType(
											"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
											draw2.NewComp("postMerge", "PostMerge", "https://google.com?q=PostMerge", cl1).AddOutput( // 2. split
												draw2.NewArrow("", "").AddDataType(
													"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
													draw2.NewComp("", "Split1", "https://google.com?q=Split1", cl1).AddOutput(
														draw2.NewArrow("", "").AddDataType(
															"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
															draw2.NewComp("", "lastMerge", "https://google.com?q=lastMerge", cl1),
														),
													),
												),
											).AddOutput( // 2. split again
												draw2.NewArrow("longNamedOutputPort", "inputPort").AddDataType(
													"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
													draw2.NewComp("", "Split2", "https://google.com?q=Split2", cl1).AddOutput(
														draw2.NewArrow("", "").AddDataType(
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
						draw2.NewArrow("out", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw2.NewComp("Mla", "Blue", "https://google.com?q=Blue", cl1).AddOutput(
								draw2.NewArrow("", "").AddDataType(
									"data2", "Data2", "https://google.com?q=Data2").AddDestination(
									brk1,
								),
							),
						),
					),
				),
			),
		).AddStartComp(
			brk1.End().AddOutput(
				draw2.NewArrow("", "in").AddDestination(
					draw2.NewComp("bla2", "Blue", "https://google.com?q=Blue", cl1).AddOutput(
						draw2.NewArrow("out", "in2").AddDataType(
							"data", "Data", "https://google.com?q=Data").MustLinkComp("bigMerge", cl1),
					),
				),
			),
		).AddStartComp(
			draw2.NewStartPort("in2").AddOutput(
				draw2.NewArrow("", "in").AddDataType(
					"data3", "Data3", "https://google.com?q=Data3").AddDestination(
					draw2.NewComp("megaParser", "MegaParser", "https://google.com?q=MegaParser", cl1).AddPluginGroup(
						draw2.NewPluginGroup("semantics").AddPlugin(
							draw2.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
						),
					).AddPluginGroup(
						draw2.NewPluginGroup("subParser").AddPlugin(
							draw2.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
						).AddPlugin(
							draw2.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
						),
					).AddOutput(
						draw2.NewArrow("out", "in3").AddDataType(
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

func buildBigTestFlowData4() *draw2.Flow {
	cl1 := draw2.NewCluster()
	brk1 := draw2.NewBreakStart(1)
	flow := draw2.NewFlow("bigTestFlow", draw2.FlowModeNoLinks, 1500, false).AddCluster(
		cl1.AddStartComp(
			draw2.NewStartPort("in").AddOutput(
				draw2.NewArrow("", "").AddDataType(
					"data", "Data", "https://google.com?q=Data").AddDestination(
					draw2.NewComp("Xa", "MiSo", "https://google.com?q=Data", cl1).AddOutput( // 1. split
						draw2.NewArrow("special", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw2.NewComp("", "To", "https://google.com?q=To", cl1).AddPluginGroup(
								draw2.NewPluginGroup("semantics").AddPlugin(
									draw2.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
								),
							).AddPluginGroup(
								draw2.NewPluginGroup("subParser").AddPlugin(
									draw2.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
								).AddPlugin(
									draw2.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
								),
							).AddOutput(
								draw2.NewArrow("out", "in1").AddDataType(
									"bigData", "BigDataType", "https://google.com?q=BigDataType").AddDestination(
									draw2.NewComp("", "bigMerge", "https://google.com?q=bigMerge", cl1),
								),
							),
						),
					).AddOutput( // 1. split again
						draw2.NewArrow("out", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw2.NewComp("Mla", "Blue", "https://google.com?q=Blue", cl1).AddOutput(
								draw2.NewArrow("", "").AddDataType(
									"data2", "Data2", "https://google.com?q=Data2").AddDestination(
									brk1,
								),
							),
						),
					),
				),
			),
		).AddStartComp(
			brk1.End().AddOutput(
				draw2.NewArrow("", "in").AddDestination(
					draw2.NewComp("bla2", "Blue", "https://google.com?q=Blue", cl1).AddOutput(
						draw2.NewArrow("out", "in2").AddDataType(
							"data", "Data", "https://google.com?q=Data").MustLinkComp("bigMerge", cl1),
					),
				),
			),
		).AddStartComp(
			draw2.NewStartPort("in2").AddOutput(
				draw2.NewArrow("", "in").AddDataType(
					"data3", "Data3", "https://google.com?q=Data3").AddDestination(
					draw2.NewComp("megaParser", "MegaParser", "https://google.com?q=MegaParser", cl1).AddPluginGroup(
						draw2.NewPluginGroup("semantics").AddPlugin(
							draw2.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
						),
					).AddPluginGroup(
						draw2.NewPluginGroup("subParser").AddPlugin(
							draw2.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
						).AddPlugin(
							draw2.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
						),
					).AddOutput(
						draw2.NewArrow("out", "in3").AddDataType(
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

func buildBigTestFlowData3() *draw2.Flow {
	flow := draw2.NewFlow("bigTestFlow", draw2.FlowModeNoLinks, 1500, false)
	flow.AddCluster(
		draw2.NewCluster().AddStartComp(
			draw2.NewStartPort("in3").AddOutput(
				draw2.NewArrow("", "").AddDataType(
					"data", "Data", "https://google.com?q=Data").AddDataType(
					"data2", "Data2", "https://google.com?q=Data2").AddDataType(
					"data3", "Data3", "https://google.com?q=Data3").AddDestination(
					draw2.NewComp("", "recursive", "https://google.com?q=recursive", nil).AddOutput(
						draw2.NewArrow("", "").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw2.NewComp("", "secondOp", "https://google.com?q=secondOp", nil).AddOutput(
								draw2.NewArrow("out", "").AddDataType(
									"data", "Data", "https://google.com?q=Data").AddDataType(
									"data2", "Data2", "https://google.com?q=Data2").AddDataType(
									"data3", "Data3", "https://google.com?q=Data3").AddDestination(
									draw2.NewLoop("recursive", "in3", "https://google.com?q=recursive"),
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

func buildBigTestFlowData2() *draw2.Flow {
	cl1 := draw2.NewCluster()
	flow := draw2.NewFlow("bigTestFlow", draw2.FlowModeNoLinks, 1500, false).AddCluster(
		cl1.AddStartComp(
			draw2.NewStartPort("in").AddOutput(
				draw2.NewArrow("", "").AddDataType(
					"data", "Data", "https://google.com?q=Data").AddDestination(
					draw2.NewComp("Xa", "MiSo", "https://google.com?q=Data", cl1).AddOutput( // 1. split
						draw2.NewArrow("special", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw2.NewComp("", "To", "https://google.com?q=To", cl1).AddPluginGroup(
								draw2.NewPluginGroup("semantics").AddPlugin(
									draw2.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
								),
							).AddPluginGroup(
								draw2.NewPluginGroup("subParser").AddPlugin(
									draw2.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
								).AddPlugin(
									draw2.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
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

func buildBigTestFlowData() *draw2.Flow {
	cl1 := draw2.NewCluster()
	brk1 := draw2.NewBreakStart(99)
	flow := draw2.NewFlow("bigTestFlow", draw2.FlowModeNoLinks, 1500, false).AddCluster(
		cl1.AddStartComp(
			draw2.NewStartPort("in").AddOutput(
				draw2.NewArrow("", "").AddDataType(
					"data", "Data", "https://google.com?q=Data").AddDestination(
					draw2.NewComp("Xa", "MiSo", "https://google.com?q=Data", cl1).AddOutput( // 1. split
						draw2.NewArrow("special", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw2.NewComp("", "To", "https://google.com?q=To", cl1).AddPluginGroup(
								draw2.NewPluginGroup("semantics").AddPlugin(
									draw2.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
								),
							).AddPluginGroup(
								draw2.NewPluginGroup("subParser").AddPlugin(
									draw2.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
								).AddPlugin(
									draw2.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
								),
							).AddOutput(
								draw2.NewArrow("out", "in1").AddDataType(
									"bigData", "BigDataType", "https://google.com?q=BigDataType").AddDestination(
									draw2.NewComp("", "bigMerge", "https://google.com?q=bigMerge", cl1).AddOutput(
										draw2.NewArrow("", "").AddDataType(
											"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
											draw2.NewComp("postMerge", "PostMerge", "https://google.com?q=PostMerge", cl1).AddOutput( // 2. split
												draw2.NewArrow("", "").AddDataType(
													"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
													draw2.NewComp("", "Split1", "https://google.com?q=Split1", cl1).AddOutput(
														draw2.NewArrow("", "").AddDataType(
															"md1", "MergedData", "https://google.com?q=MergedData").AddDestination(
															draw2.NewComp("", "lastMerge", "https://google.com?q=lastMerge", cl1).AddOutput(
																draw2.NewArrow("", "").AddDestination(draw2.NewEndPort("error")),
															),
														),
													),
												),
											).AddOutput( // 2. split again
												draw2.NewArrow("longNamedOutputPort", "inputPort").AddDataType(
													"data", "MergedData", "https://google.com?q=MergedData").AddDestination(
													draw2.NewComp("", "Split2", "https://google.com?q=Split2", cl1).AddOutput(
														draw2.NewArrow("", "").AddDataType(
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
						draw2.NewArrow("out", "in").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw2.NewComp("Mla", "Blue", "https://google.com?q=Blue", cl1).AddOutput(
								draw2.NewArrow("", "").AddDataType(
									"data2", "Data2", "https://google.com?q=Data2").AddDestination(
									brk1,
								),
							),
						),
					),
				),
			),
		).AddStartComp(
			brk1.End().AddOutput(
				draw2.NewArrow("", "in").AddDestination(
					draw2.NewComp("bla2", "Blue", "https://google.com?q=Blue", cl1).AddOutput(
						draw2.NewArrow("out", "in2").AddDataType(
							"data", "Data", "https://google.com?q=Data").MustLinkComp("bigMerge", cl1),
					),
				),
			),
		).AddStartComp(
			draw2.NewStartPort("in2").AddOutput(
				draw2.NewArrow("", "in").AddDataType(
					"data3", "Data3", "https://google.com?q=Data3").AddDestination(
					draw2.NewComp("megaParser", "MegaParser", "https://google.com?q=MegaParser", cl1).AddPluginGroup(
						draw2.NewPluginGroup("semantics").AddPlugin(
							draw2.NewPlugin("TextSemantics", "https://google.com?q=TextSemantics"),
						),
					).AddPluginGroup(
						draw2.NewPluginGroup("subParser").AddPlugin(
							draw2.NewPlugin("LiteralParser", "https://google.com?q=LiteralParser").GoLink(),
						).AddPlugin(
							draw2.NewPlugin("NaturalParser", "https://google.com?q=NaturalParser"),
						),
					).AddOutput(
						draw2.NewArrow("out", "in3").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDataType(
							"data2", "Data2", "https://google.com?q=Data2").AddDataType(
							"data3", "Data3", "https://google.com?q=Data3").MustLinkComp("bigMerge", cl1),
					),
				),
			),
		),
	)
	flow.AddCluster(
		draw2.NewCluster().AddStartComp(
			draw2.NewStartPort("in3").AddOutput(
				draw2.NewArrow("", "").AddDataType(
					"data", "Data", "https://google.com?q=Data").AddDataType(
					"data2", "data2", "https://google.com?q=data2").AddDataType(
					"data3", "Data3", "https://google.com?q=Data3").AddDestination(
					draw2.NewComp("", "recursive", "https://google.com?q=recursive", nil).AddOutput(
						draw2.NewArrow("", "").AddDataType(
							"data", "Data", "https://google.com?q=Data").AddDestination(
							draw2.NewComp("", "secondOp", "https://google.com?q=secondOp", nil).AddOutput(
								draw2.NewArrow("out", "").AddDataType(
									"data", "Data", "https://google.com?q=Data").AddDataType(
									"data2", "data2", "https://google.com?q=data2").AddDataType(
									"data3", "Data3", "https://google.com?q=Data3").AddDestination(
									draw2.NewLoop("recursive", "in3", "https://google.com?q=recursive"),
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
