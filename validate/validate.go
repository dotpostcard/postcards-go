package validate

import (
	"fmt"

	"github.com/h2non/bimg"
	"github.com/jphastings/postcard-go/internal/types"
)

const (
	smallestPx   = 640
	largestRatio = 4
	maxRatioDiff = 0.01
)

func Dimensions(meta *types.Metadata, frontImg, backImg *bimg.Image, frontDims, backDims *types.Dimensions) error {
	if frontDims != meta.FrontDimensions {
		return fmt.Errorf("the front image (%s) doesn't match the physical dimensions specified in the metadata file (%s)", frontDims, meta.FrontDimensions)
	}

	frontSize, err := frontImg.Size()
	if err != nil {
		return err
	}
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

	if !frontDims.SimilarSize(backDims, meta.PivotAxis.Heteroriented(), maxRatioDiff) {
		return fmt.Errorf("the back image (%s) doesn't match the physical dimensions of the front image (%s) when flipped about the %s", backDims, frontDims, meta.PivotAxis)
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
