package conway

func Simplify(g *Game) *Game {
	// these will be the l and r of the simple we'll return
	var ls, rs []*Game

	var add_l, add_r func(*Game)

	add_l = func(opt *Game) {
		// There is no need to add opt when it is dominated by
		// one of the existing options for Left.
		for _, l := range ls {
			if Below(opt, l) {
				return
			}
		}

		// Let us see if opt is reversible.
		for _, r0 := range opt.R {
			if Below(r0, g) {
				// opt is reversible through r0,
				// so we can add the r0.ls instead of opt.
				for _, l := range r0.L {
					add_l(l)
				}

				return
			}
		}

		// We must add opt, and can remove all options dominated by it.
		new_ls := make([]*Game, 0, len(ls)+1)
		new_ls = append(new_ls, opt)
		for _, l := range ls {
			if !Below(l, opt) {
				new_ls = append(new_ls, l)
			}
		}

		ls = new_ls
	}

	// generalising add_l and add_r is error-prone
	add_r = func(opt *Game) {
		// There is no need to add opt when it is dominated by
		// one of the existing options for Right.
		for _, r := range rs {
			if Below(r, opt) {
				return
			}
		}

		// Let us see if opt is reversible.
		for _, l0 := range opt.L {
			if Below(g, l0) {
				// opt is reversible through l0,
				// so we can add the l0.rs instead of opt.
				for _, r := range l0.R {
					add_r(r)
				}

				return
			}
		}

		// We must add opt, and can remove all options dominated by it.
		new_rs := make([]*Game, 0, len(rs)+1)
		new_rs = append(new_rs, opt)
		for _, r := range rs {
			if !Below(opt, r) {
				new_rs = append(new_rs, r)
			}
		}

		rs = new_rs
	}

	for _, l := range g.L {
		add_l(Simplify(l))
	}

	for _, r := range g.R {
		add_r(Simplify(r))
	}

	return &Game{L: ls, R: rs}
}
