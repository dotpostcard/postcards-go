package compile

import (
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
