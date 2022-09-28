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
	Description   LocalizedText `json:"description,omitempty"`
	Transcription string        `json:"transcription,omitempty"`
	Secrets       []Polygon     `json:"secrets,omitempty"`
}

type Metadata struct {
	Location        LatLong   `json:"location,omitempty"`
	PivotAxis       PivotAxis `json:"pivot_axis" yaml:"flip"`
	SentOn          Date      `json:"sent_on,omitempty" yaml:"sent_on"`
	Sender          string    `json:"sender,omitempty"`
	Recipient       string    `json:"recipient,omitempty"`
	Front           Side      `json:"front,omitempty"`
	Back            Side      `json:"back,omitempty"`
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
