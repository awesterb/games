package conway

type Gameish interface {
	Left() []Gameish
	Right() []Gameish
}

func Unfold(g Gameish) *Game {
	L := g.Left()
	ls := make([]*Game, 0, len(L))
	for _, l := range L {
		ls = append(ls, Unfold(l))
	}

	R := g.Right()
	rs := make([]*Game, 0, len(R))
	for _, r := range R {
		rs = append(rs, Unfold(r))
	}
	return &Game{R: rs, L: ls}
}

func (u *Game) Left() []Gameish {
	res := make([]Gameish, 0, len(u.L))
	for _, l := range u.L {
		res = append(res, l)
	}
	return res
}

func (u *Game) Right() []Gameish {
	res := make([]Gameish, 0, len(u.R))
	for _, r := range u.R {
		res = append(res, r)
	}
	return res
}
