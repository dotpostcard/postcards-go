package types

import (
	"encoding/json"

	"cloud.google.com/go/civil"
)

func (pts Polygon) MarshalJSON() ([]byte, error) {
	return json.Marshal(pts.toInts())
}

func (pts *Polygon) UnmarshalJSON(b []byte) error {
	var points [][]int
	if err := json.Unmarshal(b, &points); err != nil {
		return err
	}

	return pts.fromInts(points)
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

	return ll.fromFloats(floats...)
}
