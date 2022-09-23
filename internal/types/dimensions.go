package types

type Centimeters float64
type Dimensions struct {
	Width  Centimeters
	Height Centimeters
}

func (d *Dimensions) Same(other *Dimensions, heteroriented bool) bool {
	if heteroriented {
		return d.Width == other.Height && d.Height == other.Width
	} else {
		return d.Width == other.Width && d.Height == other.Height
	}
}
