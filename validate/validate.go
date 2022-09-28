package validate

import (
	"fmt"
	"image"
	"math"

	"github.com/jphastings/postcards-go/internal/types"
)

const (
	smallestPx   = 640
	largestRatio = 4
	maxRatioDiff = 0.01
)

func Dimensions(meta *types.Metadata, frontBounds, backBounds image.Rectangle, frontSize, backSize types.Size) error {
	if frontSize != meta.FrontDimensions {
		return fmt.Errorf("the front image (%s) doesn't match the physical dimensions specified in the metadata file (%s)", frontSize, meta.FrontDimensions)
	}

	fw, fh := frontBounds.Dx(), frontBounds.Dy()
	bw, bh := backBounds.Dx(), backBounds.Dy()

	if fw < smallestPx || fh < smallestPx {
		return fmt.Errorf("postcard front image is too small")
	}
	if bw < smallestPx || bh < smallestPx {
		return fmt.Errorf("postcard back image is too small")
	}

	if !similarSize(frontSize, backSize, meta.PivotAxis.Heteroriented()) {
		return fmt.Errorf("the back image (%s) doesn't match the physical dimensions of the front image (%s) when flipped about the %s", backSize, frontSize, meta.PivotAxis)
	}

	if fw > largestRatio*fh {
		return fmt.Errorf("postcard front is too wide for its height")
	}
	if fh > largestRatio*fw {
		return fmt.Errorf("postcard front is too high for its width")
	}
	if bw > largestRatio*bh {
		return fmt.Errorf("postcard back is too wide for its height")
	}
	if bh > largestRatio*bw {
		return fmt.Errorf("postcard back is too high for its width")
	}

	return nil
}

// similarSize assumes all sizes are of the same unit
func similarSize(front, back types.Size, heteroriented bool) bool {
	if heteroriented {
		return similarLength(front.Width, back.Height) && similarLength(front.Height, back.Width)
	} else {
		return similarLength(front.Width, back.Width) && similarLength(front.Height, back.Height)
	}
}

func similarLength(l1, l2 types.Length) bool {
	ratio := l1.In(types.UnitCentimetre) / l2.In(types.UnitCentimetre)
	return math.Abs(1-ratio) <= maxRatioDiff
}
