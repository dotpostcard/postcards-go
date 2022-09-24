package postcard_test

import (
	"testing"

	"github.com/jphastings/postcard-go"
)

func TestReadFile(t *testing.T) {
	_, err := postcard.ReadFile("fixtures/hello.postcard", false)
	if err != nil {
		t.Error(err)
	}
}
