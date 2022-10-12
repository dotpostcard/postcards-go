package postcards

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/dotpostcard/postcards-go/internal/types"
)

var byteOrder = binary.BigEndian

// Write creates the postcard file tarball from the in-memory object, writing to the given Writer
func Write(pc *types.Postcard, w io.Writer) error {
	if err := writeVersion(w, Version); err != nil {
		return err
	}
	if err := writeMeta(w, pc.Meta); err != nil {
		return err
	}
	if err := writeImage(w, pc.Front, "front"); err != nil {
		return err
	}
	if err := writeImage(w, pc.Back, "back"); err != nil {
		return err
	}

	return nil
}

func writeVersion(w io.Writer, ver types.Version) error {
	if _, err := w.Write(magicBytes); err != nil {
		return err
	}

	if err := binary.Write(w, byteOrder, ver.Major); err != nil {
		return err
	}
	if err := binary.Write(w, byteOrder, ver.Minor); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, ver.Patch)
}

func writeMeta(w io.Writer, meta types.Metadata) error {
	buf, err := MetadataBytes(meta, false)
	if err != nil {
		return err
	}

	if err := binary.Write(w, byteOrder, int32(len(buf))); err != nil {
		return err
	}

	_, writeErr := w.Write(buf)
	return writeErr
}

func MetadataBytes(meta types.Metadata, pretty bool) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)

	enc.SetEscapeHTML(false)
	if pretty {
		enc.SetIndent("", "  ")
	} else {
		enc.SetIndent("", "")
	}

	if err := enc.Encode(meta); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func writeImage(w io.Writer, img []byte, name string) error {
	if err := binary.Write(w, byteOrder, int32(len(img))); err != nil {
		return err
	}

	_, err := w.Write(img)
	return err
}
