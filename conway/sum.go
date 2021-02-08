package conway

func Add(gs ...*Game) *Game {
	var ls, rs []*Game

	var old *Game

	for n, g := range gs {
		old = gs[n]
		for _, l := range g.L {
			gs[n] = l
			ls = append(ls, Add(gs...))
		}
		for _, r := range g.R {
			gs[n] = r
			rs = append(rs, Add(gs...))
		}
		gs[n] = old
	}

	return &Game{L: ls, R: rs}
}
