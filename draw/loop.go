package draw

type Loop struct {
	withDrawData
	name   string
	port   string
	link   string
	goLink bool
	input  *Arrow
}

func NewLoop(name, port, link string) *Loop {
	return &Loop{
		name: name,
		port: port,
		link: link,
	}
}

func (loop *Loop) GoLink() *Loop {
	loop.goLink = true
	return loop
}

func (loop *Loop) addInput(arr *Arrow) {
	loop.input = arr
}

func (loop *Loop) switchInput(oldArr, newArr *Arrow) {
	newArr.dstComp = loop
	loop.input = newArr
}

func (loop *Loop) minRestOfRowWidth(num int) int {
	return loop.drawData.width
}

func (loop *Loop) calcHorizontalValues(x0 int) {
	txt := loop.name + loop.port
	width := BreakWidth + LoopWidth + len(txt)*CharWidth
	if loop.port != "" {
		width += CharWidth / 2
	}

	loop.withDrawData.drawData = newDrawData(x0, width)
}

func (loop *Loop) respectMaxWidth(maxWidth, num int) (newStartComps []StartComp, newNum, newWidth int) {
	return nil, num, loop.drawData.xmax()
}

func (loop *Loop) calcVerticalValues(y0, minLine int, mode FlowMode) (maxLines, newHeight int) {
	ld := loop.drawData
	ld.y0 = loop.input.drawData.ymax() - LineHeight
	ld.minLine = loop.input.drawData.maxLines() - 1
	ld.height = LineHeight
	ld.lines = 1
	return ld.maxLines(), ld.ymax()
}

func (loop *Loop) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	var svg *svgFlow
	ld := loop.drawData

	if !ld.drawLine(line) {
		return
	}

	// get or create correct SVG flow:
	if mode == FlowModeMDLinks {
		var svgLink *svgLink
		svg, svgLink = addNewSVGFlow(smf,
			ld.x0, ld.y0, ld.height, ld.width,
			"loop", line,
		)
		svgLink.Link = loop.link
	} else {
		svg = smf.svgs[""]
	}

	txt := BreakText + LoopText + loop.name
	if loop.port != "" {
		txt += ":" + loop.port
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:      ld.x0,
		Y:      ld.y0 + ld.height - arrTextOffset,
		Width:  ld.width,
		Text:   txt,
		Link:   !loop.goLink && loop.link != "",
		GoLink: loop.goLink,
	})

	smf.lastX += ld.width
}
