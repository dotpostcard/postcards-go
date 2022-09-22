package compile

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/h2non/bimg"
	"github.com/jphastings/postcard-go/internal/types"
)

const (
	smallestDim  = 1024
	largestRatio = 4
	maxRatioDiff = 0.01
)

// ReaderToImage converts any input image type to a webp image, with any secret regions blurred out
func ReaderToImage(r io.Reader, secrets []types.Polygon) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	img := bimg.NewImage(buf.Bytes())
	if len(secrets) != 0 {
		log.Println("BEWARE! Automatic blurring of secret regions of postcards is not yet implemented!")
	}

	return img.Convert(bimg.WEBP)
}

func ValidateDimensions(frontData, backData []byte) error {
	frontImg := bimg.NewImage(frontData)
	frontSize, err := frontImg.Size()
	if err != nil {
		return err
	}

	backImg := bimg.NewImage(backData)
	backSize, err := backImg.Size()
	if err != nil {
		return err
	}

	if frontSize.Width < smallestDim || frontSize.Height < smallestDim {
		return fmt.Errorf("postcard front is too small")
	}
	if backSize.Width < smallestDim || backSize.Height < smallestDim {
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

	frontRatio := float64(frontSize.Width) / float64(frontSize.Height)
	backRatio := float64(backSize.Width) / float64(backSize.Height)
	if frontRatio > 1 && backRatio < 1 || backRatio > 1 && frontRatio < 1 {
		backRatio = 1 / backRatio
	}

	ratioDiff := math.Abs(1 - frontRatio/backRatio)
	if ratioDiff > maxRatioDiff {
		return fmt.Errorf("postcard front & back are more than %.1f%% different in aspect ratio", maxRatioDiff*100)
	}
	return nil
}
