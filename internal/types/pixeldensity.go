package types

import (
	"fmt"
	"math/big"
)

type Resolution struct {
	XResolution PixelDensity
	YResolution PixelDensity
}

type PixelDensity struct {
	Count *big.Rat
	Unit  *PixelDensityUnit
}

func (r Resolution) Size(w, h int) Size {
	cmR := r.Convert(UnitPixelsPerCentimetre)

	return Size{
		Width:  Length{Count: big.NewRat(1, 1).Quo(big.NewRat(int64(w), 1), cmR.XResolution.Count), Unit: UnitCentimetre},
		Height: Length{Count: big.NewRat(1, 1).Quo(big.NewRat(int64(h), 1), cmR.YResolution.Count), Unit: UnitCentimetre},
	}
}

// Convert returns a new pixel density struct in the given units, converting without loss of fidelity.
// Will panic if the unit is of the wrong type.
func (pd PixelDensity) Convert(u *PixelDensityUnit) PixelDensity {
	scaler := big.NewRat(1, 1).Quo(&u.Rat, &pd.Unit.Rat)

	return PixelDensity{
		Count: big.NewRat(1, 1).Mul(pd.Count, scaler),
		Unit:  u,
	}
}

func (r Resolution) Convert(u *PixelDensityUnit) Resolution {
	return Resolution{
		XResolution: r.XResolution.Convert(u),
		YResolution: r.YResolution.Convert(u),
	}
}

func (pd PixelDensity) String() string {
	return fmt.Sprintf("%s%s", pd.Count.RatString(), pd.Unit.String())
}

func (r Resolution) String() string {
	return fmt.Sprintf("%s x %s", r.XResolution.String(), r.YResolution.String())
}
