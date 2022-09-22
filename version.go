package postcarder

import (
	_ "embed"

	"github.com/Masterminds/semver"
)

//go:embed VERSION
var versionData string

var Version = semver.MustParse(versionData)
