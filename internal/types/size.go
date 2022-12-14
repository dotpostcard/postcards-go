package types

import (
	"fmt"
	"math"
	"math/big"
)

const maxRatioDiff = 0.01

type Size struct {
	CmWidth  *big.Rat `json:"cmW,omitempty"`
	CmHeight *big.Rat `json:"cmH,omitempty"`
	PxWidth  int      `json:"pxW"`
	PxHeight int      `json:"pxH"`
}

func (s Size) HasPhysical() bool {
	return s.CmWidth != nil && s.CmHeight != nil
}

func (s *Size) SetResolution(xRes *big.Rat, yRes *big.Rat) {
	s.CmWidth = big.NewRat(1, 1).Quo(
		big.NewRat(int64(s.PxWidth), 1),
		xRes,
	)
	s.CmHeight = big.NewRat(1, 1).Quo(
		big.NewRat(int64(s.PxHeight), 1),
		yRes,
	)
}

// Resolution returns the pixels per centimeter
func (s Size) Resolution() (xRes *big.Rat, yRes *big.Rat) {
	xRes = big.NewRat(1, 1).Quo(big.NewRat(int64(s.PxWidth), 1), s.CmWidth)
	yRes = big.NewRat(1, 1).Quo(big.NewRat(int64(s.PxHeight), 1), s.CmHeight)
	return
}

func (s Size) SimilarPhysical(other Size, flip Flip) bool {
	if !s.HasPhysical() || !other.HasPhysical() {
		return true
	}

	if flip.Heteroriented() {
		return similar(s.CmWidth, other.CmHeight) && similar(s.CmHeight, other.CmWidth)
	} else {
		return similar(s.CmWidth, other.CmWidth) && similar(s.CmHeight, other.CmHeight)
	}
}

func similar(a, b *big.Rat) bool {
	ratio, _ := big.NewRat(1, 1).Quo(a, b).Float64()
	return math.Abs(1-ratio) <= maxRatioDiff
}

func (s Size) String() string {
	pxSize := fmt.Sprintf("%dpx x %dpx", s.PxWidth, s.PxHeight)
	if !s.HasPhysical() {
		return pxSize
	}

	fw, _ := s.CmWidth.Float64()
	fh, _ := s.CmHeight.Float64()

	return fmt.Sprintf(
		"%.1fcm x %.1fcm (%dpx x %dpx)",
		fw, fh,
		s.PxWidth, s.PxHeight,
	)
}
