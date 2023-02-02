package draw

import (
	"fmt"
)

// Split contains data for multiple paths/arrows originating from a single Comp.
type Split struct {
	Shapes   [][]Shape
	drawData *drawData
}

func (*Split) breakable() bool {
	return true
}

func (*Split) compish() bool {
	return false
}

func (split *Split) intersects(line int) bool {
	return withinShape(line, split.drawData)
}

// --------------------------------------------------------------------------
// Calculate width and height (of other shapes)
// --------------------------------------------------------------------------
func (split *Split) calcDimensions() {
	for _, row := range split.Shapes {
		for _, ishape := range row {
			ishape.calcDimensions()
		}
	}
	split.drawData = &drawData{}
}

// --------------------------------------------------------------------------
// Calculate x0, y0 and minLine
// --------------------------------------------------------------------------
func (split *Split) calcPosition(x0, y0, minLine int, outerComp *drawData,
	lastArr *Arrow, mode FlowMode, merges map[string]*Merge,
) {
	x := x0
	y := y0
	line := minLine
	xmax := x0
	ymax := y0
	maxLine := minLine
	lastComp := (*drawData)(nil)

	for i := 0; i < len(split.Shapes); i++ {
		row := split.Shapes[i]
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
		for j := 0; j < len(row); j++ {
			ishape := row[j]
			switch shape := ishape.(type) {
			case *Arrow:
				shape.calcPosition(x, y, line,
					outerComp, nil, mode, merges,
				)
				lastArr = shape
				x = growX(lastArr.drawData)
				ymax = growY(ymax, lastArr.drawData)
				if lastComp != nil {
					maxLine = growLine(maxLine, lastComp)
				}
				if j == 0 && outerComp != nil {
					growCompToDrawData(outerComp, lastArr.drawData)
				}
				if lastComp != nil {
					growCompToDrawData(lastComp, lastArr.drawData)
				}
			case *Comp:
				shape.calcPosition(x, y, line,
					outerComp, nil, mode, merges,
				)
				lastComp = shape.drawData
				merge := merges[compID(shape)]
				if j == 0 && merge != nil {
					shape.moveTo(merge.drawData)
					growCompToDrawData(lastComp, merge.drawData)

					y = lastComp.y0
					if mode != FlowModeSVGLinks {
						ymax -= LineGap
					}
					line = lastComp.minLine
				}
				if lastArr != nil {
					growCompToDrawData(lastComp, lastArr.drawData)
				}
				x = growX(lastComp)
				ymax = growY(ymax, lastComp)
				maxLine = growLine(maxLine, lastComp)
			case *Split:
				shape.calcPosition(
					x, y, line, lastComp, nil, mode, merges,
				)
				d := shape.drawData
				x = growX(d)
				ymax = growY(ymax, d)
				maxLine = growLine(maxLine, d)
				growCompToDrawData(lastComp, d)
				lastComp = nil
				lastArr = nil
			case *Merge:
				shape.calcPosition(
					x, y, line, nil, lastArr, mode, merges,
				)
				lastComp = nil
				lastArr = nil
			case *Sequel:
				if lastArr != nil {
					lad := lastArr.drawData
					shape.calcPosition(
						x,
						lad.y0+lad.height-LineHeight,
						lad.minLine+lad.lines-1,
						nil, lastArr, mode, merges,
					)
				} else {
					shape.calcPosition(x, y, line,
						nil, lastArr, mode, merges,
					)
				}
				x = growX(shape.drawData)
				lastComp = shape.drawData
				lastArr = nil
			case *Loop:
				lad := lastArr.drawData
				shape.calcPosition(
					x,
					lad.y0+lad.height-LineHeight,
					lad.minLine+lad.lines-1,
					nil, nil, mode, merges,
				)
				x = growX(shape.drawData)
				lastComp = nil
				lastArr = nil
			case *ExtPort:
				if lastArr != nil {
					lad := lastArr.drawData
					shape.calcPosition(
						x, lad.y0+lad.height-LineHeight,
						lad.minLine+lad.lines-1,
						nil, nil, mode, merges,
					)
				} else {
					shape.calcPosition(x, y, line,
						nil, nil, mode, merges,
					)
				}
				x = growX(shape.drawData)
				lastComp = shape.drawData
				lastArr = nil
			default:
				panic(fmt.Sprintf("unsupported type: %T", ishape))
			}
		}
		xmax = max(xmax, x)
	}

	sd := split.drawData
	sd.x0 = x0
	sd.y0 = y0
	sd.height = ymax - y0
	sd.width = xmax - x0
	sd.minLine = minLine
	sd.lines = maxLine - minLine + 1
}

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func (split *Split) enrich(x0, y0, minLine, level int, outerComp *drawData,
	lastArr *Arrow, global *enrichData,
) (newShapeLines [][]Shape) {
	split.calcPosition(x0, y0, minLine, outerComp, lastArr, global.mode, global.merges)
	return nil
}

func (split *Split) breakRows(x0, y0, minLine, level int, outerComp *drawData,
	lastArr *Arrow, global *enrichData,
) (newShapeLines [][]Shape) {
	s := &splitState{
		x:       x0,
		y:       y0,
		line:    minLine,
		xmax:    x0,
		ymax:    y0,
		maxLine: minLine,
	}
	/*
		IDEAS:
		- split enrich into calcDimensions, breakRows, calcPosition
		- check minimum width of the first arrows before looping
		- return early if no space for minimum width
		- merges can be mended if only a minority needs a split
		- Nightmare scenario: merge with last arrow being too long: verylonginput -> merge
	*/

	for s.i = 0; s.i < len(split.Shapes); s.i++ {
		s.row = split.Shapes[s.i]
		s.x = x0
		s.line = s.maxLine
		if s.i > 0 {
			s.line++
			if global.mode != FlowModeSVGLinks {
				s.ymax += LineGap
			}
		}
		s.y = s.ymax
		s.lastComp = nil
		s.lastArr = nil
		for s.j = 0; s.j < len(s.row); s.j++ {
			ishape := s.row[s.j]
			switch shape := ishape.(type) {
			case *Arrow:
				shape.enrich(s.x, s.y, s.line, level,
					outerComp, nil, global,
				)
				s.lastArr = shape
				s.x = growX(s.lastArr.drawData)
				s.ymax = growY(s.ymax, s.lastArr.drawData)
				if s.lastComp != nil {
					s.maxLine = growLine(s.maxLine, s.lastComp)
				}
				if s.j == 0 && outerComp != nil {
					growCompToDrawData(outerComp, s.lastArr.drawData)
				}
				if s.lastComp != nil {
					growCompToDrawData(s.lastComp, s.lastArr.drawData)
				}
			case *Comp:
				shape.enrich(s.x, s.y, s.line, level,
					outerComp, nil, global,
				)
				s.lastComp = shape.drawData
				merge := global.merges[compID(shape)]
				if s.j == 0 && merge != nil {
					shape.moveTo(merge.drawData)
					growCompToDrawData(s.lastComp, merge.drawData)

					s.y = s.lastComp.y0
					if global.mode != FlowModeSVGLinks {
						s.ymax -= LineGap
					}
					s.line = s.lastComp.minLine
				}
				if s.lastArr != nil {
					growCompToDrawData(s.lastComp, s.lastArr.drawData)
				}
				s.x = growX(s.lastComp)
				s.ymax = growY(s.ymax, s.lastComp)
				s.maxLine = growLine(s.maxLine, s.lastComp)
			case *Split:
				nsl := shape.enrich(
					s.x, s.y, s.line, level,
					s.lastComp, nil, global,
				)
				if outerComp != nil {
					newShapeLines = append(newShapeLines, nsl...)
				} else {
					split.Shapes = addRowsAfter(split.Shapes, s.i, nsl)
				}
				d := shape.drawData
				s.x = growX(d)
				s.ymax = growY(s.ymax, d)
				s.maxLine = growLine(s.maxLine, d)
				growCompToDrawData(s.lastComp, d)
				s.lastComp = nil
				s.lastArr = nil
			case *Merge:
				shape.enrich(
					s.x, s.y, s.line, level,
					nil, s.lastArr, global,
				)
				s.lastComp = nil
				s.lastArr = nil
			case *Sequel:
				if s.lastArr != nil {
					lad := s.lastArr.drawData
					shape.enrich(
						s.x,
						lad.y0+lad.height-LineHeight,
						lad.minLine+lad.lines-1, level,
						nil, s.lastArr, global,
					)
				} else {
					shape.enrich(s.x, s.y, s.line, level,
						nil, s.lastArr, global,
					)
				}
				s.x = growX(shape.drawData)
				s.lastComp = shape.drawData
				s.lastArr = nil
			case *Loop:
				lad := s.lastArr.drawData
				shape.enrich(
					s.x,
					lad.y0+lad.height-LineHeight,
					lad.minLine+lad.lines-1, level,
					nil, nil, global,
				)
				s.x = growX(shape.drawData)
				s.lastComp = nil
				s.lastArr = nil
			case *ExtPort:
				if s.lastArr != nil {
					lad := s.lastArr.drawData
					shape.enrich(
						s.x, lad.y0+lad.height-LineHeight,
						lad.minLine+lad.lines-1, level,
						nil, nil, global,
					)
				} else {
					shape.enrich(s.x, s.y, s.line, level,
						nil, nil, global,
					)
				}
				s.x = growX(shape.drawData)
				s.lastComp = shape.drawData
				s.lastArr = nil
			default:
				panic(fmt.Sprintf("unsupported type: %T", ishape))
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

// copyState returns a SHALLOW copy of the given state.
func copyState(stat *splitState) *splitState {
	newStat := *stat
	return &newStat
}

// --------------------------------------------------------------------------
// Convert To SVG and MD
// --------------------------------------------------------------------------
func (split *Split) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	for _, row := range split.Shapes {
		for _, shape := range row {
			if shape.intersects(line) {
				shape.toSVG(smf, line, mode)
			}
		}
	}
}

func addRowsAfter(shapes [][]Shape, i int, newShapes [][]Shape) [][]Shape {
	i++
	shapes = append(shapes, newShapes...)       // grow bigShapes
	copy(shapes[i+len(newShapes):], shapes[i:]) // move everything after i
	copy(shapes[i:], newShapes)                 // add new content
	return shapes
}
