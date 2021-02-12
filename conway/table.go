package conway

import (
	"errors"
	"log"
	"sort"
	"strconv"
	"strings"
)

type Serial int

func sortSerials(s []Serial) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

func (s Serial) String() string {
	return strconv.Itoa(int(s))
}

var NotSimpleErr = errors.New("game is not simple")

// Table represents a table of all simple games.
type Table interface {
	// Serial looks up the serial of a game
	Serial(g *Game) Serial

	// SerialSimple looks up the code of a simple game.
	SerialSimple(g *Game) (Serial, error)

	// Left returns the left options of the given serial.
	Left(s Serial) []Serial

	// Right returns the right options of the given serial.
	Right(s Serial) []Serial

	// Describe describes the game associated to the given serial
	Describe(s Serial) string
}

func NewTable() Table {
	return &table{
		entries: []*entry{
			&entry{}, // zero
		},
		lookup: map[string]Serial{"|": 0},
		below: map[xy]bool{
			xy{x: 0, y: 0}: true, // 0 <= 0
		},
		gift: map[xy]bool{
			xy{x: 0, y: 0}: false, // not 0 <| 0
		},
		next: 0,
	}
}

// table lists all simple games up to a certain complexity; the index
// of a game in this table is called its serial number.
// The table always contains all games G for which the G^R and G^L
// all have serial number smaller than `next`.
//
// The table is extended in steps (by incrementing `next`) by adding
// the simple games from
//   { H, G^L | G^R },  { H, G^L | G^R, H },  { G^L | G^R, H }   (*)
// where H is the game with serial `next`.
//
// These games (*) are added in the following order:  the G's with lower serial
// number go first, and games with the same G,  are added in the order
// listed in (*).
//
// The first game, with serial number 0 is, of course, 0 = {|}.
//
// Together these stipulations determine the start and extensions of
// the table uniquely:  the first few entries are:
//
// 0,  1={0|},  -1={|0},  *={0|0}
type table struct {
	entries     []*entry
	lookup      map[string]Serial
	below, gift map[xy]bool
	next        Serial
}

type entry struct {
	L, R []Serial
}

func (t *table) ensureContains(s Serial) {
	for len(t.entries) <= int(s) {
		t.extend()
	}
}

func (t *table) Left(s Serial) []Serial {
	t.ensureContains(s)
	return t.entries[s].L
}

func (t *table) Right(s Serial) []Serial {
	t.ensureContains(s)
	return t.entries[s].R
}

func (t *table) Describe(s Serial) string {
	sb := strings.Builder{}
	sb.WriteRune('{')
	for n, l := range t.Left(s) {
		if n > 0 {
			sb.WriteRune(',')
		}
		sb.WriteRune('#')
		sb.WriteString(l.String())
	}
	sb.WriteRune('|')
	for n, r := range t.Right(s) {
		if n > 0 {
			sb.WriteRune(',')
		}
		sb.WriteRune('#')
		sb.WriteString(r.String())
	}
	sb.WriteRune('}')
	return sb.String()
}

// serialish represents a serial (an entry from the table),
// or an entry not yet added to the table whose options are all
// already present in the table.
type serialish interface {
	leftIn(*table) []Serial
	rightIn(*table) []Serial
}

func (e *entry) leftIn(t *table) []Serial  { return e.L }
func (e *entry) rightIn(t *table) []Serial { return e.R }

func (s Serial) leftIn(t *table) []Serial  { return t.entries[s].L }
func (s Serial) rightIn(t *table) []Serial { return t.entries[s].R }

// key returns a key representing the entry to be used in map types;
// entry, containing slices, can not be a key type of a map.
func (e entry) key() string {
	var sb strings.Builder
	for n, l := range e.L {
		if n > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(l.String())
	}
	sb.WriteRune('|')
	for n, r := range e.R {
		if n > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(r.String())
	}
	return sb.String()
}

type xy struct {
	x, y Serial
}

func (t *table) Serial(g *Game) Serial {
	s, err := t.SerialSimple(Simplify(g))
	if err != nil {
		log.Fatal("unexpected error")
	}
	return s
}

func (t *table) SerialSimple(g *Game) (serial Serial, err error) {

	lserials := make([]Serial, 0, len(g.L))
	rserials := make([]Serial, 0, len(g.R))

	for _, l := range g.L {
		var ls Serial
		ls, err = t.SerialSimple(l)
		if err != nil {
			return
		}
		lserials = append(lserials, ls)
	}

	for _, r := range g.R {
		var rs Serial
		rs, err = t.SerialSimple(r)
		if err != nil {
			return
		}
		rserials = append(rserials, rs)
	}

	sortSerials(lserials)
	sortSerials(rserials)

	// extend the table until our game is in it
	for (len(lserials) > 0) && lserials[len(lserials)-1] >= t.next ||
		(len(rserials) > 0) && rserials[len(rserials)-1] >= t.next {

		t.extend()
	}

	key := entry{L: lserials, R: rserials}.key()
	serial, ok := t.lookup[key]
	if !ok {
		err = NotSimpleErr
		return
	}
	return
}

func (t *table) extend() (count int) {
	N := len(t.entries)
	for gs := 0; gs < N; gs++ {
		count += t.extendWith(Serial(gs), true, false)
		count += t.extendWith(Serial(gs), true, true)
		count += t.extendWith(Serial(gs), false, true)
	}
	t.next++
	return count
}

func (t *table) extendWith(gs Serial, addl, addr bool) int {
	ge := t.entries[gs]

	// If by adding t.next as option to ge, we create
	// a dominated option, abort.
	if addl {
		for _, l := range ge.L {
			if t.Below(Serial(l), t.next) ||
				t.Below(t.next, l) {
				return 0
			}
		}
	}
	if addr {
		for _, r := range ge.R {
			if t.Below(Serial(r), t.next) ||
				t.Below(t.next, r) {
				return 0
			}
		}
	}

	gem := &entry{
		L: make([]Serial, 0, len(ge.L)+1),
		R: make([]Serial, 0, len(ge.R)+1),
	}
	gem.L = append(gem.L, ge.L...)
	if addl {
		gem.L = append(gem.L, t.next)
		sortSerials(gem.L)
	}
	gem.R = append(gem.R, ge.R...)
	if addr {
		gem.R = append(gem.R, t.next)
		sortSerials(gem.R)
	}

	// Check if gem has reversible options, and if so, abort.
	for _, l := range gem.L {
		for _, lr := range t.entries[l].R {
			if t.belowish(lr, gem) {
				// l is reversible through lr
				return 0
			}
		}
	}

	for _, r := range gem.L {
		for _, rl := range t.entries[r].L {
			if t.belowish(gem, rl) {
				// r is reversible through rl
				return 0
			}
		}
	}

	// gem is simple, so can be added to the table

	serial := Serial(len(t.entries))
	t.entries = append(t.entries, gem)
	t.lookup[gem.key()] = serial

	for i := 0; i <= int(serial); i++ {
		t.setBelow(Serial(i), serial)
		t.setBelow(serial, Serial(i))
		t.setGift(Serial(i), serial)
		t.setGift(serial, Serial(i))
	}

	return 1
}

func (t *table) setBelow(g, h Serial) {
	t.below[xy{x: g, y: h}] = func() bool {
		for _, hr := range t.entries[h].R {
			if !t.Gift(g, hr) {
				return false
			}
		}
		for _, gl := range t.entries[g].L {
			if !t.Gift(gl, h) {
				return false
			}
		}
		return true
	}()
}

func (t *table) setGift(g, h Serial) {
	t.gift[xy{x: g, y: h}] = func() bool {
		for _, gr := range t.entries[h].R {
			if t.Below(gr, h) {
				return true
			}
		}
		for _, hl := range t.entries[h].L {
			if t.Below(g, hl) {
				return true
			}
		}
		return false
	}()
}

func (t *table) Below(g, h Serial) bool {
	return t.below[xy{x: g, y: h}]
}

func (t *table) Gift(g, h Serial) bool {
	return t.gift[xy{x: g, y: h}]
}

func (t *table) belowish(g, h serialish) bool {
	gs, okg := g.(Serial)
	hs, okh := h.(Serial)

	if okg && okh {
		return t.Below(gs, hs)
	}

	for _, hr := range h.rightIn(t) {
		if !t.giftish(g, hr) {
			return false
		}
	}
	for _, gl := range g.leftIn(t) {
		if !t.giftish(gl, h) {
			return false
		}
	}
	return true
}

func (t *table) giftish(g, h serialish) bool {
	gs, okg := g.(Serial)
	hs, okh := h.(Serial)

	if okg && okh {
		return t.Gift(gs, hs)
	}

	for _, gr := range g.rightIn(t) {
		if t.belowish(gr, h) {
			return true
		}
	}
	for _, hl := range h.leftIn(t) {
		if t.belowish(g, hl) {
			return true
		}
	}
	return false
}
