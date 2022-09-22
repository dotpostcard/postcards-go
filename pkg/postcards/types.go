package postcards

import (
	"encoding/json"
	"image"

	"cloud.google.com/go/civil"
	"github.com/Masterminds/semver"
	"github.com/h2non/bimg"
)

var Version = semver.MustParse("0.0.0")

type Postcard struct {
	Version *semver.Version
	Meta    PostcardMetadata
	Front   *bimg.Image
	Back    *bimg.Image
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
	SentOn     Date         `json:"sent_on"`
	Senders    []string     `json:"senders"`
	Recipients []string     `json:"recipients"`
	Front      PostcardSide `json:"front"`
	Back       PostcardSide `json:"back"`
}

var _ json.Unmarshaler = (*LatLong)(nil)
var _ json.Unmarshaler = (*Date)(nil)
var _ json.Unmarshaler = (*Polygon)(nil)
