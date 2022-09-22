package loader

import (
	"fmt"
	"image"
	"io"
	"strconv"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/jphastings/postcarder/pkg/postcards"
	"gopkg.in/yaml.v3"
)

type yamlSecret struct {
	HeightBottomLeftWidthHeight string `yaml:"hblwh"`
}

type yamlPostcardSide struct {
	Description   postcards.LocalizedText
	Transcription string
	Secrets       []yamlSecret
}

type yamlMetadata struct {
	Location  string
	SentOn    civil.Date `yaml:"sent_on"`
	Sender    string
	Recipient string
	Front     yamlPostcardSide
	Back      yamlPostcardSide
}

func readerToMeta(r io.Reader) (postcards.PostcardMetadata, error) {
	var meta postcards.PostcardMetadata

	d := yaml.NewDecoder(r)
	var yMeta yamlMetadata
	if err := d.Decode(&yMeta); err != nil {
		return meta, err
	}

	frontSecrets, err := parseSecrets(yMeta.Front.Secrets)
	if err != nil {
		return meta, err
	}
	backSecrets, err := parseSecrets(yMeta.Back.Secrets)
	if err != nil {
		return meta, err
	}

	var lat, lng float64
	if _, err := fmt.Sscanf(yMeta.Location, "%f,%f", &lat, &lng); err != nil {
		return meta, fmt.Errorf("given location isn't a comma separated lat,long: %w", err)
	}

	meta.Location = postcards.LatLong{Latitude: lat, Longitude: lng}
	meta.SentOn = postcards.Date(yMeta.SentOn)
	meta.Senders = []string{yMeta.Sender}
	meta.Recipients = []string{yMeta.Recipient}
	meta.Front = postcards.PostcardSide{
		Description:   yMeta.Front.Description,
		Transcription: yMeta.Front.Transcription,
		Secrets:       frontSecrets,
	}
	meta.Back = postcards.PostcardSide{
		Description:   yMeta.Back.Description,
		Transcription: yMeta.Back.Transcription,
		Secrets:       backSecrets,
	}

	return meta, nil
}

func parseSecrets(secrets []yamlSecret) ([]postcards.Polygon, error) {
	polys := make([]postcards.Polygon, len(secrets))
	for i, secret := range secrets {
		poly, err := parseSecret(secret)
		if err != nil {
			return polys, err
		}
		polys[i] = poly
	}
	return polys, nil
}

func parseSecret(secret yamlSecret) (postcards.Polygon, error) {
	ints, err := splitNInts(secret.HeightBottomLeftWidthHeight, ",", 5)
	if err != nil {
		return nil, err
	}

	top := ints[0] - ints[1]
	left := ints[2]
	bottom := top - ints[4]
	right := left + ints[3]

	return []image.Point{
		{left, top},
		{right, top},
		{right, bottom},
		{left, bottom},
	}, nil
}

func splitNInts(in string, sep string, n int) ([]int, error) {
	parts := strings.SplitN(in, ",", n)
	intParts := make([]int, n)
	for i, part := range parts {
		intPart, err := strconv.Atoi(part)
		if err != nil {
			return intParts, err
		}
		intParts[i] = intPart
	}
	return intParts, nil
}
