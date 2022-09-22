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

func (pts Polygon) MarshalYAML() (interface{}, error) {
	secret := struct {
		Type   string
		Points [][]int
	}{
		Type:   "polygon",
		Points: pts.toInts(),
	}

	return yaml.Marshal(secret)
}

func (pts *Polygon) UnmarshalYAML(y *yaml.Node) error {
	// TODO: Process secrets sections

	if y.ShortTag() != "!!map" {
		return fmt.Errorf("invalid secret section definition type")
	}

	return nil
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

// func parseSecrets(secrets []yamlSecret) ([]postcards.Polygon, error) {
// 	polys := make([]postcards.Polygon, len(secrets))
// 	for i, secret := range secrets {
// 		poly, err := parseSecret(secret)
// 		if err != nil {
// 			return polys, err
// 		}
// 		polys[i] = poly
// 	}
// 	return polys, nil
// }

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
