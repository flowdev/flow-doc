package draw

import (
	"fmt"
	"strconv"
)

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
			case *Sequel:
				enrichSequel(s, x, y, line)
				x = growX(s.drawData)
				lastOp = nil
				lastArr = nil
			case *Loop:
				enrichLoop(s, x, y, line)
				x = growX(s.drawData)
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

func enrichSequel(seq *Sequel, x0, y0, minLine int) {
	width := ParenWidth*3 + len(strconv.Itoa(seq.Number))*CharWidth

	seq.drawData = &drawData{
		x0:      x0,
		y0:      y0,
		width:   width,
		height:  LineHeight,
		minLine: minLine,
		lines:   1,
	}
}

func enrichLoop(loop *Loop, x0, y0, minLine int) {
	txt := "back to: " + loop.Name + loop.Port
	width := ParenWidth*3 + len(txt)*CharWidth
	if loop.Port != "" {
		width += ParenWidth
	}

	loop.drawData = &drawData{
		x0:      x0,
		y0:      y0,
		width:   width,
		height:  LineHeight,
		minLine: minLine,
		lines:   1,
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
			case *Sequel:
				if withinShape(line, s.drawData) {
					sequelToSVG(smf, line, mode, s)
					smf.lastX += s.drawData.width
				}
			case *Loop:
				if withinShape(line, s.drawData) {
					loopToSVG(smf, line, mode, s)
					smf.lastX += s.drawData.width
				}
			default:
				panic(fmt.Sprintf("unsupported type: %T", is))
			}
		}
	}
}

func sequelToSVG(smf *svgMDFlow, line int, mode FlowMode, seq *Sequel) {
	var svg *svgFlow
	sd := seq.drawData

	// get or create correct SVG flow:
	if mode == FlowModeSVGLinks {
		svg = newSVGFlow(sd.x0, sd.y0, sd.height, sd.width, tinyDiagramSize)
		name := svgFileName(smf, "sequel", line)
		smf.svgs[name] = svg
		addSVGLinkToMDFlowLines(smf, line, name, "sequel")
	} else {
		svg = smf.svgs[""]
	}

	svg.Texts = append(svg.Texts, &svgText{
		X:     sd.x0,
		Y:     sd.y0 + sd.height - TextOffset,
		Width: sd.width,
		Text:  "..." + strconv.Itoa(seq.Number),
	})
}

func loopToSVG(smf *svgMDFlow, line int, mode FlowMode, loop *Loop) {
	var svg *svgFlow
	ld := loop.drawData

	// get or create correct SVG flow:
	if mode == FlowModeSVGLinks {
		svg = newSVGFlow(ld.x0, ld.y0, ld.height, ld.width, tinyDiagramSize)
		name := svgFileName(smf, "loop", line)
		smf.svgs[name] = svg
		addSVGLinkToMDFlowLines(smf, line, name, "loop")
	} else {
		svg = smf.svgs[""]
	}

	txt := "...back to: " + loop.Name
	if loop.Port != "" {
		txt += ":" + loop.Port
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:     ld.x0,
		Y:     ld.y0 + ld.height - TextOffset,
		Width: ld.width,
		Text:  txt,
	})
}
