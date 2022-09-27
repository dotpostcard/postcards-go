package types

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type Postcard struct {
	Meta  Metadata
	Front []byte
	Back  []byte
}

type LatLong struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

type LocalizedText map[string]string
type Polygon []Point

type Side struct {
	Description   LocalizedText `json:"description"`
	Transcription string        `json:"transcription"`
	Secrets       []Polygon     `json:"secrets,omitempty"`
}

type Metadata struct {
	Location        LatLong   `json:"location"`
	PivotAxis       PivotAxis `json:"pivot_axis" yaml:"flip"`
	SentOn          Date      `json:"sent_on" yaml:"sent_on"`
	Sender          string    `json:"sender"`
	Recipient       string    `json:"recipient"`
	Front           Side      `json:"front"`
	Back            Side      `json:"back"`
	FrontDimensions Size      `json:"front_size" yaml:",omitempty"`
}

var _ json.Marshaler = (*LatLong)(nil)
var _ yaml.Marshaler = (*LatLong)(nil)
var _ json.Unmarshaler = (*LatLong)(nil)
var _ yaml.Unmarshaler = (*LatLong)(nil)

var _ json.Marshaler = (*Polygon)(nil)
var _ yaml.Marshaler = (*Polygon)(nil)
var _ json.Unmarshaler = (*Polygon)(nil)
var _ yaml.Unmarshaler = (*Polygon)(nil)
