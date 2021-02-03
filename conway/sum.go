package conway

func add(gs ...*unfolded) *unfolded {
	var ls, rs []*unfolded

	var old *unfolded

	for n, g := range gs {
		old = gs[n]
		for _, l := range g.ls {
			gs[n] = l
			ls = append(ls, add(gs...))
		}
		for _, r := range g.rs {
			gs[n] = r
			rs = append(rs, add(gs...))
		}
		gs[n] = old
	}

	return &unfolded{ls: ls, rs: rs}
}

func Add(gs ...Game) Game {
	//return &sum{gs: gs}
	ugs := make([]*unfolded, 0, len(gs))
	for _, g := range gs {
		ugs = append(ugs, unfold(g))
	}
	return add(ugs...)
}

type sum struct {
	gs []Game
}

func (s *sum) L() (res []Game) {
	for n, g := range s.gs {
		for _, l := range g.L() {
			summands := make([]Game, 0, len(s.gs))
			summands = append(summands, s.gs[0:n]...)
			summands = append(summands, l)
			summands = append(summands, s.gs[n+1:len(s.gs)]...)
			res = append(res, Add(summands...))
		}
	}
	return
}

func (s *sum) R() (res []Game) {
	for n, g := range s.gs {
		for _, r := range g.R() {
			summands := make([]Game, 0, len(s.gs))
			summands = append(summands, s.gs[0:n]...)
			summands = append(summands, r)
			summands = append(summands, s.gs[n+1:len(s.gs)]...)
			res = append(res, Add(summands...))
		}
	}
	return
}
