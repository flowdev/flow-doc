package draw2

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"
)

type svgColors struct {
	Background string
	Text       string
	Link       string
	GoLink     string
	Comp       string
	Plugin     string
	PluginType string
}

var lightColors = svgColors{
	Background: "rgb(255,255,255)",
	Text:       "rgb(0,0,0)",
	Link:       "rgb(32,48,128)",
	GoLink:     "rgb(0,96,0)",
	Comp:       "rgb(96,192,255)",
	Plugin:     "rgb(224,224,32)",
	PluginType: "rgb(32,224,32)",
}

var darkColors = svgColors{
	Background: "rgb(13,17,23)",
	Text:       "rgb(201,209,217)",
	Link:       "rgb(96,192,255)",
	GoLink:     "rgb(32,224,32)",
	Comp:       "rgb(32,48,128)",
	Plugin:     "rgb(96,96,0)",
	PluginType: "rgb(0,96,0)",
}

const svgDiagram = `<?xml version="1.0" ?>
<svg version="1.1" xmlns="http://www.w3.org/2000/svg" viewBox="{{.X0}} {{.Y0}} {{.TotalWidth}} {{.TotalHeight}}" width="{{.TotalWidth}}px" height="{{.TotalHeight}}px">
    <!-- Generated by FlowDev tool. -->
    <rect fill="{{.Colors.Background}}" fill-opacity="1" width="{{.TotalWidth}}" height="{{.TotalHeight}}" x="{{.X0}}" y="{{.Y0}}"/>
{{$colors := .Colors}}
{{- range .Arrows}}
    <line stroke="{{$colors.Text}}" stroke-opacity="1.0" stroke-width="2" x1="{{.X1}}" y1="{{.Y1}}" x2="{{.X2}}" y2="{{.Y2}}"/>
    <line stroke="{{$colors.Text}}" stroke-opacity="1.0" stroke-width="2" x1="{{.XTip1}}" y1="{{.YTip1}}" x2="{{.X2}}" y2="{{.Y2}}"/>
    <line stroke="{{$colors.Text}}" stroke-opacity="1.0" stroke-width="2" x1="{{.XTip2}}" y1="{{.YTip2}}" x2="{{.X2}}" y2="{{.Y2}}"/>
{{end -}}
{{- range .Rects}}
{{- if .SubRect}}
    <rect fill="{{$colors.PluginType}}" fill-opacity="1.0" stroke="{{$colors.Text}}" stroke-opacity="1.0" stroke-width="1" width="{{.Width}}" height="{{.Height}}" x="{{.X}}" y="{{.Y}}" rx="10"/>
{{- else -}}
    {{- if .Plugin}}
    <rect fill="{{$colors.Plugin}}" fill-opacity="1.0" stroke="{{$colors.Text}}" stroke-opacity="1.0" stroke-width="2" width="{{.Width}}" height="{{.Height}}" x="{{.X}}" y="{{.Y}}" rx="10"/>
    {{- else}}
    <rect fill="{{$colors.Comp}}" fill-opacity="1.0" stroke="{{$colors.Text}}" stroke-opacity="1.0" stroke-width="2" width="{{.Width}}" height="{{.Height}}" x="{{.X}}" y="{{.Y}}" rx="10"/>
    {{- end}}
{{- end}}
{{- end -}}
{{- if .Texts}}
{{end -}}
{{- range .Texts}}
{{- if .Small}}
    {{- if .GoLink}}
    <text fill="{{$colors.GoLink}}" fill-opacity="1.0" font-size="14" x="{{.X}}" y="{{.Y}}" textLength="{{.Width}}" lengthAdjust="spacingAndGlyphs">{{.Text}}</text>
    {{- else if .Link}}
    <text fill="{{$colors.Link}}" fill-opacity="1.0" font-size="14" x="{{.X}}" y="{{.Y}}" textLength="{{.Width}}" lengthAdjust="spacingAndGlyphs">{{.Text}}</text>
    {{- else}}
    <text fill="{{$colors.Text}}" fill-opacity="1.0" font-size="14" x="{{.X}}" y="{{.Y}}" textLength="{{.Width}}" lengthAdjust="spacingAndGlyphs">{{.Text}}</text>
    {{- end}}
{{- else}}
    {{- if .GoLink}}
    <text fill="{{$colors.GoLink}}" fill-opacity="1.0" font-size="16" x="{{.X}}" y="{{.Y}}" textLength="{{.Width}}" lengthAdjust="spacingAndGlyphs">{{.Text}}</text>
    {{- else if .Link}}
    <text fill="{{$colors.Link}}" fill-opacity="1.0" font-size="16" x="{{.X}}" y="{{.Y}}" textLength="{{.Width}}" lengthAdjust="spacingAndGlyphs">{{.Text}}</text>
    {{- else}}
    <text fill="{{$colors.Text}}" fill-opacity="1.0" font-size="16" x="{{.X}}" y="{{.Y}}" textLength="{{.Width}}" lengthAdjust="spacingAndGlyphs">{{.Text}}</text>
    {{- end}}
{{- end}}
{{- end -}}
{{- if or .Texts .Rects}}
{{end -}}
</svg>
`

var svgTmpl = template.Must(template.New("svgDiagram").Parse(svgDiagram))

const mdDiagram = `
{{- if .FlowLines}}
{{- $maxLine := .MaxLine -}}
{{range $i, $flowLine := .FlowLines}}
    {{- range $cell := $flowLine -}}
        {{- if $cell.Link -}}
            [![{{$cell.Name}}]({{$cell.SVG}})]({{$cell.Link}})
        {{- else -}}
            ![{{$cell.Name}}]({{$cell.SVG}})
        {{- end -}}
    {{- end -}}
    {{- if ne $i $maxLine}}\{{end}}
{{end}}
{{- else}}
![{{.Flow.Name}}]({{.Flow.SVG}})
{{- end}}
{{- if .DataTypes}}

#### Data Types
{{range $name, $link := .DataTypes}}[{{$name}}]({{$link}}), {{end}}
{{end}}
{{- if .Subflows}}

#### Subflows
{{range $name, $link := .Flows}}[{{$name}}]({{$link}}), {{end}}
{{end}}
{{- if .GoFuncs}}

#### Go Functions and Methods
{{range $name, $link := .GoFuncs}}[{{$name}}]({{$link}}), {{end}}
{{end}}
`

var mdTmpl = template.Must(template.New("mdDiagram").Parse(mdDiagram))

const (
	bigDiagramSize  = 256
	tinyDiagramSize = 8
)

type svgArrow struct {
	X1, Y1       int
	X2, Y2       int
	XTip1, YTip1 int
	XTip2, YTip2 int
}

type svgRect struct {
	X, Y    int
	Height  int
	Width   int
	Plugin  bool
	SubRect bool
}

type svgText struct {
	X, Y   int
	Width  int
	Text   string
	Small  bool
	Link   bool
	GoLink bool
}

type svgFlow struct {
	X0, Y0      int
	TotalHeight int
	TotalWidth  int
	Arrows      []*svgArrow
	Rects       []*svgRect
	Texts       []*svgText
	Colors      svgColors
}

func newSVGFlow(x0, y0, height, width, size int) *svgFlow {
	return &svgFlow{
		X0:          x0,
		Y0:          y0,
		TotalHeight: height,
		TotalWidth:  width,
		Arrows:      make([]*svgArrow, 0, size),
		Rects:       make([]*svgRect, 0, size),
		Texts:       make([]*svgText, 0, size),
	}
}

type svgLink struct {
	Name string
	SVG  string
	Link string
}

type mdFlow struct {
	Flow      svgLink
	FlowLines [][]*svgLink
	MaxLine   int
	DataTypes map[string]string
	Subflows  map[string]string
	GoFuncs   map[string]string
}

func newMDFlow() *mdFlow {
	return &mdFlow{
		FlowLines: make([][]*svgLink, 0, 128),
		DataTypes: make(map[string]string, 256),
		Subflows:  make(map[string]string, 256),
		GoFuncs:   make(map[string]string, 256),
	}
}

type svgMDFlow struct {
	svgs          map[string]*svgFlow
	md            *mdFlow
	svgFilePrefix string
	lastX         int
}

func flowToSVGs(f *Flow) *svgMDFlow {
	smf := &svgMDFlow{
		svgs:          make(map[string]*svgFlow, 256),
		md:            newMDFlow(),
		svgFilePrefix: filepath.Join(".", "flowdev", "flow-"+f.name),
	}
	fd := f.getDrawData()

	if f.mode != FlowModeSVGLinks {
		smf.md.Flow = svgLink{
			Name: f.name,
			SVG:  smf.svgFilePrefix + ".svg",
		}
		svg := newSVGFlow(0, 0, fd.height, fd.width+1, bigDiagramSize)
		smf.svgs[""] = svg
	}

	minLine := fd.minLine
	maxLine := minLine + fd.lines - 1
	for line := minLine; line <= maxLine; line++ {
		smf.lastX = 0
		f.toSVG(smf, line, f.mode)

		if f.mode == FlowModeSVGLinks &&
			smf.lastX < fd.width {

			addFillerSVG(smf, line, smf.lastX, LineHeight,
				fd.width-smf.lastX)
		}
	}
	return smf
}

func svgFlowsToBytes(sfs map[string]*svgFlow, dark bool) (map[string][]byte, error) {
	sfbs := make(map[string][]byte)
	for key, sf := range sfs {
		bs, err := svgFlowToBytes(sf, dark)
		if err != nil {
			return nil, fmt.Errorf("unable to create SVG file %q: %w", key, err)
		}
		sfbs[key] = bs
	}
	return sfbs, nil
}

func svgFlowToBytes(sf *svgFlow, dark bool) ([]byte, error) {
	buf := bytes.Buffer{}
	if dark {
		sf.Colors = darkColors
	} else {
		sf.Colors = lightColors
	}
	err := svgTmpl.Execute(&buf, sf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func mdFlowToBytes(mdf *mdFlow) ([]byte, error) {
	buf := bytes.Buffer{}
	mdf.MaxLine = len(mdf.FlowLines) - 1
	err := mdTmpl.Execute(&buf, mdf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// --------------------------------------------------------------------------
//
//	U T I L s :
//
// --------------------------------------------------------------------------
func addNewSVGFlow(
	smf *svgMDFlow,
	x0, y0, height, width int,
	compName string, line int,
) (*svgFlow, *svgLink) {

	svg := newSVGFlow(x0, y0, height, width, tinyDiagramSize)
	name := svgFileName(smf, compName, line)
	smf.svgs[name] = svg
	sl := addSVGLinkToMDFlowLines(smf, line, name, compName)
	return svg, sl
}

func addFillerSVG(smf *svgMDFlow, line, x, height, width int) {
	svg := newSVGFlow(0, 0, height, width, tinyDiagramSize)
	name := fmt.Sprintf("%s-filler-%d-%d.svg", smf.svgFilePrefix, width, height)

	smf.svgs[name] = svg
	addSVGLinkToMDFlowLines(smf, line, name, "filler")
}

func addSVGLinkToMDFlowLines(smf *svgMDFlow, line int, svgName, desc string) *svgLink {
	l := len(smf.md.FlowLines)
	if l <= line {
		for i := l; i <= line; i++ { // fill up lines as necessary (should only be one)
			smf.md.FlowLines = append(smf.md.FlowLines, make([]*svgLink, 0, 32))
		}
	}
	sl := &svgLink{
		Name: desc,
		SVG:  svgName,
	}
	smf.md.FlowLines[line] = append(smf.md.FlowLines[line], sl)

	return sl
}

func svgFileName(smf *svgMDFlow, compName string, line int) string {
	idx := 0
	if len(smf.md.FlowLines) > line {
		idx = len(smf.md.FlowLines[line])
	}
	if compName == "" {
		return fmt.Sprintf("%s-%d-%d.svg", smf.svgFilePrefix, idx, line)
	}
	return fmt.Sprintf("%s-%d-%d-%s.svg", smf.svgFilePrefix, idx, line, compName)
}
