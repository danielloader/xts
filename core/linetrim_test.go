package core

import (
	"strings"
	"testing"
)

// trimLayout is a 108pt band (nine 12pt rows above a 113.39pt content area)
// filled with lines of 15.6pt, of which about 2.2pt below the text are
// leading: seven lines need 109.2pt, but the seventh line's text ends at
// about 107pt.
func trimLayout(style, body string) string {
	return layoutHead + `
  <PageFormat width="100mm" height="60mm"/>
  <SetGrid width="5mm" height="12pt"/>
  <StyleSheet>p { margin: 0; font-family: serif; font-size: 10pt; line-height: 15.6pt; ` + style + ` }</StyleSheet>
  <Record match="data">` + body + `</Record>
</Layout>`
}

func paragraphs(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString(`<Paragraph><Value>Line</Value></Paragraph>`)
	}
	return b.String()
}

// In a Flow, text-box-trim: trim-end lets the band's last line fit by its
// text, and the leading below it lies past the band without the warning that
// an object protrudes into the bottom margin.
func TestTextBoxTrimInAFlow(t *testing.T) {
	for _, c := range []struct {
		style string
		lines int
	}{
		{``, 6},
		{`text-box-trim: trim-end`, 7},
	} {
		log := runLayoutLog(t, trimLayout(c.style, `<Flow>`+paragraphs(10)+`</Flow><Message select="concat('page ', sd:current-page())"/>`))
		if strings.Contains(log, "protrudes into the bottom margin") {
			t.Errorf("%q: the trimmed leading is reported as protruding", c.style)
		}
		root := renderDump(t, trimLayout(c.style, `<Flow>`+paragraphs(10)+`</Flow>`))
		if got := strings.Count(dumpText(root.Children[0]), "Line"); got != c.lines {
			t.Errorf("%q: page 1 holds %d lines, want %d", c.style, got, c.lines)
		}
	}
}

// A PlaceObject with a text block fits where its last line's text does, when
// that line is trimmed: it is placed at the top without the warning.
func TestTextBoxTrimInAPlaceObject(t *testing.T) {
	for _, c := range []struct {
		style string
		warns bool
	}{
		{``, true},
		{`text-box-trim: trim-end`, false},
	} {
		log := runLayoutLog(t, trimLayout(c.style, `<PlaceObject><TextBlock>`+paragraphs(7)+`</TextBlock></PlaceObject>`))
		if got := strings.Contains(log, "protrudes into the bottom margin"); got != c.warns {
			t.Errorf("%q: warned %v, want %v", c.style, got, c.warns)
		}
	}
}

// dumpText is the glyphs of a dump node and its children, in order.
func dumpText(n dumpNode) string {
	s := n.attr("components")
	for _, c := range n.Children {
		s += dumpText(c)
	}
	return s
}
