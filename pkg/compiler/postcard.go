package compiler

import (
	"io"
)

type Postcard struct {
	Front io.Reader
	Back io.Reader
}