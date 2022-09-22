package postcard

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"io"

	"github.com/Masterminds/semver"
	"github.com/jphastings/postcard-go/internal/types"
)

// Write creates the postcard file tarball from the in-memory object, writing to the given Writer
func Write(pc *types.Postcard, w io.Writer) error {
	ar := tar.NewWriter(w)
	defer ar.Close()

	if err := writeVersion(ar, Version); err != nil {
		return err
	}
	if err := writeMeta(ar, pc.Meta); err != nil {
		return err
	}
	if err := writeImage(ar, pc.Front, "front"); err != nil {
		return err
	}
	if err := writeImage(ar, pc.Back, "back"); err != nil {
		return err
	}

	return nil
}

func writeVersion(ar *tar.Writer, ver *semver.Version) error {
	v := []byte(ver.String())

	hdr := &tar.Header{
		Name: "VERSION",
		Mode: 0444,
		Size: int64(len(v)),
	}
	if err := ar.WriteHeader(hdr); err != nil {
		return err
	}

	_, err := ar.Write(v)
	return err
}

func writeMeta(ar *tar.Writer, meta types.Metadata) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "")
	if err := enc.Encode(meta); err != nil {
		return err
	}

	hdr := &tar.Header{
		Name: "meta.json",
		Mode: 0444,
		Size: int64(buf.Len()),
	}
	if err := ar.WriteHeader(hdr); err != nil {
		return err
	}

	_, err := ar.Write(buf.Bytes())
	return err
}

func writeImage(ar *tar.Writer, img []byte, name string) error {
	hdr := &tar.Header{
		Name: name + ".webp",
		Mode: 0444,
		Size: int64(len(img)),
	}
	if err := ar.WriteHeader(hdr); err != nil {
		return err
	}

	_, wErr := ar.Write(img)
	return wErr
}
