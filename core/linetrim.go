package core

import (
	"github.com/boxesandglue/boxesandglue/backend/bag"
	"github.com/boxesandglue/boxesandglue/backend/node"
	"github.com/boxesandglue/htmlbag"
)

// trailingTrim is how far the line vl ends with may reach below a frame:
// the leading below its text that htmlbag records for a block with
// text-box-trim: trim-end (htmlbag.LineTrimEnd). It follows the chain of
// boxes a text block comes in down to its last line; anything with height
// after that line, a border or padding below the text, leaves none.
func trailingTrim(vl *node.VList) bag.ScaledPoint {
	var cur node.Node = vl
	for depth := 0; depth <= 12; depth++ {
		switch t := cur.(type) {
		case *node.HList:
			if trim, _ := t.Attributes[htmlbag.LineTrimEnd].(bag.ScaledPoint); trim > 0 {
				return trim
			}
			// A wrapper, or a text block's line holding the HTML: only its
			// sole box sets its depth.
			var sole node.Node
			for c := t.List; c != nil; c = c.Next() {
				switch c.(type) {
				case *node.VList, *node.HList:
					if sole != nil {
						return 0
					}
					sole = c
				case *node.Glue, *node.Kern, *node.StartStop, *node.Penalty:
				default:
					return 0
				}
			}
			if sole == nil {
				return 0
			}
			cur = sole
		case *node.VList:
			last := node.Tail(t.List)
		back:
			for last != nil {
				switch n := last.(type) {
				case *node.Glue:
					if n.Width != 0 {
						return 0
					}
				case *node.Kern:
					if n.Kern != 0 {
						return 0
					}
				case *node.StartStop, *node.Penalty:
				default:
					break back
				}
				last = last.Prev()
			}
			if last == nil {
				return 0
			}
			cur = last
		default:
			return 0
		}
	}
	return 0
}
