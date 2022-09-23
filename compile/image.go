package compile

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math/big"

	"github.com/h2non/bimg"
	"github.com/jphastings/postcard-go/internal/types"
	"golang.org/x/image/draw"
)

func readerToImage(r io.Reader) (*image.NRGBA, *types.Dimensions, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, nil, err
	}

	// TODO: is there a way to do this without converting a billion times??
	tmp, err := bimg.NewImage(buf.Bytes()).Convert(bimg.WEBP)
	if err != nil {
		return nil, nil, err
	}
	vips := bimg.NewImage(tmp)
	vMeta, err := vips.Metadata()
	if err != nil {
		return nil, nil, err
	}

	scaler, err := exifResolutionScaler(vMeta.EXIF.ResolutionUnit)
	if err != nil {
		return nil, nil, err
	}

	horizontalRes, err := exifResolutionToFloat(vMeta.EXIF.XResolution)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid horizontal resolution in EXIF data: %v", err)
	}

	verticalRes, err := exifResolutionToFloat(vMeta.EXIF.YResolution)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid vertical resolution in EXIF data: %v", err)
	}

	gImg, _, err := image.Decode(buf)
	if err != nil {
		return nil, nil, err
	}

	img, ok := gImg.(*image.NRGBA)
	if !ok {
		img = image.NewNRGBA(gImg.Bounds())
		draw.Copy(img, image.Point{}, gImg, gImg.Bounds(), draw.Over, nil)
	}

	dims := &types.Dimensions{
		Width:  resolutionToCentimeters(img.Bounds().Dx(), horizontalRes, scaler),
		Height: resolutionToCentimeters(img.Bounds().Dy(), verticalRes, scaler),
	}

	return img, dims, nil
}

func resolutionToCentimeters(pixels int, res, scaler *big.Rat) types.Centimeters {
	scaledRes := res.Mul(res, scaler)
	cms := scaledRes.Quo(big.NewRat(int64(pixels), 1), scaledRes)
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
