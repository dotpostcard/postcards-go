package validate

import (
	"fmt"
	"math/big"

	"github.com/h2non/bimg"
	"github.com/jphastings/postcard-go/internal/types"
)

const (
	smallestPx   = 640
	largestRatio = 4
	maxRatioDiff = 0.01
)

func Dimensions(pc *types.Postcard) error {
	frontImg := bimg.NewImage(pc.Front)
	frontSize, err := frontImg.Size()
	if err != nil {
		return err
	}

	backImg := bimg.NewImage(pc.Back)
	backSize, err := backImg.Size()
	if err != nil {
		return err
	}

	if frontSize.Width < smallestPx || frontSize.Height < smallestPx {
		return fmt.Errorf("postcard front is too small")
	}
	if backSize.Width < smallestPx || backSize.Height < smallestPx {
		return fmt.Errorf("postcard back is too small")
	}

	frontDim, err := dimensions(frontImg, frontSize)
	if err != nil {
		return err
	}
	backDim, err := dimensions(backImg, backSize)
	if err != nil {
		return err
	}

	if pc.Meta.FrontDimensions == nil {
		pc.Meta.FrontDimensions = frontDim
	} else if frontDim != pc.Meta.FrontDimensions {
		return fmt.Errorf("the front image (%s) doesn't match the physical dimensions specified in the metadata file (%s)", frontDim, pc.Meta.FrontDimensions)
	}

	if !frontDim.SimilarSize(backDim, pc.Meta.PivotAxis.Heteroriented(), maxRatioDiff) {
		return fmt.Errorf("the back image (%s) doesn't match the physical dimensions of the front image (%s) when flipped about the %s", backDim, frontDim, pc.Meta.PivotAxis)
	}

	if frontSize.Width > largestRatio*frontSize.Height {
		return fmt.Errorf("postcard front is too wide for its height")
	}
	if frontSize.Height > largestRatio*frontSize.Width {
		return fmt.Errorf("postcard front is too high for its width")
	}
	if backSize.Width > largestRatio*backSize.Height {
		return fmt.Errorf("postcard back is too wide for its height")
	}
	if backSize.Height > largestRatio*backSize.Width {
		return fmt.Errorf("postcard back is too high for its width")
	}

	return nil
}

func dimensions(img *bimg.Image, size bimg.ImageSize) (*types.Dimensions, error) {
	meta, err := img.Metadata()
	if err != nil {
		return nil, err
	}

	scaler, err := exifResolutionScaler(meta.EXIF.ResolutionUnit)
	if err != nil {
		return nil, err
	}

	horizontalRes, err := exifResolutionToFloat(meta.EXIF.XResolution)
	if err != nil {
		return nil, fmt.Errorf("invalid horizontal resolution in EXIF data: %v", err)
	}

	verticalRes, err := exifResolutionToFloat(meta.EXIF.YResolution)
	if err != nil {
		return nil, fmt.Errorf("invalid vertical resolution in EXIF data: %v", err)
	}

	return &types.Dimensions{
		Width:  resolutionToCentimeters(size.Width, horizontalRes, scaler),
		Height: resolutionToCentimeters(size.Height, verticalRes, scaler),
	}, nil
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
	case 1: // None
		return &big.Rat{}, fmt.Errorf("no resolution unit in EXIf data for physical dimensions of image")
	case 2: // Inches
		return big.NewRat(100, 254), nil // Who knew, an inch is *exactly* 2.54 cm, as of 1959?
	case 3: // Centimeters
		return big.NewRat(1, 1), nil
	default:
		return &big.Rat{}, fmt.Errorf("invalid unit in EXIf data for physical dimensions of image")
	}
}
