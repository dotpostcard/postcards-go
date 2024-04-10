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
)

var nameRegex = regexp.MustCompile(`(.+)-(?:front|back|meta)+\.[a-z]+`)

// Files accepts a path to one of the three needed files, attempts to find the others, and provides the conventional name and bytes for the file.
func Files(part string, skipIfPresent bool, webFormat bool) ([]string, [][]byte, error) {
	dir := filepath.Dir(part)
	parts := nameRegex.FindStringSubmatch(filepath.Base(part))
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("given filename not of the form *-{front,back,meta}.ext")
	}
	prefix := parts[1]
	var outputFilenames []string
	if webFormat {
		outputFilenames = []string{
			fmt.Sprintf("%s.webp", prefix),
			fmt.Sprintf("%s.json", prefix),
		}
	} else {
		outputFilenames = []string{fmt.Sprintf("%s.postcard", prefix)}
	}

	exists, err := anyFilesExist(outputFilenames...)
	if err != nil {
		return outputFilenames, nil, nil
	}
	if skipIfPresent && exists {
		return outputFilenames, nil, fmt.Errorf("output file already exists: %s", strings.Join(outputFilenames, ", "))
	}

	metaRaw, metaExt, err := openVagueFilename(dir, prefix, "meta", "json", "yml", "yaml")
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't load metadata: %w", err)
	}
	meta, err := metaReader(metaRaw, metaExt)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't parse metadata: %w", err)
	}

	front, _, err := openVagueFilename(dir, prefix, "front", "png", "jpg", "tif", "tiff")
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't load postcard front: %w", err)
	}
	back, _, err := openVagueFilename(dir, prefix, "back", "png", "jpg", "tif", "tiff")
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't load postcard back: %w", err)
	}

	if webFormat {
		img, md, err := CompileWeb(front, back, meta)
		if err != nil {
			return nil, nil, err
		}

		return outputFilenames, [][]byte{img, md}, nil
	} else {
		pc, err := Readers(front, back, meta)
		if err != nil {
			return nil, nil, err
		}

		buf := new(bytes.Buffer)
		if err := postcards.Write(pc, buf); err != nil {
			return nil, nil, err
		}

		return outputFilenames, [][]byte{buf.Bytes()}, nil
	}
}

func anyFilesExist(filenames ...string) (bool, error) {
	for _, filename := range filenames {
		info, err := os.Stat(filename)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return false, err
		}
		if info.IsDir() {
			return true, fmt.Errorf("file %s is a directory", filename)
		}
		return true, nil
	}
	return false, nil
}

func processImages(frontReader, backReader io.Reader, mp MetadataProvider) (image.Image, image.Image, types.Size, types.Size, types.Metadata, error) {
	meta, err := mp.Metadata()
	if err != nil {
		return nil, nil, types.Size{}, types.Size{}, types.Metadata{}, fmt.Errorf("unable to obtain the metadata: %w", err)
	}

	if err := validateMetadata(meta); err != nil {
		return nil, nil, types.Size{}, types.Size{}, types.Metadata{}, fmt.Errorf("metadata invalid: %w", err)
	}

	frontRaw, frontDims, err := readerToImage(frontReader)
	if err != nil {
		return nil, nil, types.Size{}, types.Size{}, types.Metadata{}, fmt.Errorf("unable to parse image for front image: %w", err)
	}
	backRaw, backDims, err := readerToImage(backReader)
	if err != nil {
		return nil, nil, types.Size{}, types.Size{}, types.Metadata{}, fmt.Errorf("unable to parse image for back image: %w", err)
	}

	meta.FrontDimensions = bestFrontDimensions(meta.FrontDimensions, frontDims, backDims, meta.Flip.Heteroriented())

	if err := validate.Dimensions(&meta, frontRaw.Bounds(), backRaw.Bounds(), frontDims, backDims); err != nil {
		return nil, nil, types.Size{}, types.Size{}, types.Metadata{}, err
	}
	if isOversized(frontDims) {
		log.Printf("WARNING! This postcard is very large (%s), do the images have the correct ppi/ppcm?\n", frontDims)
	}

	frontImg, err := hideSecrets(frontRaw, meta.Front.Secrets)
	if err != nil {
		return nil, nil, types.Size{}, types.Size{}, types.Metadata{}, fmt.Errorf("unable to hide the secret areas specified on the postcard front: %w", err)
	}
	backImg, err := hideSecrets(backRaw, meta.Back.Secrets)
	if err != nil {
		return nil, nil, types.Size{}, types.Size{}, types.Metadata{}, fmt.Errorf("unable to hide the secret areas specified on the postcard back: %w", err)
	}

	return frontImg, backImg, frontDims, backDims, meta, nil
}

// Readers accepts reader objects for each of the components of a postcard file, and creates an in-memory Postcard object.
func Readers(frontReader, backReader io.Reader, mp MetadataProvider) (*types.Postcard, error) {
	frontImg, backImg, frontDims, backDims, meta, err := processImages(frontReader, backReader, mp)
	if err != nil {
		return nil, err
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

func openVagueFilename(dir, prefix, suffix string, extensions ...string) (io.Reader, string, error) {
	for _, ext := range extensions {
		filename := path.Join(dir, fmt.Sprintf("%s-%s.%s", prefix, suffix, ext))
		r, err := os.Open(filename)
		if err == nil {
			return r, ext, nil
		}
	}
	return nil, "", fmt.Errorf("no file '%s-%s.{%s}' in %s", prefix, suffix, strings.Join(extensions, ","), dir)
}

func metaReader(r io.Reader, ext string) (MetadataProvider, error) {
	switch ext {
	case "json":
		return MetadataFromJSON(r), nil
	case "yaml", "yml":
		return MetadataFromYaml(r), nil
	default:
		return nil, fmt.Errorf("unknown metadata format: %s", ext)
	}
}

var oversized float64 = 30 // Centimetres

func isOversized(s types.Size) bool {
	if !s.HasPhysical() {
		return false
	}

	if w, _ := s.CmWidth.Float64(); w > oversized {
		return true
	}
	if h, _ := s.CmHeight.Float64(); h > oversized {
		return true
	}

	return false
}

func bestFrontDimensions(fromMeta, fromFront, fromBack types.Size, isHeteroriented bool) types.Size {
	bestSize := types.Size{
		CmWidth:  fromMeta.CmWidth,
		CmHeight: fromMeta.CmHeight,
		PxWidth:  fromFront.PxWidth,
		PxHeight: fromFront.PxHeight,
	}

	if bestSize.HasPhysical() {
		// TODO: Flag if resolutions are wildly off
		return bestSize
	}

	bestSize.CmWidth = fromFront.CmWidth
	bestSize.CmHeight = fromFront.CmHeight

	if bestSize.HasPhysical() {
		return bestSize
	}

	if isHeteroriented {
		bestSize.CmWidth = fromBack.CmHeight
		bestSize.CmHeight = fromBack.CmWidth
	} else {
		bestSize.CmWidth = fromBack.CmWidth
		bestSize.CmHeight = fromBack.CmHeight
	}

	return bestSize
}
