package types

type Flip string

const (
	FlipBook      Flip = "book"
	FlipLeftHand  Flip = "left-hand"
	FlipCalendar  Flip = "calendar"
	FlipRightHand Flip = "right-hand"
)

// Heteroriented will be true if the card need to pivot about a diagonal axis for the front and back to remain upright.
// the negation of this method is always whether the card is homoriented or not.
func (flip Flip) Heteroriented() bool {
	return flip == FlipLeftHand || flip == FlipRightHand
}

func (flip Flip) String() string {
	switch flip {
	case FlipBook:
		return "vertical axis"
	case FlipLeftHand:
		return "diagonal (up-right) axis"
	case FlipCalendar:
		return "horizontal axis"
	case FlipRightHand:
		return "diagonal (down-right) axis"
	default:
		panic("unknown pivot axis")
	}
}
