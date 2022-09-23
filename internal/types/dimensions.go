package types

import (
	"fmt"
	"math"
	"math/big"
)

var bigPostcardCm float64 = 30

type Centimeters *big.Rat
type Dimensions struct {
	Width  Centimeters `json:"w"`
	Height Centimeters `json:"h"`
}

func (d *Dimensions) AsFloats() (float64, float64) {
	var wr, hr *big.Rat
	wr = d.Width
	hr = d.Height
	w, _ := wr.Float64()
	h, _ := hr.Float64()
	return w, h
}

func (d *Dimensions) AspectRatio() float64 {
	w, h := d.AsFloats()
	return w / h
}

func (d *Dimensions) SimilarSize(other *Dimensions, heteroriented bool, acceptableDiff float64) bool {
	if d == nil || other == nil {
		return false
	}

	var ratio float64
	if heteroriented {
		ratio = d.AspectRatio() * other.AspectRatio()
	} else {
		ratio = d.AspectRatio() / other.AspectRatio()
	}

	return math.Abs(1-ratio) <= acceptableDiff
}

func (d *Dimensions) String() string {
	if d == nil {
		return "unknown dimensions"
	}

	w, h := d.AsFloats()
	return fmt.Sprintf("%.1fcm x %.1fcm", w, h)
}

func (d *Dimensions) IsBig() bool {
	if d == nil {
		return false
	}

	w, h := d.AsFloats()

	return w >= bigPostcardCm || h >= bigPostcardCm
}
