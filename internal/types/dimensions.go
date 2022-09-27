package types

import "math/big"

// Unit holds conversion factor to SI units
type Unit struct{ big.Rat }
type LengthUnit Unit
type PixelDensityUnit Unit

var (
	metersIn1Centimeter = *big.NewRat(1, 100)
	metersIn1Inch       = *big.NewRat(254, 10000)

	UnitMetre      = &LengthUnit{Rat: *big.NewRat(1, 1)}
	UnitCentimetre = &LengthUnit{Rat: metersIn1Centimeter}
	UnitInch       = &LengthUnit{Rat: metersIn1Inch}

	UnitPixelsPerMetre      = &PixelDensityUnit{Rat: *big.NewRat(1, 1)}
	UnitPixelsPerCentimetre = &PixelDensityUnit{Rat: metersIn1Centimeter}
	UnitPixelsPerInch       = &PixelDensityUnit{Rat: metersIn1Inch}
)

func (l *LengthUnit) String() string {
	switch l {
	case UnitMetre:
		return "m"
	case UnitCentimetre:
		return "cm"
	case UnitInch:
		return "in"
	default:
		return "(unknown unit)"
	}
}

func (pd *PixelDensityUnit) String() string {
	switch pd {
	case UnitPixelsPerMetre:
		return "ppm"
	case UnitPixelsPerCentimetre:
		return "ppcm"
	case UnitPixelsPerInch:
		return "ppi"
	default:
		return "(unknown unit)"
	}
}
