package conway

type Int int

func (i Int) L() []Game {
	if int(i) > 0 {
		return []Game{Int(i - 1)}
	}
	return []Game{}
}

func (i Int) R() []Game {
	if int(i) < 0 {
		return []Game{Int(i + 1)}
	}
	return []Game{}
}
