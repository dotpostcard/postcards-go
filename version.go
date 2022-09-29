package postcards

import (
	_ "embed"

	"github.com/dotpostcard/postcards-go/internal/types"
)

//go:embed VERSION
var versionData string

// Version is the semantic version of this reference implementation of postcard file format reading and writing
var Version = types.MustParseVersion(versionData)
