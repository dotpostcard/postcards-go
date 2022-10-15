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

type Polygon []Point

type Side struct {
	Description   LocalizedText `json:"description,omitempty"`
	Transcription LocalizedText `json:"transcription,omitempty"`
	Secrets       []Polygon     `json:"secrets,omitempty"`
}

type Context struct {
	Author      Person        `json:"author"`
	Description LocalizedText `json:"description"`
}

type Person struct {
	Name string `json:"name,omitempty"`
	Uri  string `json:"uri,omitempty" yaml:"link,omitempty"`
}

type Metadata struct {
	Location        LatLong `json:"location,omitempty"`
	Flip            Flip    `json:"flip" yaml:"flip"`
	SentOn          Date    `json:"sentOn,omitempty" yaml:"sent_on"`
	Sender          Person  `json:"sender,omitempty"`
	Recipient       Person  `json:"recipient,omitempty"`
	Front           Side    `json:"front,omitempty"`
	Back            Side    `json:"back,omitempty"`
	FrontDimensions Size    `json:"frontSize" yaml:",omitempty"`
	Context         Context `json:"context,omitempty"`
}

var _ json.Marshaler = (*LatLong)(nil)
var _ yaml.Marshaler = (*LatLong)(nil)
var _ json.Unmarshaler = (*LatLong)(nil)
var _ yaml.Unmarshaler = (*LatLong)(nil)

var _ json.Marshaler = (*Polygon)(nil)
var _ yaml.Marshaler = (*Polygon)(nil)
var _ json.Unmarshaler = (*Polygon)(nil)
var _ yaml.Unmarshaler = (*Polygon)(nil)
