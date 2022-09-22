package postcard

import (
	"encoding/json"
	"image"

	"cloud.google.com/go/civil"
	"github.com/h2non/bimg"
	"gopkg.in/yaml.v3"
)

type Postcard struct {
	Meta  PostcardMetadata
	Front *bimg.Image
	Back  *bimg.Image
}

type LatLong struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

type LocalizedText map[string]string
type Polygon []image.Point
type Date civil.Date

type PivotAxis uint

const (
	PivotAxisUp PivotAxis = iota
	PivotAxisUpRight
	PivotAxisRight
	PivotAxisDownRight
)

type PostcardSide struct {
	Description   LocalizedText `json:"description"`
	Transcription string        `json:"transcription"`
	Secrets       []Polygon     `json:"secrets"`
}

type PostcardMetadata struct {
	Location   LatLong      `json:"location"`
	PivotAxis  PivotAxis    `json:"pivot_axis"`
	SentOn     Date         `json:"sent_on" yaml:"sent_on"`
	Senders    []string     `json:"senders"`
	Recipients []string     `json:"recipients"`
	Front      PostcardSide `json:"front"`
	Back       PostcardSide `json:"back"`
}

var _ json.Marshaler = (*LatLong)(nil)
var _ yaml.Marshaler = (*LatLong)(nil)
var _ json.Unmarshaler = (*LatLong)(nil)
var _ yaml.Unmarshaler = (*LatLong)(nil)

var _ json.Marshaler = (*Polygon)(nil)
var _ yaml.Marshaler = (*Polygon)(nil)
var _ json.Unmarshaler = (*Polygon)(nil)
var _ yaml.Unmarshaler = (*Polygon)(nil)

var _ json.Marshaler = (*Date)(nil)
var _ yaml.Marshaler = (*Date)(nil)
var _ json.Unmarshaler = (*Date)(nil)
var _ yaml.Unmarshaler = (*Date)(nil)
