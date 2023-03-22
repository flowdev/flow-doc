package draw2

type Loop struct {
	name       string
	port       string
	link       string
	input      *Arrow
	inReturned bool
	goLink     bool
	drawData   *drawData
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

func (loop *Loop) prevArrow() *Arrow {
	if loop.input == nil {
		return nil
	}

	loop.inReturned = !loop.inReturned
	if loop.inReturned {
		return loop.input
	}
	return nil
}

func (loop *Loop) minRestOfRowWidth(num int) int {
	return loop.drawData.width
}

func (loop *Loop) intersects(line int) bool {
	return withinShape(line, loop.drawData)
}

func (loop *Loop) respectMaxWidth(maxWidth, num int) ([]StartComp, int) {
	return nil, num
}

func (loop *Loop) calcHorizontalValues(x0 int) {
	txt := loop.name + loop.port
	width := BreakWidth + LoopWidth + len(txt)*CharWidth
	if loop.port != "" {
		width += CharWidth / 2
	}

	loop.drawData = &drawData{
		x0:    x0,
		width: width,
	}
}

func (loop *Loop) calcVerticalValues(y0, minLine int, mode FlowMode) {
	ld := loop.drawData
	ld.y0 = y0
	ld.minLine = minLine
	ld.height = LineHeight
	ld.lines = 1
}
