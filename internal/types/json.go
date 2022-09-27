package types

import (
	"encoding/json"
)

func (pts Polygon) MarshalJSON() ([]byte, error) {
	return json.Marshal(pts.toFloats())
}

func (pts *Polygon) UnmarshalJSON(b []byte) error {
	var points [][]float64
	if err := json.Unmarshal(b, &points); err != nil {
		return err
	}

	return pts.fromFloats(points)
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

func (s Size) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Size) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	return s.fromString(str)
}
