package doc

type lengthCache struct {
	calculated bool
	length     int
	canBeFlat  bool
}

func (lc lengthCache) flatLength(calculate func() (int, bool)) (int, bool) {
	if lc.calculated {
		return lc.length, lc.canBeFlat
	}

	lc.length, lc.canBeFlat = calculate()
	lc.calculated = true

	return lc.length, lc.canBeFlat
}
