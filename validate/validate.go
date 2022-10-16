package validate

import (
	"fmt"
	"image"

	"github.com/dotpostcard/postcards-go/internal/types"
)

const (
	smallestPx   = 640
	largestRatio = 4
)

func Dimensions(meta *types.Metadata, frontBounds, backBounds image.Rectangle, frontSize, backSize types.Size) error {
	if frontSize.HasPhysical() && frontSize != meta.FrontDimensions {
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

	if !frontSize.SimilarPhysical(backSize, meta.Flip) {
		return fmt.Errorf("the back image (%s) doesn't match the physical dimensions of the front image (%s) when flipped about the %s", backSize, frontSize, meta.Flip)
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
