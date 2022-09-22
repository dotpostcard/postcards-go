package types

import (
	"fmt"

	"cloud.google.com/go/civil"
	"gopkg.in/yaml.v3"
)

func (ll LatLong) MarshalYAML() (interface{}, error) {
	return ll.String(), nil
}

func (ll *LatLong) UnmarshalYAML(y *yaml.Node) error {
	if y.ShortTag() != "!!str" {
		return fmt.Errorf("invalid lat,long type, expected a comma separated string")
	}

	return ll.fromString(y.Value)
}

type SecretType struct {
	Type string `yaml:"type"`
}

type SecretPolygon struct {
	Type   string      `yaml:"type"`
	Points [][]float64 `yaml:"points"`
}

type SecretBox struct {
	Type   string  `yaml:"type"`
	Width  float64 `yaml:"width"`
	Height float64 `yaml:"height"`
	Left   float64 `yaml:"left"`
	Top    float64 `yaml:"top"`
}

func (pts Polygon) MarshalYAML() (interface{}, error) {
	secret := SecretPolygon{
		Type:   "polygon",
		Points: pts.toFloats(),
	}

	return yaml.Marshal(secret)
}

func (pts *Polygon) UnmarshalYAML(y *yaml.Node) error {
	var typer SecretType
	if err := y.Decode(&typer); err != nil {
		return fmt.Errorf("invalid secret definition")
	}

	switch typer.Type {
	case "box":
		var box SecretBox
		if err := y.Decode(&box); err != nil {
			return fmt.Errorf("invalid box secret definition")
		}

		return box.intoPolygon(pts)
	case "polygon":
		var polygon SecretPolygon
		if err := y.Decode(&polygon); err != nil {
			return fmt.Errorf("invalid polygon secret definition")
		}

		return pts.fromFloats(polygon.Points)
	default:
		return fmt.Errorf("unknown secret type: %s", typer.Type)
	}
}

func (d Date) MarshalYAML() (interface{}, error) {
	return civil.Date(d).String(), nil
}

func (d *Date) UnmarshalYAML(y *yaml.Node) error {
	switch y.ShortTag() {
	case "!!timestamp", "!!str":
		parsed, err := civil.ParseDate(y.Value)
		if err != nil {
			return err
		}

		*d = Date(parsed)
	default:
		return fmt.Errorf("invalid date type, expected a string")
	}

	return nil
}

func (box SecretBox) intoPolygon(pts *Polygon) error {
	if outOfBounds(box.Width) {
		return fmt.Errorf("width of box secret is larger than 100%% of the postcard")
	}
	if outOfBounds(box.Height) {
		return fmt.Errorf("height of box secret is larger than 100%% of the postcard")
	}
	if outOfBounds(box.Left) {
		return fmt.Errorf("left edge of box secret is outside the postcard")
	}
	if outOfBounds(box.Top) {
		return fmt.Errorf("top edge of box secret is outside the postcard")
	}

	bottom := box.Top + box.Height
	if outOfBounds(bottom) {
		return fmt.Errorf("bottom edge of box secret is outside the postcard")
	}
	right := box.Left + box.Width
	if outOfBounds(right) {
		return fmt.Errorf("right edge of box secret is outside the postcard")
	}

	*pts = append(*pts, Point{X: box.Left, Y: box.Top})
	*pts = append(*pts, Point{X: right, Y: box.Top})
	*pts = append(*pts, Point{X: right, Y: bottom})
	*pts = append(*pts, Point{X: box.Left, Y: bottom})

	return nil
}

func outOfBounds(d float64) bool {
	return d < 0.0 || d > 1.0
}

// func parseSecret(secret yamlSecret) (postcards.Polygon, error) {
// 	ints, err := splitNInts(secret.HeightBottomLeftWidthHeight, ",", 5)
// 	if err != nil {
// 		return nil, err
// 	}

// 	top := ints[0] - ints[1]
// 	left := ints[2]
// 	bottom := top - ints[4]
// 	right := left + ints[3]

// 	return []image.Point{
// 		{left, top},
// 		{right, top},
// 		{right, bottom},
// 		{left, bottom},
// 	}, nil
// }

// func splitNInts(in string, sep string, n int) ([]int, error) {
// 	parts := strings.SplitN(in, ",", n)
// 	intParts := make([]int, n)
// 	for i, part := range parts {
// 		intPart, err := strconv.Atoi(part)
// 		if err != nil {
// 			return intParts, err
// 		}
// 		intParts[i] = intPart
// 	}
// 	return intParts, nil
// }
