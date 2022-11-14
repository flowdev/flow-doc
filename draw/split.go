package draw

import (
	"fmt"
	"strconv"
)

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
type splitState struct {
	lastComp, lastArr               *drawData
	x, y, line, xmax, ymax, maxLine int
	i, j                            int
	ss                              []any
}

func enrichSplit(split *Split, x0, y0, minLine, width int, outerComp *drawData,
	mode FlowMode, merges map[string]*Merge,
) (newShapeLines [][]any) {
	s := &splitState{
		x:       x0,
		y:       y0,
		line:    minLine,
		xmax:    x0,
		ymax:    y0,
		maxLine: minLine,
	}

	for s.i = 0; s.i < len(split.Shapes); s.i++ {
		s.ss = split.Shapes[s.i]
		s.x = x0
		s.line = s.maxLine
		if s.i > 0 {
			s.line++
			if mode != FlowModeSVGLinks {
				s.ymax += LineGap
			}
		}
		s.y = s.ymax
		s.lastComp = nil
		s.lastArr = nil
		for s.j = 0; s.j < len(s.ss); s.j++ {
			is := s.ss[s.j]
			switch sh := is.(type) {
			case *Arrow:
				enrichArrow(sh, s.x, s.y, s.line)
				s.lastArr = sh.drawData
				s.x = growX(s.lastArr)
				s.ymax = growY(s.ymax, s.lastArr)
				if s.lastComp != nil {
					s.maxLine = growLine(s.maxLine, s.lastComp)
				}
				if s.j == 0 && outerComp != nil {
					growCompToDrawData(outerComp, s.lastArr)
				}
				if s.lastComp != nil {
					growCompToDrawData(s.lastComp, s.lastArr)
				}
			case *Comp:
				enrichComp(sh, s.x, s.y, s.line)
				s.lastComp = sh.drawData
				merge := mergeForComp(sh, merges)
				if s.j == 0 && merge != nil {
					moveComp(sh, merge.drawData)
					growCompToDrawData(s.lastComp, merge.drawData)

					s.y = s.lastComp.y0
					if mode != FlowModeSVGLinks {
						s.ymax -= LineGap
					}
					s.line = s.lastComp.minLine
				}
				if s.lastArr != nil {
					growCompToDrawData(s.lastComp, s.lastArr)
				}
				s.x = growX(s.lastComp)
				s.ymax = growY(s.ymax, s.lastComp)
				s.maxLine = growLine(s.maxLine, s.lastComp)
			case *Split:
				nsl := enrichSplit(
					sh, s.x, s.y, s.line, width,
					s.lastComp, mode, merges,
				)
				newShapeLines = append(newShapeLines, nsl...)
				d := sh.drawData
				s.x = growX(d)
				s.ymax = growY(s.ymax, d)
				s.maxLine = growLine(s.maxLine, d)
				growCompToDrawData(s.lastComp, d)
				s.lastComp = nil
				s.lastArr = nil
			case *Merge:
				enrichMerge(sh, s.lastArr, merges)
				s.lastComp = nil
				s.lastArr = nil
			case *Sequel:
				if s.lastArr != nil {
					enrichSequel(
						sh, s.x,
						s.lastArr.y0+s.lastArr.height-LineHeight,
						s.lastArr.minLine+s.lastArr.lines-1,
					)
				} else {
					enrichSequel(sh, s.x, s.y, s.line)
				}
				s.x = growX(sh.drawData)
				s.lastComp = nil
				s.lastArr = nil
			case *Loop:
				enrichLoop(
					sh, s.x,
					s.lastArr.y0+s.lastArr.height-LineHeight,
					s.lastArr.minLine+s.lastArr.lines-1,
				)
				s.x = growX(sh.drawData)
				s.lastComp = nil
				s.lastArr = nil
			default:
				panic(fmt.Sprintf("unsupported type: %T", is))
			}
		}
		s.xmax = max(s.xmax, s.x)
	}

	split.drawData = &drawData{
		x0:      x0,
		y0:      y0,
		height:  s.ymax - y0,
		width:   s.xmax - x0,
		minLine: minLine,
		lines:   s.maxLine - minLine + 1,
	}

	return newShapeLines
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
		svg, _ = addNewSVGFlow(smf,
			sd.x0, sd.y0, sd.height, sd.width,
			"sequel", line,
		)
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
		var svgLink *svgLink
		svg, svgLink = addNewSVGFlow(smf,
			ld.x0, ld.y0, ld.height, ld.width,
			"loop", line,
		)
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
