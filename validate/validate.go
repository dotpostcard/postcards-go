package validate

import (
	"fmt"
	"image"

	"github.com/jphastings/postcard-go/internal/types"
)

const (
	smallestPx   = 640
	largestRatio = 4
	maxRatioDiff = 0.01
)

func Dimensions(meta *types.Metadata, frontImg, backImg image.Image, frontDims, backDims *types.Dimensions) error {
	if frontDims != meta.FrontDimensions {
		return fmt.Errorf("the front image (%s) doesn't match the physical dimensions specified in the metadata file (%s)", frontDims, meta.FrontDimensions)
	}

	frontSize := frontImg.Bounds()
	backSize := backImg.Bounds()

	if frontSize.Dx() < smallestPx || frontSize.Dy() < smallestPx {
		return fmt.Errorf("postcard front is too small")
	}
	if backSize.Dx() < smallestPx || backSize.Dy() < smallestPx {
		return fmt.Errorf("postcard back is too small")
	}

	if !frontDims.SimilarSize(backDims, meta.PivotAxis.Heteroriented(), maxRatioDiff) {
		return fmt.Errorf("the back image (%s) doesn't match the physical dimensions of the front image (%s) when flipped about the %s", backDims, frontDims, meta.PivotAxis)
	}

	if frontSize.Dx() > largestRatio*frontSize.Dy() {
		return fmt.Errorf("postcard front is too wide for its height")
	}
	if frontSize.Dy() > largestRatio*frontSize.Dx() {
		return fmt.Errorf("postcard front is too high for its width")
	}
	if backSize.Dx() > largestRatio*backSize.Dy() {
		return fmt.Errorf("postcard back is too wide for its height")
	}
	if backSize.Dy() > largestRatio*backSize.Dx() {
		return fmt.Errorf("postcard back is too high for its width")
	}

	return nil
}
