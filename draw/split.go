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
// Add drawData
// --------------------------------------------------------------------------
func enrichSplit(split *Split, x0, y0, minLine, level int, outerComp *drawData,
	global *enrichData,
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
				//splitFit := enrichArrow(shape, s.x, s.y, s.line)
				enrichArrow(shape, s.x, s.y, s.line)
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
				enrichComp(shape, s.x, s.y, s.line)
				s.lastComp = shape.drawData
				merge := global.merges[compID(shape)]
				if s.j == 0 && merge != nil {
					moveComp(shape, merge.drawData)
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
				nsl := enrichSplit(
					shape, s.x, s.y, s.line, level,
					s.lastComp, global,
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
				enrichMerge(shape, s.lastArr, global.merges)
				s.lastComp = nil
				s.lastArr = nil
			case *Sequel:
				if s.lastArr != nil {
					lad := s.lastArr.drawData
					enrichSequel(
						shape, s.x,
						lad.y0+lad.height-LineHeight,
						lad.minLine+lad.lines-1,
					)
				} else {
					enrichSequel(shape, s.x, s.y, s.line)
				}
				s.x = growX(shape.drawData)
				s.lastComp = shape.drawData
				s.lastArr = nil
			case *Loop:
				lad := s.lastArr.drawData
				enrichLoop(
					shape, s.x,
					lad.y0+lad.height-LineHeight,
					lad.minLine+lad.lines-1,
				)
				s.x = growX(shape.drawData)
				s.lastComp = nil
				s.lastArr = nil
			case *ExtPort:
				if s.lastArr != nil {
					lad := s.lastArr.drawData
					enrichExtPort(
						shape, s.x,
						lad.y0+lad.height-LineHeight,
						lad.minLine+lad.lines-1,
					)
				} else {
					enrichExtPort(shape, s.x, s.y, s.line)
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
