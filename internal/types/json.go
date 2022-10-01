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

func (l Length) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

func (l *Length) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	return l.fromString(str)
}
