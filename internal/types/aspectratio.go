package types

import "fmt"

// AspectRatio is the width over the height; larger if wider than tall
type AspectRatio float64

func NewAspectRatio(w, h float64) AspectRatio {
	return AspectRatio(w / h)
}

func (ar AspectRatio) String() string {
	return fmt.Sprintf("1:%.3f", ar)
}
