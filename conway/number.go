package conway

func Int(n int) *Game {
	result := &Game{}

	if int(n) > 0 {
		result.L = []*Game{Int(n - 1)}
	}
	if int(n) < 0 {
		result.R = []*Game{Int(n + 1)}
	}

	return result
}
