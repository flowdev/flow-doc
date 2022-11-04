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
