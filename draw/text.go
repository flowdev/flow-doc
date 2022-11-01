package draw

func textDataToSVG(t Text, sf *svgFlow, x int, y int) (nsf *svgFlow, nx, ny int) {
	txt := "... back to: " + t.Text
	width := len(txt) * 8

	y += 12 + 24 - 6
	sf.Texts = append(sf.Texts, &svgText{
		X: x, Y: y,
		Width: width,
		Text:  txt,
	})

	return sf, x + width + 8, y + 12
}
