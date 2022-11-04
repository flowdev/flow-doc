package draw

import "fmt"

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func enrichSplit(split *Split, x0, y0, minLine int, outerOp *drawData,
	mode FlowMode, merges map[string]*Merge,
	pluginEnrichArrow func(*Arrow, int, int, int),
	pluginEnrichOp func(*Op, int, int, int),
	pluginEnrichMerge func(*Merge, *drawData, map[string]*Merge),
) {
	var lastOp *drawData
	var lastArr *drawData
	x, y, line, xmax, ymax, maxLine := x0, y0, minLine, x0, y0, minLine

	for i, ss := range split.Shapes {
		x = x0
		line = maxLine
		if i > 0 {
			line++
			if mode != FlowModeSVGLinks {
				ymax += LineGap
			}
		}
		y = ymax
		lastOp = nil
		lastArr = nil
		for j, is := range ss {
			switch s := is.(type) {
			case *Arrow:
				pluginEnrichArrow(s, x, y, line)
				lastArr = s.drawData
				x = growX(lastArr)
				ymax = growY(ymax, lastArr)
				if lastOp != nil {
					maxLine = growLine(maxLine, lastOp)
				}
				if j == 0 && outerOp != nil {
					growOpToDrawData(outerOp, lastArr)
				}
				if lastOp != nil {
					growOpToDrawData(lastOp, lastArr)
				}
			case *Op:
				pluginEnrichOp(s, x, y, line)
				lastOp = s.drawData
				merge := mergeForOp(s, merges)
				if j == 0 && merge != nil {
					moveOp(s, merge.drawData)
					growOpToDrawData(lastOp, merge.drawData)

					y = lastOp.y0
					if mode != FlowModeSVGLinks {
						ymax -= LineGap
					}
					line = lastOp.minLine
				}
				x = growX(lastOp)
				ymax = growY(ymax, lastOp)
				maxLine = growLine(maxLine, lastOp)
				if lastArr != nil {
					growOpToDrawData(lastOp, lastArr)
				}
			case *Split:
				enrichSplit(s, x, y, line, lastOp, mode, merges,
					pluginEnrichArrow, pluginEnrichOp, pluginEnrichMerge)
				d := s.drawData
				x = growX(d)
				ymax = growY(ymax, d)
				maxLine = growLine(maxLine, d)
				growOpToDrawData(lastOp, d)
				lastOp = nil
				lastArr = nil
			case *Merge:
				pluginEnrichMerge(s, lastArr, merges)
				lastOp = nil
				lastArr = nil
			default:
				panic(fmt.Sprintf("unsupported type: %T", is))
			}
		}
		xmax = max(xmax, x)
	}

	split.drawData = &drawData{
		x0:      x0,
		y0:      y0,
		height:  ymax - y0,
		width:   xmax - x0,
		minLine: minLine,
		lines:   maxLine - minLine + 1,
	}
}

func growX(d *drawData) int {
	return d.x0 + d.width
}

func growY(ymax int, d *drawData) int {
	return max(ymax, d.y0+d.height)
}

func growLine(maxLine int, d *drawData) int {
	return max(maxLine, d.minLine+d.lines-1)
}

func growOpToDrawData(op *drawData, d *drawData) {
	op.height = max(op.height, d.y0+d.height-op.y0)
	op.lines = max(op.lines, d.minLine+d.lines-op.minLine)
}

func mergeForOp(op *Op, merges map[string]*Merge) *Merge {
	if op.Main.Name != "" {
		return merges[op.Main.Name]
	} else {
		return merges[op.Main.Type]
	}
}

// --------------------------------------------------------------------------
// Convert To SVG and MD
// --------------------------------------------------------------------------
func splitToSVG(sfs map[string]*svgFlow, mdf *mdFlow, mode FlowMode, split *Split) {
	for _, ss := range split.Shapes {
		for _, is := range ss {
			switch s := is.(type) {
			case *Arrow:
				arrowToSVG(sfs, mdf, mode, s)
			case *Op:
				opToSVG(sfs, mdf, mode, s)
			case *Split:
				splitToSVG(sfs, mdf, mode, s)
			case *Merge:
				// no SVG to create
			default:
				panic(fmt.Sprintf("unsupported type: %T", is))
			}
		}
	}
}

func splitDataToSVG(s Split, sf *svgFlow, lsr *svgRect, x0, y0 int,
) (nsf *svgFlow, xn, yn int) {
	nsf, xn, yn = shapesToSVG(
		s.Shapes,
		sf, x0, y0,
		arrowDataToSVG,
		opDataToSVG,
		textDataToSVG,
		splitDataToSVG,
		mergeDataToSVG,
	)
	adjustLastRect(lsr, yn)
	return
}

func shapesToSVG(
	shapes [][]interface{}, sf *svgFlow, x0 int, y0 int,
	pluginArrowDataToSVG func(Arrow, *svgFlow, *svgRect, int, int) (*svgFlow, int, int, *moveData),
	pluginOpDataToSVG func(Op, *svgFlow, int, int, int) (*svgFlow, *svgRect, int, int, int),
	pluginTextDataToSVG func(Text, *svgFlow, int, int) (*svgFlow, int, int),
	pluginSplitDataToSVG func(Split, *svgFlow, *svgRect, int, int) (*svgFlow, int, int),
	pluginMergeDataToSVG func(Merge, *svgFlow, *moveData, int, int) *myMergeData,
) (nsf *svgFlow, xn, yn int) {
	var xmax, ymax int
	var mod *moveData
	var lsr *svgRect

	for _, ss := range shapes {
		x := x0
		lsr = nil
		if len(ss) < 1 {
			y0 += 48
			continue
		}
		ya := y0
		for _, is := range ss {
			y := y0
			switch s := is.(type) {
			case Arrow:
				sf, x, y, mod = pluginArrowDataToSVG(s, sf, lsr, x, y)
				ya = y - 48 // use the upper arrow Y not the lowest Y
				lsr = nil
			case Op:
				sf, lsr, y0, x, y = pluginOpDataToSVG(s, sf, x, y0, ymax)
				sf.completedMerge = nil
				ya = y0
			case Text:
				sf, x, y = pluginTextDataToSVG(s, sf, x, ya)
			case Split:
				sf, x, y = pluginSplitDataToSVG(s, sf, lsr, x, y)
				lsr = nil
				ya = y0
			case Merge:
				sf.completedMerge = pluginMergeDataToSVG(s, sf, mod, x, y)
				mod = nil
				ya = y0
			default:
				panic(fmt.Sprintf("unsupported type: %T", is))
			}

			ymax = max(ymax, y)
		}
		xmax = max(xmax, x)
		y0 = ymax + 5
	}
	return sf, xmax, ymax
}

func adjustLastRect(lsr *svgRect, yn int) {
	if lsr != nil {
		if lsr.Y+lsr.Height < yn {
			lsr.Height = yn - lsr.Y
		}
	}
}
