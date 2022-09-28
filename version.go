package postcards

import (
	_ "embed"

	"github.com/Masterminds/semver"
)

//go:embed VERSION
var versionData string

// Version is the semantic version of this reference implementation of postcard file format reading and writing
var Version = semver.MustParse(versionData)
