package draw

import "fmt"

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func enrichSplit(split *Split, x0, y0, minLine int, outerOp *drawData,
	mode FlowMode, merges map[string]*Merge,
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
		if y != line*LineHeight {
			fmt.Printf("y should be %d lines * 24 = %d, got: %d\n", line, line*LineHeight, y)
		}
		for j, is := range ss {
			switch s := is.(type) {
			case *Arrow:
				enrichArrow(s, x, y, line)
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
				enrichOp(s, x, y, line)
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
				enrichSplit(s, x, y, line, lastOp, mode, merges)
				d := s.drawData
				x = growX(d)
				ymax = growY(ymax, d)
				maxLine = growLine(maxLine, d)
				growOpToDrawData(lastOp, d)
				lastOp = nil
				lastArr = nil
			case *Merge:
				enrichMerge(s, lastArr, merges)
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
func splitToSVG(smf *svgMDFlow, line int, mode FlowMode, split *Split) {
	for _, ss := range split.Shapes {
		for _, is := range ss {
			switch s := is.(type) {
			case *Arrow:
				if withinShape(line, s.drawData) {
					arrowToSVG(smf, line, mode, s)
					smf.lastX += s.drawData.width
				}
			case *Op:
				if withinShape(line, s.drawData) {
					xDiff := s.drawData.x0 - smf.lastX
					if mode == FlowModeSVGLinks && xDiff > 0 {
						addFillerSVG(smf, line, smf.lastX, LineHeight, xDiff)
						smf.lastX += xDiff
					}
					opToSVG(smf, line, mode, s)
					smf.lastX += s.drawData.width
				}
			case *Split:
				if withinShape(line, s.drawData) {
					splitToSVG(smf, line, mode, s)
				}
			case *Merge:
				// no SVG to create
			default:
				panic(fmt.Sprintf("unsupported type: %T", is))
			}
		}
	}
}
