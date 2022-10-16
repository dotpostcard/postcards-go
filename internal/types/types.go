package types

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Postcard struct {
	Meta  Metadata
	Front []byte
	Back  []byte
}

type Location struct {
	Name      string   `json:"name"`
	Latitude  *float64 `json:"lat,omitempty" yaml:"latitude,omitempty"`
	Longitude *float64 `json:"long,omitempty" yaml:"longitude,omitempty"`
}

func (l Location) String() string {
	if l.Latitude == nil || l.Longitude == nil {
		return l.Name
	}

	return fmt.Sprintf("%s (%.5f,%.5f)", l.Name, *l.Latitude, *l.Longitude)
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
	Name string `json:"name"`
	Uri  string `json:"uri,omitempty" yaml:"link,omitempty"`
}

func (p Person) String() string {
	if p.Uri == "" {
		return p.Name
	}

	return fmt.Sprintf("%s (%s)", p.Name, p.Uri)
}

type Metadata struct {
	Location        Location `json:"location,omitempty"`
	Flip            Flip     `json:"flip" yaml:"flip"`
	SentOn          Date     `json:"sentOn,omitempty" yaml:"sent_on"`
	Sender          Person   `json:"sender,omitempty"`
	Recipient       Person   `json:"recipient,omitempty"`
	Front           Side     `json:"front,omitempty"`
	Back            Side     `json:"back,omitempty"`
	FrontDimensions Size     `json:"frontSize" yaml:"front_size,omitempty"`
	Context         Context  `json:"context,omitempty"`
}

var _ json.Marshaler = (*Polygon)(nil)
var _ yaml.Marshaler = (*Polygon)(nil)
var _ json.Unmarshaler = (*Polygon)(nil)
var _ yaml.Unmarshaler = (*Polygon)(nil)
