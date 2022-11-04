package draw

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"
)

const svgDiagram = `<?xml version="1.0" ?>
<svg version="1.1" xmlns="http://www.w3.org/2000/svg" width="{{.TotalWidth}}px" height="{{.TotalHeight}}px">
    <!-- Generated by FlowDev tool. -->
	<rect fill="rgb(255,255,255)" fill-opacity="1" stroke="none" stroke-opacity="1" stroke-width="0" width="{{.TotalWidth}}" height="{{.TotalHeight}}" x="0" y="0"/>
{{- if .Arrows}}
{{end}}
{{- range .Arrows}}
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2" x1="{{.X1}}" y1="{{.Y1}}" x2="{{.X2}}" y2="{{.Y2}}"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2" x1="{{.XTip1}}" y1="{{.YTip1}}" x2="{{.X2}}" y2="{{.Y2}}"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2" x1="{{.XTip2}}" y1="{{.YTip2}}" x2="{{.X2}}" y2="{{.Y2}}"/>
{{end}}
{{- range .Rects}}
{{- if .IsSubRect}}
	<rect fill="rgb(32,224,32)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="1" width="{{.Width}}" height="{{.Height}}" x="{{.X}}" y="{{.Y}}" rx="10"/>
{{- else -}}
	{{- if .IsPlugin}}
	<rect fill="rgb(224,224,32)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2" width="{{.Width}}" height="{{.Height}}" x="{{.X}}" y="{{.Y}}" rx="10"/>
	{{- else}}
	<rect fill="rgb(96,192,255)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2" width="{{.Width}}" height="{{.Height}}" x="{{.X}}" y="{{.Y}}" rx="10"/>
	{{- end}}
{{- end}}
{{- end}}
{{range .Texts}}
{{- if .Small}}
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-size="14" x="{{.X}}" y="{{.Y}}" textLength="{{.Width}}" lengthAdjust="spacingAndGlyphs">{{.Text}}</text>
{{- else}}
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-size="16" x="{{.X}}" y="{{.Y}}" textLength="{{.Width}}" lengthAdjust="spacingAndGlyphs">{{.Text}}</text>
{{- end}}
{{- end}}
</svg>
`

const mdDiagram = `
{{- if .FlowLines}}
{{- $n := len .FlowLines -}}
{{range $i, $flowLine := .FlowLines}}
	{{- range $cell := $flowLine -}}
		{{- if $cell.Link -}}
			[![{{$cell.Name}}]({{$cell.SVG}})]({{$cell.Link}})
		{{- else -}}
			![{{$cell.Name}}]({{$cell.SVG}})
		{{- end -}}
	{{- end -}}
	{{- if ne $i $n}}\{{end}}
{{end}}
{{else}}
![{{.Flow.Name}}]({{.Flow.SVG}})
{{end}}
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

type TextType int

const (
	TextTypeText TextType = iota
	TextTypeLink
	TextTypeGoLink
)

type svgArrow struct {
	X1, Y1       int
	X2, Y2       int
	XTip1, YTip1 int
	XTip2, YTip2 int
}

type svgRect struct {
	X, Y      int
	Width     int
	Height    int
	IsPlugin  bool
	IsSubRect bool
}

type svgText struct {
	X, Y  int
	Width int
	Text  string
	Small bool
	Type  TextType
}

type svgFlow struct {
	X0, Y0      int
	TotalWidth  int
	TotalHeight int
	Arrows      []*svgArrow
	Rects       []*svgRect
	Texts       []*svgText
}

type svgLink struct {
	Name string
	SVG  string
	Link string
}

type mdFlow struct {
	Flow      svgLink
	FlowLines [][]svgLink
	DataTypes map[string]string
	Subflows  map[string]string
	GoFuncs   map[string]string
}

var svgTmpl = template.Must(template.New("svgDiagram").Parse(svgDiagram))

var mdTmpl = template.Must(template.New("mdDiagram").Parse(mdDiagram))

// FromFlowData creates a set of SVG diagrams and a MarkDown file from flow
// data. If the flow data isn't valid or the SVG diagrams or the MarkDown file
// can't be created with their template, an error is returned.
func FromFlowData(f *Flow) (svgContents map[string][]byte, mdContent []byte, err error) {
	err = validateFlowData(f)
	if err != nil {
		return nil, nil, err
	}

	enrichFlow(f)
	sfs, mdf := flowToSVGs(f)
	if f.Mode != FlowModeSVGLinks {
		sfs["flow-"+f.Name+".svg"] = sfs[""]
		delete(sfs, "")
	}

	svgContents, err = svgFlowsToBytes(sfs)
	if err != nil {
		return nil, nil, err
	}
	mdContent, err = mdFlowToBytes(mdf)
	if err != nil {
		return nil, nil,
			fmt.Errorf("unable to create MarkDown content for %q flow: %w", f.Name, err)
	}
	return svgContents, mdContent, nil
}

func enrichFlow(f *Flow) {
	merges := make(map[string]*Merge)
	enrichSplit(f.Shapes, 0, 0, 0, nil, FlowModeNoLinks, merges,
		enrichArrow, enrichOp, enrichMerge)
}

func flowToSVGs(f *Flow) (map[string]*svgFlow, *mdFlow) {
	sfs := make(map[string]*svgFlow, 256)
	mdf := &mdFlow{
		FlowLines: make([][]svgLink, 0, 128),
		DataTypes: make(map[string]string),
		Subflows:  make(map[string]string),
		GoFuncs:   make(map[string]string),
	}
	if f.Mode != FlowModeSVGLinks {
		mdf.Flow = svgLink{
			Name: f.Name,
			SVG:  filepath.Join(".", "flowdev", "flow-"+f.Name+".svg"),
		}
		d := f.Shapes.drawData
		svg := &svgFlow{
			X0:          0,
			Y0:          0,
			TotalHeight: d.height,
			TotalWidth:  d.width,
		}
		sfs[""] = svg
	}

	minLine := f.Shapes.drawData.minLine
	maxLine := minLine + f.Shapes.drawData.lines - 1
	for line := minLine; line <= maxLine; line++ {
		splitToSVG(sfs, mdf, line, f.Mode, f.Shapes)
	}
	return sfs, mdf
}

func svgFlowsToBytes(sfs map[string]*svgFlow) (map[string][]byte, error) {
	sfbs := make(map[string][]byte)
	for key, sf := range sfs {
		bs, err := svgFlowToBytes(sf)
		if err != nil {
			return nil, fmt.Errorf("unable to create SVG file %q: %w", key, err)
		}
		sfbs[key] = bs
	}
	return sfbs, nil
}

func svgFlowToBytes(sf *svgFlow) ([]byte, error) {
	buf := bytes.Buffer{}
	err := svgTmpl.Execute(&buf, sf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func mdFlowToBytes(mdf *mdFlow) ([]byte, error) {
	buf := bytes.Buffer{}
	err := mdTmpl.Execute(&buf, mdf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
