package draw

import (
	"fmt"
	"strconv"
)

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func enrichSplit(split *Split, x0, y0, minLine int, outerComp *drawData,
	mode FlowMode, merges map[string]*Merge,
) {
	var lastComp *drawData
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
		lastComp = nil
		lastArr = nil
		for j, is := range ss {
			switch s := is.(type) {
			case *Arrow:
				enrichArrow(s, x, y, line)
				lastArr = s.drawData
				x = growX(lastArr)
				ymax = growY(ymax, lastArr)
				if lastComp != nil {
					maxLine = growLine(maxLine, lastComp)
				}
				if j == 0 && outerComp != nil {
					growCompToDrawData(outerComp, lastArr)
				}
				if lastComp != nil {
					growCompToDrawData(lastComp, lastArr)
				}
			case *Comp:
				enrichComp(s, x, y, line)
				lastComp = s.drawData
				merge := mergeForComp(s, merges)
				if j == 0 && merge != nil {
					moveComp(s, merge.drawData)
					growCompToDrawData(lastComp, merge.drawData)

					y = lastComp.y0
					if mode != FlowModeSVGLinks {
						ymax -= LineGap
					}
					line = lastComp.minLine
				}
				if lastArr != nil {
					growCompToDrawData(lastComp, lastArr)
				}
				x = growX(lastComp)
				ymax = growY(ymax, lastComp)
				maxLine = growLine(maxLine, lastComp)
			case *Split:
				enrichSplit(s, x, y, line, lastComp, mode, merges)
				d := s.drawData
				x = growX(d)
				ymax = growY(ymax, d)
				maxLine = growLine(maxLine, d)
				growCompToDrawData(lastComp, d)
				lastComp = nil
				lastArr = nil
			case *Merge:
				enrichMerge(s, lastArr, merges)
				lastComp = nil
				lastArr = nil
			case *Sequel:
				if lastArr != nil {
					enrichSequel(
						s, x,
						lastArr.y0+lastArr.height-LineHeight,
						lastArr.minLine+lastArr.lines-1,
					)
				} else {
					enrichSequel(s, x, y, line)
				}
				x = growX(s.drawData)
				lastComp = nil
				lastArr = nil
			case *Loop:
				enrichLoop(
					s, x,
					lastArr.y0+lastArr.height-LineHeight,
					lastArr.minLine+lastArr.lines-1,
				)
				x = growX(s.drawData)
				lastComp = nil
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

func growCompToDrawData(comp *drawData, d *drawData) {
	comp.height = max(comp.height, d.y0+d.height-comp.y0)
	comp.lines = max(comp.lines, d.minLine+d.lines-comp.minLine)
}

func mergeForComp(comp *Comp, merges map[string]*Merge) *Merge {
	if comp.Main.Name != "" {
		return merges[comp.Main.Name]
	} else {
		return merges[comp.Main.Type]
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
			case *Comp:
				if withinShape(line, s.drawData) {
					xDiff := s.drawData.x0 - smf.lastX
					if mode == FlowModeSVGLinks && xDiff > 0 {
						addFillerSVG(smf, line, smf.lastX, LineHeight, xDiff)
						smf.lastX += xDiff
					}
					compToSVG(smf, line, mode, s)
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
		Y:     sd.y0 + sd.height - arrTextOffset,
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
		svgLink := addSVGLinkToMDFlowLines(smf, line, name, "loop")
		svgLink.Link = loop.Link
	} else {
		svg = smf.svgs[""]
	}

	txt := "...back to: " + loop.Name
	if loop.Port != "" {
		txt += ":" + loop.Port
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:      ld.x0,
		Y:      ld.y0 + ld.height - arrTextOffset,
		Width:  ld.width,
		Text:   txt,
		Link:   !loop.GoLink && loop.Link != "",
		GoLink: loop.GoLink,
	})
}
