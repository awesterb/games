package conway

// Simple represents a game in its simplest form, without any dominated
// and reversible options.
type Simple interface {
	Game
}

type simple struct {
	l, r []*simple
}

func (s *simple) L() []Game {
	res := make([]Game, 0, len(s.l))
	for _, l := range s.l {
		res = append(res, l)
	}
	return res
}

func (s *simple) R() []Game {
	res := make([]Game, 0, len(s.r))
	for _, r := range s.r {
		res = append(res, r)
	}
	return res
}

func Simplify(g Game) Simple {
	s, ok := g.(*simple)
	if ok {
		return s
	}
	ug, ok := g.(*unfolded)
	if !ok {
		ug = unfold(g)
	}
	return simplify(ug)
}

func simplify(g *unfolded) *simple {
	// these will be the l and r of the simple we'll return
	var ls, rs []*simple

	var add_l, add_r func(*simple)

	add_l = func(opt *simple) {
		// There is no need to add opt when it is dominated by
		// one of the existing options for Left.
		for _, l := range ls {
			if Below(opt, l) {
				return
			}
		}

		// Let us see if opt is reversible.
		for _, r0 := range opt.r {
			if Below(r0, g) {
				// opt is reversible through r0,
				// so we can add the r0.ls instead of opt.
				for _, l := range r0.l {
					add_l(l)
				}

				return
			}
		}

		// We must add opt, and can remove all options dominated by it.
		new_ls := make([]*simple, 0, len(ls)+1)
		new_ls = append(new_ls, opt)
		for _, l := range ls {
			if !Below(l, opt) {
				new_ls = append(new_ls, l)
			}
		}

		ls = new_ls
	}

	// generalising add_l and add_r is error-prone
	add_r = func(opt *simple) {
		// There is no need to add opt when it is dominated by
		// one of the existing options for Right.
		for _, r := range rs {
			if Below(r, opt) {
				return
			}
		}

		// Let us see if opt is reversible.
		for _, l0 := range opt.l {
			if Below(g, l0) {
				// opt is reversible through l0,
				// so we can add the l0.rs instead of opt.
				for _, r := range l0.r {
					add_r(r)
				}

				return
			}
		}

		// We must add opt, and can remove all options dominated by it.
		new_rs := make([]*simple, 0, len(rs)+1)
		new_rs = append(new_rs, opt)
		for _, r := range rs {
			if !Below(opt, r) {
				new_rs = append(new_rs, r)
			}
		}

		rs = new_rs
	}

	for _, l := range g.ls {
		add_l(simplify(l))
	}

	for _, r := range g.rs {
		add_r(simplify(r))
	}

	return &simple{l: ls, r: rs}

}
