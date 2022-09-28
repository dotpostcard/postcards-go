package postcards_test

import (
	"testing"

	"github.com/dotpostcard/postcards-go"
)

func TestReadFile(t *testing.T) {
	_, err := postcards.ReadFile("fixtures/hello.postcard", false)
	if err != nil {
		t.Error(err)
	}
}
