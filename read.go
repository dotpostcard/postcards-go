package postcards

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dotpostcard/postcards-go/internal/types"
)

var magicBytes = []byte("postcard")

var (
	cannotRead = types.Version{Major: Version.Major + 1, Minor: Version.Minor, Patch: Version.Patch}
	warnOnRead = types.Version{Major: Version.Major, Minor: Version.Minor + 1, Patch: Version.Patch}
)

// Read will parse a Postcard struct from a Reader
func Read(r io.Reader, metaOnly bool) (*types.Postcard, error) {
	pc := &types.Postcard{}

	// TODO: Skip ahead if metaOnly is true
	frontBytes, err := readImage(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read front image: %v", err)
	}
	pc.Front = frontBytes

	if !hasMagicBytes(r) {
		return nil, fmt.Errorf("not valid postcard file")
	}

	version, err := readVersion(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read version data: %w", err)
	}

	if cannotRead.LessThan(version) {
		return nil, fmt.Errorf("postcard is too new to be processed (postcard is v%s, library cannot read v%v or above)", version, cannotRead)
	}

	if warnOnRead.LessThan(version) {
		log.Printf("This postcard (v%s) may have features this library cannot make use of. Upgrade to v%v or greater to remove this warning.", version, warnOnRead)
	}

	meta, err := readMeta(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read metadata: %v", err)
	}
	pc.Meta = meta

	if metaOnly {
		return pc, nil
	}

	backBytes, err := readImage(r)
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

func hasMagicBytes(r io.Reader) bool {
	b := make([]byte, len(magicBytes))
	if n, err := r.Read(b); err != nil || n < len(magicBytes) {
		return false
	}

	for i := 0; i < len(magicBytes); i++ {
		if b[i] != magicBytes[i] {
			return false
		}
	}

	return true
}

func readVersion(r io.Reader) (types.Version, error) {
	b := make([]byte, 3)
	if n, err := r.Read(b); err != nil {
		return types.Version{}, err
	} else if n < 3 {
		return types.Version{}, fmt.Errorf("not valid postcard version number")
	}

	return types.Version{Major: b[0], Minor: b[1], Patch: b[2]}, nil
}

func readSize(r io.Reader) (uint32, error) {
	b := make([]byte, 4)
	if n, err := r.Read(b); err != nil {
		return 0, err
	} else if n < 2 {
		return 0, io.EOF
	}

	return byteOrder.Uint32(b), nil
}

func readMeta(r io.Reader) (types.Metadata, error) {
	var meta types.Metadata

	size, err := readSize(r)
	if err != nil {
		return meta, err
	}

	d := json.NewDecoder(&io.LimitedReader{R: r, N: int64(size)})
	if err := d.Decode(&meta); err != nil {
		return meta, err
	}

	return meta, nil
}

func readImage(r io.Reader) ([]byte, error) {
	var header bytes.Buffer
	t := io.TeeReader(r, &header)

	fourByte := make([]byte, 4)
	if _, err := t.Read(fourByte); err != nil {
		return nil, err
	}
	if string(fourByte) != "RIFF" {
		return nil, fmt.Errorf("not valid postcard file, expected WebP image")
	}

	if _, err := r.Read(fourByte); err != nil {
		return nil, err
	}
	webpSize := binary.LittleEndian.Uint32(fourByte)

	if _, err := t.Read(fourByte); err != nil {
		return nil, err
	}
	if string(fourByte) != "WEBP" {
		return nil, fmt.Errorf("not valid postcard file, expected WebP image")
	}

	rest, err := io.ReadAll(&io.LimitedReader{R: r, N: int64(webpSize - 4)})
	if err != nil {
		return nil, err
	}

	return append(header.Bytes(), rest...), nil
}
