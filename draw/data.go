package draw

import (
	"fmt"
)

const (
	RowGap     = 8
	LineHeight = 24
	CharWidth  = 8
	ParenWidth = 6
	WordGap    = 6
	TextOffset = 6
	BreakText  = "…"
	BreakWidth = 2 * ParenWidth
	LoopText   = "back to: "
	LoopWidth  = 6*CharWidth + 3*CharWidth/2
)

type FlowMode int

const (
	FlowModeNoLinks FlowMode = iota
	FlowModeMDLinks
	FlowModeSVGLinks // not implemented yet
)

// Shapes:
// Name     | Start | Mid | End
// ============================
// Comp     | Yes   | Yes | Yes
// StartPort| Yes   | No  | No
// BreakEnd | Yes   | No  | No
// BreakStart No    | No  | Yes
// EndPort  | No    | No  | Yes
// Loop     | No    | No  | Yes
// ----------------------------
// Arrow    | No    | Yes | No

type anyComp interface {
	calcHorizontalValues(x0 int)

	// maxWidth is constant and newWidth the full width (x0 + width)
	respectMaxWidth(maxWidth, num int) (newStartComps []StartComp, newNum, newWidth int)
	resetDrawData()

	// maxLines is full number of lines (minLine + lines) and newHeight is the full height (y0 + height)
	calcVerticalValues(y0, minLine int, mode FlowMode) (maxLines, newHeight int)
	toSVG(smf *svgMDFlow, line int, mode FlowMode)
	getDrawData() *drawData
}

type StartComp interface {
	anyComp
}

type EndComp interface {
	anyComp
	addInput(*Arrow)
	switchInput(oldArr, newArr *Arrow)
	minRestOfRowWidth(num int) int
}

// drawData contains all data needed for positioning an element correctly.
type drawData struct {
	x0, y0         int
	height, width  int
	minLine, lines int
	drawnLines     map[int]bool
}

func newDrawData(x0, width int) *drawData {
	return &drawData{
		x0:         x0,
		width:      width,
		drawnLines: make(map[int]bool),
	}
}
func (d *drawData) xmax() int {
	return d.x0 + d.width
}
func (d *drawData) ymax() int {
	return d.y0 + d.height
}
func (d *drawData) maxLines() int {
	return d.minLine + d.lines
}
func (d *drawData) contains(line int) bool {
	return d.minLine <= line && line < d.minLine+d.lines
}
func (d *drawData) drawLine(line int) bool {
	return d.contains(line) && !d.drawnLines[line]
}

type withDrawData struct {
	drawData *drawData
}

func (wd *withDrawData) getDrawData() *drawData {
	return wd.drawData
}
func (wd *withDrawData) resetDrawData() {
	wd.drawData = nil
}

type CompRegistry interface {
	register(*Comp)
	lookup(id string) *Comp
}

type ShapeCluster struct {
	withDrawData
	shapeRows    []StartComp
	compRegistry map[string]*Comp
}

func NewCluster() *ShapeCluster {
	return &ShapeCluster{
		shapeRows:    make([]StartComp, 0, 32),
		compRegistry: make(map[string]*Comp, 128),
	}
}

func (cl *ShapeCluster) AddStartComp(comp StartComp) *ShapeCluster {
	cl.shapeRows = append(cl.shapeRows, comp)
	return cl
}

func (cl *ShapeCluster) lookup(id string) *Comp {
	return cl.compRegistry[id]
}

func (cl *ShapeCluster) register(comp *Comp) {
	cl.compRegistry[comp.ID()] = comp
}

func (cl *ShapeCluster) resetDrawData() {
	for _, cl := range cl.shapeRows {
		cl.resetDrawData()
	}
	cl.withDrawData.resetDrawData()
}

func (cl *ShapeCluster) calcHorizontalValues() {
	for _, comp := range cl.shapeRows {
		comp.calcHorizontalValues(0)
	}
}

func (cl *ShapeCluster) respectMaxWidth(maxWidth, num int) (newNum, newWidth int) {
	var newRows []StartComp
	addRows := make([]StartComp, 0, 64)
	width := 0

	for _, comp := range cl.shapeRows {
		newRows, num, width = comp.respectMaxWidth(maxWidth, num)
		addRows = append(addRows, newRows...)
		newWidth = max(newWidth, width)
	}

	for len(addRows) > 0 {
		n := len(cl.shapeRows)
		cl.shapeRows = append(cl.shapeRows, addRows...)
		addRows = addRows[:0]
		for i := n; i < len(cl.shapeRows); i++ {
			start := cl.shapeRows[i]
			start.resetDrawData()
			start.calcHorizontalValues(0)
			newRows, num, width = start.respectMaxWidth(maxWidth, num)
			addRows = append(addRows, newRows...)
			newWidth = max(newWidth, width)
		}
	}
	cl.drawData = newDrawData(0, newWidth)
	return num, newWidth
}

func (cl *ShapeCluster) calcVerticalValues(y0, minLine int, mode FlowMode) (maxLines, newHeight int) {
	cd := cl.drawData
	cd.y0 = y0
	cd.minLine = minLine
	for i, comp := range cl.shapeRows {
		if i > 0 && mode != FlowModeSVGLinks {
			y0 += RowGap
		}
		minLine, y0 = comp.calcVerticalValues(y0, minLine, mode)
	}
	cl.drawData.lines = minLine - cd.minLine
	cl.drawData.height = y0 - cd.y0
	return minLine, y0
}

func (cl *ShapeCluster) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	if !cl.drawData.drawLine(line) {
		return
	}
	for i := len(cl.shapeRows) - 1; i >= 0; i-- {
		cl.shapeRows[i].toSVG(smf, line, mode)
	}
	cl.drawData.drawnLines[line] = true
}

type Flow struct {
	withDrawData
	name     string
	mode     FlowMode
	width    int
	dark     bool
	clusters []*ShapeCluster
}

func NewFlow(name string, mode FlowMode, width int, dark bool) *Flow {
	return &Flow{
		name:     name,
		mode:     mode,
		width:    width,
		dark:     dark,
		clusters: make([]*ShapeCluster, 0, 64),
	}
}

func (flow *Flow) ChangeConfig(name string, mode FlowMode, width int, dark bool) {
	flow.name = name
	flow.mode = mode
	flow.width = width
	flow.dark = dark
}

func (flow *Flow) AddCluster(cl *ShapeCluster) *Flow {
	flow.clusters = append(flow.clusters, cl)
	return flow
}

// Draw creates a set of SVG diagrams and a MarkDown file for this flow.
// If the flow data isn't valid or the SVG diagrams or the MarkDown file
// can't be created with their template, an error is returned.
func (flow *Flow) Draw() (svgContents map[string][]byte, mdContent []byte, err error) {
	err = flow.validate()
	if err != nil {
		return nil, nil, err
	}

	flow.resetDrawData()
	flow.calcHorizontalValues()
	flow.respectMaxWidth()
	flow.calcVerticalValues()

	smf := flowToSVGs(flow)
	if flow.mode != FlowModeSVGLinks {
		svgName := smf.svgFilePrefix + ".svg"
		smf.svgs[svgName] = smf.svgs[""]
		delete(smf.svgs, "")
		smf.md.FlowLines = append(smf.md.FlowLines, make([]*svgLink, 1))
		smf.md.FlowLines[0][0] = &svgLink{
			Name: flow.name,
			SVG:  svgName,
		}
	}

	svgContents, err = svgFlowsToBytes(smf.svgs, flow.dark)
	if err != nil {
		return nil, nil, err
	}
	mdContent, err = mdFlowToBytes(smf.md)
	if err != nil {
		return nil, nil,
			fmt.Errorf("unable to create MarkDown content for %q flow: %w", flow.name, err)
	}
	return svgContents, mdContent, nil
}

func (flow *Flow) validate() error {
	if len(flow.clusters) == 0 {
		return fmt.Errorf("no shape clusters found in flow")
	}

	for i, cl := range flow.clusters {
		if len(cl.shapeRows) == 0 {
			return fmt.Errorf("no shapes found in the %d-th cluster of the flow", i+1)
		}
	}

	return nil
}

func (flow *Flow) resetDrawData() {
	for _, cl := range flow.clusters {
		cl.resetDrawData()
	}
	flow.withDrawData.resetDrawData()
}

func (flow *Flow) calcHorizontalValues() {
	for _, cl := range flow.clusters {
		cl.calcHorizontalValues()
	}
}

func (flow *Flow) respectMaxWidth() {
	num := 1 // breaks start with 1
	width, maxWidth := 0, 0
	for _, cl := range flow.clusters {
		num, width = cl.respectMaxWidth(flow.width, num)
		maxWidth = max(maxWidth, width)
	}
	flow.drawData = &drawData{
		x0:      0,
		y0:      0,
		width:   maxWidth,
		minLine: 0,
	}
}

func (flow *Flow) calcVerticalValues() {
	height, lines := 0, 0
	for i, cl := range flow.clusters {
		if i > 0 && flow.mode != FlowModeSVGLinks {
			height += LineHeight - RowGap
		}
		lines, height = cl.calcVerticalValues(height, lines, flow.mode)
	}
	flow.drawData.height = height
	flow.drawData.lines = lines
}

func (flow *Flow) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	for _, cl := range flow.clusters {
		cl.toSVG(smf, line, mode)
	}
}
