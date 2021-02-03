package conway

// Unfolds unfolds a given game, computing all its positions, and storing
// them in memory, returning an equivalent game, for which L() and R()
// are swiftly computable.
func Unfold(g Game) Game {
	return unfold(g)
}

func unfold(g Game) *unfolded {
	L := g.L()
	ls := make([]*unfolded, 0, len(L))
	for _, l := range L {
		ls = append(ls, unfold(l))
	}

	R := g.R()
	rs := make([]*unfolded, 0, len(R))
	for _, r := range R {
		rs = append(rs, unfold(r))
	}
	return &unfolded{rs: rs, ls: ls}
}

type unfolded struct {
	ls, rs []*unfolded
}

func (u *unfolded) L() []Game {
	res := make([]Game, 0, len(u.ls))
	for _, l := range u.ls {
		res = append(res, l)
	}
	return res
}

func (u *unfolded) R() []Game {
	res := make([]Game, 0, len(u.rs))
	for _, r := range u.rs {
		res = append(res, r)
	}
	return res
}
