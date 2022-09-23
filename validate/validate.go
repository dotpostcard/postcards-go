package validate

import (
	"fmt"
	"math"
	"strconv"

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

	frontDim, err := dimensions(frontImg)
	if err != nil {
		return err
	}
	backDim, err := dimensions(frontImg)
	if err != nil {
		return err
	}

	fmt.Println("Dimensions", frontDim, backDim)
	if pc.Meta.FrontDimensions == nil {
		pc.Meta.FrontDimensions = frontDim

		fmt.Println(frontDim)
	} else if frontDim != pc.Meta.FrontDimensions {
		return fmt.Errorf("the front image doesn't match the physical dimensions specified in the metadata file")
	}

	if !frontDim.Same(backDim, pc.Meta.PivotAxis.Heteroriented()) {
		return fmt.Errorf("the back image doesn't match the physical dimensions of the front image (when flipped about the %s, as defined)", pc.Meta.PivotAxis)
	}

	frontRatio := types.NewAspectRatio(float64(frontSize.Width), float64(frontSize.Height))
	backRatio := types.NewAspectRatio(float64(frontSize.Width), float64(frontSize.Height))

	var flippedBack types.AspectRatio
	if pc.Meta.PivotAxis.Heteroriented() {
		flippedBack = 1 / backRatio
	} else {
		flippedBack = backRatio
	}

	ratioDiff := math.Abs(1 - float64(frontRatio)/float64(flippedBack))
	if ratioDiff > maxRatioDiff {
		return fmt.Errorf(
			"image sizes don't align: aspect ratios of postcard front (%s) & back (%s) are different "+
				"by %.1f%% (when flipped along the %s, as defined)", frontRatio, backRatio, ratioDiff*100, pc.Meta.PivotAxis)
	}
	return nil
}

func dimensions(img *bimg.Image) (*types.Dimensions, error) {
	meta, err := img.Metadata()
	if err != nil {
		return nil, err
	}

	scaler, err := exifResolutionScaler(meta.EXIF.ResolutionUnit)
	if err != nil {
		return nil, err
	}

	width, err := strconv.ParseFloat(meta.EXIF.XResolution, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid width resolution in EXIF data: %v", err)
	}

	height, err := strconv.ParseFloat(meta.EXIF.YResolution, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid height resolution in EXIF data: %v", err)
	}

	return &types.Dimensions{
		Width:  types.Centimeters(width * scaler),
		Height: types.Centimeters(height * scaler),
	}, nil
}

// As defined by https://exiftool.org/TagNames/EXIF.html#:~:text=0x0128-,ResolutionUnit,-int16u%3A
func exifResolutionScaler(unit int) (float64, error) {
	switch unit {
	case 1: // None
		return 0, fmt.Errorf("no resolution unit in EXIf data for physical dimensions of image")
	case 2: // Inches
		return 0.393701, nil
	case 3: // Centimeters
		return 1, nil
	default:
		return 0, fmt.Errorf("invalid unit in EXIf data for physical dimensions of image")
	}
}
