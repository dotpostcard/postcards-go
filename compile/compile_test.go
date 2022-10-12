package compile_test

import (
	"crypto/md5"
	"fmt"
	"os"
	"testing"

	"github.com/dotpostcard/postcards-go/compile"
	"github.com/dotpostcard/postcards-go/internal/types"
)

func ExampleFiles() {
	filename, data, err := compile.Files("../fixtures/hello-meta.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s has checksum %x", filename, md5.Sum(data))
	// Output: hello.postcard has checksum 0fa011fe11258b08ad70471d9b2ec054
}

func checkBadSetup(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("bad setup: %v", err)
	}
}

func TestReaders(t *testing.T) {
	meta, err := os.Open("../fixtures/hello-meta.yaml")
	checkBadSetup(t, err)
	front, err := os.Open("../fixtures/hello-front.png")
	checkBadSetup(t, err)
	back, err := os.Open("../fixtures/hello-back.png")
	checkBadSetup(t, err)

	pc, err := compile.Readers(front, back, compile.MetadataFromYaml(meta))
	if err != nil {
		t.Error(err)
	}

	// Metadata checks

	checks := []struct {
		name     string
		actual   interface{}
		expected interface{}
	}{
		{"Latitude", pc.Meta.Location.Latitude, 40.41365195362523},
		{"Longitude", pc.Meta.Location.Longitude, -3.6818597177370997},
		{"Sender", pc.Meta.Sender, "https://dotpostcards.org"},
		{"Recipient", pc.Meta.Recipient, "https://github.com/dotpostcard/postcards-go"},
		{"Pivot axis", pc.Meta.Flip, types.FlipBook},
		{"Sent date", pc.Meta.SentOn, types.Date("2022-09-21")},
		{"Front description", pc.Meta.Front.Description["en-GB"], "A polaroid-style framed photo of the Palacio de Cristal in Madrid's Retiro Park in Autumn."},
		{"Back description", pc.Meta.Back.Description["en-GB"], `A plain postcard back. Text at the top left declares this postcard 0033 of Madrid, "Parque del Retiro". Text at the bottom explains artwork is by Hans Löhr.`},
		{"Back transcription original", pc.Meta.Back.Transcription["original"], "en-GB"},
		{"Back transcription", pc.Meta.Back.Transcription["en-GB"], "Hello world!\n\nI hope you like this postcard from Madrid!\n\nx JP\n"},
	}
	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s should have been %v but was %s", check.name, check.expected, check.actual)
		}
	}
	_ = checks
}
