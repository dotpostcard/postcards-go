package compile

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math/big"

	"github.com/h2non/bimg"
	"github.com/jphastings/postcard-go/internal/types"
)

func readerToImage(r io.Reader) (*bimg.Image, *types.Dimensions, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}
	vips := bimg.NewImage(b)
	size, err := vips.Size()
	if err != nil {
		return nil, nil, err
	}

	hRes, vRes, err := extractResolution(vips)
	if err != nil {
		return nil, nil, err
	}

	dims := &types.Dimensions{
		Width:  resolutionToCentimeters(size.Width, hRes),
		Height: resolutionToCentimeters(size.Height, vRes),
	}

	switch vips.Type() {
	case "jpeg", "png":
		// These are useable image types
	default:
		pngData, err := vips.Convert(bimg.PNG)
		if err != nil {
			return nil, nil, err
		}
		vips = bimg.NewImage(pngData)
	}

	return vips, dims, nil
}

func extractResolution(vips *bimg.Image) (*big.Rat, *big.Rat, error) {
	var (
		imMeta bimg.ImageMetadata
		err    error
	)

	if vips.Type() == "png" {
		// This is mega annoying; bimg can't extract exif data from PNG files, apparently
		exifData, err := vips.Convert(bimg.JPEG)
		if err != nil {
			return nil, nil, err
		}

		if imMeta, err = bimg.Metadata(exifData); err != nil {
			return nil, nil, err
		}
	} else if imMeta, err = vips.Metadata(); err != nil {
		return nil, nil, err
	}

	scaler, err := exifResolutionScaler(imMeta.EXIF.ResolutionUnit)
	if err != nil {
		return nil, nil, err
	}

	hRes, err := exifResolutionToFloat(imMeta.EXIF.XResolution)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid horizontal resolution in EXIF data: %v", err)
	}

	vRes, err := exifResolutionToFloat(imMeta.EXIF.YResolution)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid vertical resolution in EXIF data: %v", err)
	}

	return hRes.Mul(hRes, scaler), vRes.Mul(vRes, scaler), nil
}

func resolutionToCentimeters(pixels int, res *big.Rat) types.Centimeters {
	cms := res.Quo(big.NewRat(int64(pixels), 1), res)
	return types.Centimeters(cms)
}

// Resolutions are specified in 'rational64u' format: https://exiftool.org/TagNames/EXIF.html#:~:text=0x011a-,XResolution,-rational64u%3A
func exifResolutionToFloat(res string) (*big.Rat, error) {
	var a, b int64
	if _, err := fmt.Sscanf(res, "%d/%d", &a, &b); err != nil {
		return &big.Rat{}, fmt.Errorf("invalid width resolution in EXIF data: %v", err)
	}

	return big.NewRat(a, b), nil
}

// As defined by https://exiftool.org/TagNames/EXIF.html#:~:text=0x0128-,ResolutionUnit,-int16u%3A
func exifResolutionScaler(unit int) (*big.Rat, error) {
	switch unit {
	case 0, 1: // None
		return &big.Rat{}, fmt.Errorf("no resolution unit in EXIF data for physical dimensions of image")
	case 2: // Inches
		return big.NewRat(100, 254), nil // Who knew, an inch is *exactly* 2.54 cm, as of 1959?
	case 3: // Centimeters
		return big.NewRat(1, 1), nil
	default:
		return &big.Rat{}, fmt.Errorf("invalid unit in EXIF data for physical dimensions of image (%d)", unit)
	}
}
