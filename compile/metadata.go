package compile

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dotpostcard/postcards-go/internal/types"
	"gopkg.in/yaml.v3"
)

type MetadataProvider interface {
	Metadata() (types.Metadata, error)
}

func MetadataFromYaml(r io.Reader) yamlMetadata {
	return yamlMetadata{reader: r}
}

type yamlMetadata struct {
	reader io.Reader
}

func (ym yamlMetadata) Metadata() (types.Metadata, error) {
	var meta types.Metadata
	err := yaml.NewDecoder(ym.reader).Decode(&meta)
	return meta, err
}

func MetadataFromJSON(r io.Reader) jsonMetadata {
	return jsonMetadata{reader: r}
}

type jsonMetadata struct {
	reader io.Reader
}

func (jm jsonMetadata) Metadata() (types.Metadata, error) {
	var meta types.Metadata
	err := json.NewDecoder(jm.reader).Decode(&meta)
	return meta, err
}

func validateMetadata(meta types.Metadata) error {
	switch meta.Flip {
	case types.FlipBook, types.FlipCalendar, types.FlipLeftHand, types.FlipRightHand:
	case "":
		return fmt.Errorf("missing flip type")
	default:
		return fmt.Errorf("invalid flip type: %s", meta.Flip)
	}

	return nil
}
