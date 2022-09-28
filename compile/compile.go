package compile

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dotpostcard/postcards-go"
	"github.com/dotpostcard/postcards-go/internal/types"
	"github.com/dotpostcard/postcards-go/validate"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/kolesa-team/goexiv"
)

var nameRegex = regexp.MustCompile(`(.+)-(?:front|back|meta)+\.[a-z]+`)

// Files accepts a path to one of the three needed files, attempts to find the others, and provides the conventional name and bytes for the file.
func Files(part string) (string, []byte, error) {
	dir := filepath.Dir(part)
	parts := nameRegex.FindStringSubmatch(filepath.Base(part))
	if len(parts) != 2 {
		return "", nil, fmt.Errorf("given filename not of the form *-{front,back,meta}.ext")
	}
	prefix := parts[1]

	meta, err := openVagueFilename(dir, prefix, "meta", "yml", "yaml")
	if err != nil {
		return "", nil, fmt.Errorf("couldn't load metadata: %w", err)
	}
	front, err := openVagueFilename(dir, prefix, "front", "png", "jpg", "tif", "tiff")
	if err != nil {
		return "", nil, fmt.Errorf("couldn't load postcard front: %w", err)
	}
	back, err := openVagueFilename(dir, prefix, "back", "png", "jpg", "tif", "tiff")
	if err != nil {
		return "", nil, fmt.Errorf("couldn't load postcard back: %w", err)
	}

	pc, err := Readers(front, back, MetadataFromYaml(meta))
	if err != nil {
		return "", nil, err
	}

	buf := new(bytes.Buffer)
	if err := postcards.Write(pc, buf); err != nil {
		return "", nil, err
	}

	return fmt.Sprintf("%s.postcard", prefix), buf.Bytes(), nil
}

// Readers accepts reader objects for each of the components of a postcard file, and creates an in-memory Postcard object.
func Readers(frontReader, backReader io.Reader, mp MetadataProvider) (*types.Postcard, error) {
	meta, err := mp.Metadata()
	if err != nil {
		return nil, fmt.Errorf("unable to obtain the metadata: %w", err)
	}

	frontRaw, frontDims, err := readerToImage(frontReader)
	if err != nil {
		return nil, fmt.Errorf("unable to parse image for front image: %w", err)
	}
	backRaw, backDims, err := readerToImage(backReader)
	if err != nil {
		return nil, fmt.Errorf("unable to parse image for back image: %w", err)
	}

	meta.FrontDimensions = frontDims

	if err := validate.Dimensions(&meta, frontRaw.Bounds(), backRaw.Bounds(), frontDims, backDims); err != nil {
		return nil, err
	}

	if isOversized(frontDims) {
		log.Printf("WARNING! This postcard is very large (%s), do the images have the correct ppi/ppcm?\n", frontDims)
	}

	frontImg, err := hideSecrets(frontRaw, meta.Front.Secrets)
	if err != nil {
		return nil, fmt.Errorf("unable to hide the secret areas specified on the postcard front: %w", err)
	}
	backImg, err := hideSecrets(backRaw, meta.Back.Secrets)
	if err != nil {
		return nil, fmt.Errorf("unable to hide the secret areas specified on the postcard back: %w", err)
	}

	frontWebp, err := encodeWebp(frontImg, frontDims)
	if err != nil {
		return nil, fmt.Errorf("unable to convert front image to WebP: %w", err)
	}
	backWebp, err := encodeWebp(backImg, backDims)
	if err != nil {
		return nil, fmt.Errorf("unable to convert back image to WebP: %w", err)
	}

	pc := &types.Postcard{
		Front: frontWebp,
		Back:  backWebp,
		Meta:  meta,
	}

	return pc, nil
}

func openVagueFilename(dir, prefix, suffix string, extensions ...string) (io.Reader, error) {
	for _, ext := range extensions {
		r, err := os.Open(path.Join(dir, fmt.Sprintf("%s-%s.%s", prefix, suffix, ext)))
		if err == nil {
			return r, nil
		}
	}
	return nil, fmt.Errorf("no file '%s-%s.{%s}' in %s", prefix, suffix, strings.Join(extensions, ","), dir)
}

var webpEncoderOpts *encoder.Options

func init() {
	opts, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 85)
	if err != nil {
		panic(err)
	}
	webpEncoderOpts = opts
}

// encodeWebp turns a the image.Image into bytes in Webp format. Currently does *not* write the resolution
// bytes into exif tags, as I can't find a good library for completing this (goexiv doesn't support writing
// rational numbers, which XResolution and YResolution are.)
func encodeWebp(img image.Image, size types.Size) ([]byte, error) {
	data := new(bytes.Buffer)
	if err := webp.Encode(data, img, webpEncoderOpts); err != nil {
		return nil, err
	}

	goIm, err := goexiv.OpenBytes(data.Bytes())
	if err != nil {
		return nil, err
	}

	return goIm.GetBytes(), nil
}

var oversized float64 = 30 // Centimetres

func isOversized(s types.Size) bool {
	return s.Width.In(types.UnitCentimetre) >= oversized || s.Height.In(types.UnitCentimetre) >= oversized
}
