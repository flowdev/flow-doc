package draw

import (
	"strconv"
)

// --------------------------------------------------------------------------
// BreakStart
// --------------------------------------------------------------------------

type BreakStart struct {
	withDrawData
	number     int
	input      *Arrow
	inReturned bool
	end        *BreakEnd
}

func NewBreakStart(num int) *BreakStart {
	return &BreakStart{
		number: num,
	}
}

func (brk *BreakStart) addInput(arr *Arrow) {
	brk.input = arr
}

func (brk *BreakStart) switchInput(oldArr, newArr *Arrow) {
	newArr.dstComp = brk
	brk.input = newArr
}

func (brk *BreakStart) minRestOfRowWidth(num int) int {
	return brk.drawData.width
}

func (brk *BreakStart) calcHorizontalValues(x0 int) {
	brk.drawData = breakHorizontalValues(x0, brk.number)
}

func (brk *BreakStart) respectMaxWidth(maxWidth, num int) (newStartComps []StartComp, newNum, newWidth int) {
	return nil, num, brk.drawData.xmax()
}

func (brk *BreakStart) calcVerticalValues(y0, minLine int, _ FlowMode) (maxLines, newHeight int) {
	bd := brk.drawData
	breakVerticalValues(bd, brk.input.drawData.ymax(), brk.input.drawData.maxLines())
	return bd.maxLines(), bd.ymax()
}

func (brk *BreakStart) End() *BreakEnd {
	if brk.end == nil {
		brk.end = &BreakEnd{
			number: brk.number,
		}
	}
	return brk.end
}

func (brk *BreakStart) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	if brk.drawData.drawLine(line) {
		breakToSVG(smf, line, mode, brk.drawData, brk.number)
		brk.drawData.drawnLines[line] = true
	}
}

// --------------------------------------------------------------------------
// BreakEnd
// --------------------------------------------------------------------------

type BreakEnd struct {
	withDrawData
	number      int
	output      *Arrow
	outReturned bool
}

func (brk *BreakEnd) AddOutput(arr *Arrow) *BreakEnd {
	arr.srcComp = brk
	brk.output = arr
	return brk
}

func (brk *BreakEnd) minRestOfRowWidth(num int) int {
	return brk.drawData.width + brk.output.minRestOfRowWidth(num)
}

func (brk *BreakEnd) calcHorizontalValues(x0 int) {
	brk.withDrawData.drawData = breakHorizontalValues(x0, brk.number)
	brk.output.calcHorizontalValues(brk.drawData.x0 + brk.drawData.width)
}

func (brk *BreakEnd) respectMaxWidth(maxWidth, num int) (newStartComps []StartComp, newNum, newWidth int) {
	return brk.output.respectMaxWidth(maxWidth, num)
}

func (brk *BreakEnd) resetDrawData() {
	brk.withDrawData.resetDrawData()
	brk.output.resetDrawData()
}

func (brk *BreakEnd) calcVerticalValues(y0, minLine int, mode FlowMode) (maxLines, newHeight int) {
	bd := brk.drawData
	maxLines, newHeight = brk.output.calcVerticalValues(y0, minLine, mode)
	breakVerticalValues(bd, brk.output.drawData.ymax(), brk.output.drawData.maxLines())
	return maxLines, newHeight
}

func (brk *BreakEnd) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	if brk.drawData.drawLine(line) {
		breakToSVG(smf, line, mode, brk.drawData, brk.number)
		brk.drawData.drawnLines[line] = true
	}
	brk.output.toSVG(smf, line, mode)
}

// --------------------------------------------------------------------------
// Helpers:
// --------------------------------------------------------------------------
func breakHorizontalValues(x0, num int) *drawData {
	return newDrawData(x0, breakWidth(num))
}

func breakVerticalValues(d *drawData, ymax, maxLines int) {
	d.y0 = ymax - LineHeight
	d.minLine = maxLines - 1
	d.height = LineHeight
	d.lines = 1
}

func breakWidth(num int) int {
	return BreakWidth + numWidth(num)
}

// numWidth returns the width to the given number.
// It panics for negative numbers and numbers bigger than 9,999.
func numWidth(num int) int {
	if num < 0 {
		panic("unable to calculate the width of a negative number")
	}

	if num < 10 {
		return CharWidth
	}
	if num < 100 {
		return 2 * CharWidth
	}
	if num < 1000 {
		return 3 * CharWidth
	}
	if num < 10_000 {
		return 4 * CharWidth
	}
	panic("unable to calculate the width of a number bigger than 9,999")
}

func breakToSVG(smf *svgMDFlow, line int, mode FlowMode, bd *drawData, number int) {
	var svg *svgFlow

	// get or create correct SVG flow:
	if mode == FlowModeSVGLinks {
		svg, _ = addNewSVGFlow(smf,
			bd.x0, bd.y0, bd.height, bd.width,
			"sequel", line,
		)
	} else {
		svg = smf.svgs[""]
	}

	svg.Texts = append(svg.Texts, &svgText{
		X:     bd.x0,
		Y:     bd.y0 + bd.height - arrTextOffset,
		Width: bd.width,
		Text:  BreakText + strconv.Itoa(number),
	})

	smf.lastX += bd.width
}
