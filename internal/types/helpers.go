package types

import (
	"fmt"
	"image"
)

func (pts Polygon) toInts() [][]int {
	points := make([][]int, len(pts))
	for i, pt := range pts {
		points[i] = []int{pt.X, pt.Y}
	}
	return points
}

func (pts *Polygon) fromInts(points [][]int) error {
	for _, pt := range points {
		if len(pt) != 2 {
			return fmt.Errorf("%dD point given instead of 2D", len(pt))
		}

		*pts = append(*pts, image.Point{X: pt[0], Y: pt[1]})
	}

	return nil
}

func (ll LatLong) String() string {
	return fmt.Sprintf("%f,%f", ll.Latitude, ll.Longitude)
}

func (ll *LatLong) fromString(str string) error {
	var lat, lng float64
	_, err := fmt.Sscanf(str, "%f,%f", &lat, &lng)
	if err != nil {
		return err
	}

	return ll.fromFloats(lat, lng)
}

func (ll *LatLong) fromFloats(floats ...float64) error {
	if len(floats) != 2 {
		return fmt.Errorf("expected a latitude and a longitude, but %d items given", len(floats))
	}

	if floats[0] > 90 || floats[0] < -90 {
		return fmt.Errorf("given latitude (%f) is outside -90º to 90º", floats[0])
	}

	if floats[1] > 180 || floats[1] < -180 {
		return fmt.Errorf("given longitude (%f) is outside -180º to 180º", floats[1])
	}

	ll.Latitude = floats[0]
	ll.Longitude = floats[1]

	return nil
}
