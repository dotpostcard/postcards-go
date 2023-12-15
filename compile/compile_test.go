package compile_test

import (
	"crypto/md5"
	"fmt"
	"os"
	"testing"

	"github.com/dotpostcard/postcards-go"
	"github.com/dotpostcard/postcards-go/compile"
	"github.com/dotpostcard/postcards-go/internal/types"
)

func hashOfPostcardInnards(data []byte) [16]byte {
	if data[8] == postcards.Version.Major {
		data[8] = 0
	}
	if data[9] == postcards.Version.Minor {
		data[9] = 0
	}
	if data[10] == postcards.Version.Patch {
		data[10] = 0
	}

	return md5.Sum(data)
}

func ExampleFiles_postcard() {
	filenames, datas, err := compile.Files("../fixtures/hello-meta.yaml", false, false)
	if err != nil {
		panic(err)
	}

	fmt.Printf("(%d file) %s has checksum %x", len(filenames), filenames[0], hashOfPostcardInnards(datas[0]))
	// Output: (1 file) hello.postcard has checksum 659a65db02a600f70d3bf0b438d10297
}

func ExampleFiles_web() {
	filenames, datas, err := compile.Files("../fixtures/hello-meta.yaml", false, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("(%d files) %s has checksum %x, %s has checksum %x",
		len(filenames),
		filenames[0], hashOfPostcardInnards(datas[0]),
		filenames[1], hashOfPostcardInnards(datas[1]),
	)
	// Output: (2 files) hello.webp has checksum c085551191a9bf20e65392ed239d990e, hello.md has checksum 47ab2ccda9585222cc0625f5298e12c7
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
		{"Locale", pc.Meta.Locale, "en-GB"},
		{"Location name", pc.Meta.Location.Name, "Palacio de Cristal, Retiro Park, Madrid, Spain"},
		{"Location latitude", *pc.Meta.Location.Latitude, float64(40.41365195362523)},
		{"Location longitude", *pc.Meta.Location.Longitude, float64(-3.6818597177370997)},
		{"Sender (name)", pc.Meta.Sender.Name, "JP"},
		{"Sender (uri)", pc.Meta.Sender.Uri, ""},
		{"Recipient (name)", pc.Meta.Recipient.Name, "Users of @dotpostcard code"},
		{"Recipient (uri)", pc.Meta.Recipient.Uri, "https://github.com/dotpostcard/postcards-go"},
		{"Pivot axis", pc.Meta.Flip, types.FlipBook},
		{"Sent date", pc.Meta.SentOn, types.Date("2022-09-21")},
		{"Front description", pc.Meta.Front.Description, "A polaroid-style framed photo of the Palacio de Cristal in Madrid's Retiro Park in Autumn."},
		{"Back description", pc.Meta.Back.Description, `A plain postcard back. Text at the top left declares this postcard 0033 of Madrid, "Parque del Retiro". Text at the bottom explains artwork is by Hans Löhr.`},
		{"Back transcription", pc.Meta.Back.Transcription, "Hello world!\n\nI hope you like this postcard from <span lang=\"es-ES\">Madrid</span>!\n\nx JP\n"},
		{"First commentary author (name)", pc.Meta.Context.Author.Name, "JP"},
		{"First commentary author (url)", pc.Meta.Context.Author.Uri, "https://www.byJP.me"},
		{"First commentary description", pc.Meta.Context.Description, "This is a postcard I wrote, but never sent, as a fixture to use for the software repository at https://github.com/dotpostcard/postcards-go."},
	}
	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s should have been %v but was %s", check.name, check.expected, check.actual)
		}
	}
	_ = checks
}
