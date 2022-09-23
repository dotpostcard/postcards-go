package types

type PivotAxis uint

const (
	PivotAxisUp PivotAxis = iota
	PivotAxisUpRight
	PivotAxisRight
	PivotAxisDownRight
)

// Heteroriented will be true if the card need to pivot about a diagonal axis for the front and back to remain upright.
// the negation of this method is always whether the card is homoriented or not.
func (pa PivotAxis) Heteroriented() bool {
	return pa == PivotAxisUpRight || pa == PivotAxisDownRight
}

func (pa PivotAxis) String() string {
	switch pa {
	case PivotAxisUp:
		return "vertical axis"
	case PivotAxisUpRight:
		return "diagonal (up-right) axis"
	case PivotAxisRight:
		return "horizontal axis"
	case PivotAxisDownRight:
		return "diagonal (down-right) axis"
	default:
		panic("unknown pivot axis")
	}
}
