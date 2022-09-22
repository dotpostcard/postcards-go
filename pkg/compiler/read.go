package compiler

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/h2non/bimg"
	"github.com/jphastings/postcarder/pkg/postcards"
)

func ReadPostcard(r io.Reader, metaOnly bool) (*postcards.Postcard, error) {
	pc := &postcards.Postcard{}
	ar := tar.NewReader(r)

	version, err := readVersion(ar)
	if err != nil {
		return nil, fmt.Errorf("unable to read version file: %w", err)
	}
	pc.Version = version

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

func ReadPostcardFile(path string, metaOnly bool) (*postcards.Postcard, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return ReadPostcard(f, metaOnly)
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

func readMeta(ar *tar.Reader) (postcards.PostcardMetadata, error) {
	var meta postcards.PostcardMetadata

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
