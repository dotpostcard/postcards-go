package postcard

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/jphastings/postcard-go/internal/types"
)

var (
	cannotRead = semver.MustParse("1.0.0")
	warnOnRead = semver.MustParse("0.1.0")
)

// Read will parse a Postcard struct from a Reader
func Read(r io.Reader, metaOnly bool) (*types.Postcard, error) {
	pc := &types.Postcard{}
	ar := tar.NewReader(r)

	version, err := readVersion(ar)
	if err != nil {
		return nil, fmt.Errorf("unable to read version file: %w", err)
	}

	if cannotRead.LessThan(version) {
		return nil, fmt.Errorf("postcard is too new to be processed (postcard is v%s, library cannot read v%s or above)", version, cannotRead)
	}

	if warnOnRead.LessThan(version) {
		log.Printf("This postcard (v%s) may have features this library cannot make use of. Upgrade to v%s or greater to remove this warning.", version, warnOnRead)
	}

	meta, err := readMeta(ar)
	if err != nil {
		return nil, fmt.Errorf("unable to read metadata: %v", err)
	}
	pc.Meta = meta

	if metaOnly {
		return pc, nil
	}

	frontBytes, err := readImage(ar, "front")
	if err != nil {
		return nil, fmt.Errorf("unable to read front image: %v", err)
	}
	pc.Front = frontBytes

	backBytes, err := readImage(ar, "back")
	if err != nil {
		return nil, fmt.Errorf("unable to read back image: %v", err)
	}
	pc.Back = backBytes

	return pc, nil
}

// ReadFile is a convenience method for reading from a .postcard file from disk
func ReadFile(path string, metaOnly bool) (*types.Postcard, error) {
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

func readMeta(ar *tar.Reader) (types.Metadata, error) {
	var meta types.Metadata

	hdr, err := ar.Next()
	if err != nil {
		return meta, fmt.Errorf("not a valid postcard tarball, missing metadata: %w", err)
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

func readImage(ar *tar.Reader, name string) ([]byte, error) {
	hdr, err := ar.Next()
	if err != nil {
		return nil, fmt.Errorf("not a valid postcard tarball, missing %s: %w", name, err)
	}
	if hdr.Name != name+".webp" {
		return nil, fmt.Errorf("missing %s image file, got %s first instead", name, hdr.Name)
	}

	return io.ReadAll(ar)
}
