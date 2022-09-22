package postcards

import (
	"encoding/json"
	"fmt"
	"image"

	"cloud.google.com/go/civil"
)

func (pts Polygon) MarshalJSON() ([]byte, error) {
	points := make([][]int, len(pts))
	for i, pt := range pts {
		points[i] = []int{pt.X, pt.Y}
	}
	return json.Marshal(points)
}

func (pts *Polygon) UnmarshalJSON(b []byte) error {
	var points [][]int
	if err := json.Unmarshal(b, &points); err != nil {
		return err
	}

	for _, pt := range points {
		if len(pt) != 2 {
			return fmt.Errorf("%dD point given instead of 2D", len(pt))
		}

		*pts = append(*pts, image.Point{X: pt[0], Y: pt[1]})
	}

	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(civil.Date(d).String())
}

func (d *Date) UnmarshalJSON(b []byte) error {
	var dateStr string
	if err := json.Unmarshal(b, &dateStr); err != nil {
		return err
	}

	parsed, err := civil.ParseDate(dateStr)
	if err != nil {
		return err
	}

	*d = Date(parsed)

	return nil
}

func (ll LatLong) MarshalJSON() ([]byte, error) {
	return json.Marshal([]float64{ll.Latitude, ll.Longitude})
}

func (ll *LatLong) UnmarshalJSON(b []byte) error {
	var floats []float64
	if err := json.Unmarshal(b, &floats); err != nil {
		return err
	}
	if len(floats) != 2 {
		return fmt.Errorf("expected a latitude and a longitude float, but %d floats given", len(floats))
	}

	ll.Latitude = floats[0]
	ll.Longitude = floats[1]

	return nil
}
