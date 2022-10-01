package types

import (
	"fmt"
	"math/big"
)

type Size struct {
	Width  Length `json:"w"`
	Height Length `json:"h"`
}

type Length struct {
	Count *big.Rat
	Unit  *LengthUnit
}

func (s Size) Resolution(w, h int) Resolution {
	cmS := s.Convert(UnitCentimetre)

	return Resolution{
		XResolution: PixelDensity{Count: big.NewRat(1, 1).Quo(big.NewRat(int64(w), 1), cmS.Width.Count), Unit: UnitPixelsPerCentimetre},
		YResolution: PixelDensity{Count: big.NewRat(1, 1).Quo(big.NewRat(int64(h), 1), cmS.Height.Count), Unit: UnitPixelsPerCentimetre},
	}
}

func (l Length) Convert(u *LengthUnit) Length {
	scaler := big.NewRat(1, 1).Quo(&u.Rat, &l.Unit.Rat)

	return Length{
		Count: big.NewRat(1, 1).Mul(l.Count, scaler),
		Unit:  u,
	}
}

func (l Length) In(u *LengthUnit) float64 {
	fl, _ := l.Convert(u).Count.Float64()
	return fl
}

func (s Size) Convert(u *LengthUnit) Size {
	return Size{
		Width:  s.Width.Convert(u),
		Height: s.Height.Convert(u),
	}
}

func (l Length) String() string {
	return fmt.Sprintf("%s%s", l.Count.RatString(), l.Unit.String())
}

func (l *Length) fromString(str string) error {
	var count big.Rat
	if _, err := fmt.Sscanf(str, "%fcm", &count); err != nil {
		return err
	}

	l.Count = &count
	l.Unit = UnitCentimetre

	return nil
}

func (s Size) String() string {
	return fmt.Sprintf("%s x %s", s.Width.String(), s.Height.String())
}
