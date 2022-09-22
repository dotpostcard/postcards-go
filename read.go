package postcard

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/h2non/bimg"
)

var maxReadableVersion = Version.IncMinor()

// Read will parse a Postcard struct from a Reader
func Read(r io.Reader, metaOnly bool) (*Postcard, error) {
	pc := &Postcard{}
	ar := tar.NewReader(r)

	version, err := readVersion(ar)
	if err != nil {
		return nil, fmt.Errorf("unable to read version file: %w", err)
	}

	if maxReadableVersion.LessThan(version) {
		return nil, fmt.Errorf("postcard is too new to be processed (postcard is %s, library is %s)", version, Version)
	}

	meta, err := readMeta(ar)
	if err != nil {
		return nil, err
	}
	pc.Meta = meta

	if metaOnly {
		return pc, nil
	}

	var frontBytes []byte
	if _, err := ar.Read(frontBytes); err != nil {
		return nil, err
	}
	pc.Front = bimg.NewImage(frontBytes)

	var backBytes []byte
	if _, err := ar.Read(backBytes); err != nil {
		return nil, err
	}
	pc.Back = bimg.NewImage(backBytes)

	return pc, nil
}

// ReadFile is a convenience method for reading from a .postcard file from disk
func ReadFile(path string, metaOnly bool) (*Postcard, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return Read(f, metaOnly)
}

func readVersion(ar *tar.Reader) (*semver.Version, error) {
	hdr, err := ar.Next()
	if err != nil {
		return nil, fmt.Errorf("not a valid postcard archive: %w", err)
	}
	if hdr.Name != "VERSION" {
		return nil, fmt.Errorf("missing VERSION file, got %s first instead", hdr.Name)
	}

	buf := new(strings.Builder)
	if _, err := io.Copy(buf, ar); err != nil {
		return nil, fmt.Errorf("unable to read version data: %w", err)
	}

	return semver.NewVersion(buf.String())
}

func readMeta(ar *tar.Reader) (PostcardMetadata, error) {
	var meta PostcardMetadata

	hdr, err := ar.Next()
	if err != nil {
		return meta, fmt.Errorf("not a valid tarball: %w", err)
	}
	if hdr.Name != "meta.json" {
		return meta, fmt.Errorf("missing metadata json file, got %s first instead", hdr.Name)
	}

	d := json.NewDecoder(ar)
	if err := d.Decode(&meta); err != nil {
		return meta, err
	}

	return meta, nil
}
