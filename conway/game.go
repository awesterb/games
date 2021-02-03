package conway

import (
	"strings"
)

// Represents a short Conway game (see his On Numbers and Games.)
type Game interface {
	L() []Game
	R() []Game
}

// Below(g,h) returns g<=h, that is, whether h dominates g as option for Left.
func Below(g, h Game) bool {
	for _, hr := range h.R() {
		if !Gift(g, hr) {
			return false
		}
	}
	for _, gl := range g.L() {
		if !Gift(gl, h) {
			return false
		}
	}
	return true
}

// Gift(g,h) returns g<||h, that is, whether g is for Left a gift horse for h.
func Gift(g, h Game) bool {
	for _, gr := range g.R() {
		if Below(gr, h) {
			return true
		}
	}
	for _, hl := range h.L() {
		if Below(g, hl) {
			return true
		}
	}
	return false
}

// Eq(g,h) return whether g and h are equivalent, that is,  whether g<=h<=g.
func Eq(g, h Game) bool {
	return Below(g, h) && Below(h, g)
}

func Describe(g Game) string {
	sb := strings.Builder{}
	sb.WriteString("{")
	for n, l := range g.L() {
		if n > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(Describe(l))
	}
	sb.WriteString("|")
	for n, r := range g.R() {
		if n > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(Describe(r))
	}
	sb.WriteString("}")
	return sb.String()
}
