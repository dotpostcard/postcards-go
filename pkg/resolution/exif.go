package resolution

import (
	"fmt"
	"math/big"

	"github.com/jphastings/postcards-go/internal/types"
	"github.com/kolesa-team/goexiv"
)

const (
	resUnitCode = "Exif.Image.ResolutionUnit"
	resXCode    = "Exif.Image.XResolution"
	resYCode    = "Exif.Image.YResolution"
)

func decodeExif(data []byte) (types.Resolution, error) {
	im, err := goexiv.OpenBytes(data)
	if err != nil {
		return types.Resolution{}, err
	}

	if err := im.ReadMetadata(); err != nil {
		return types.Resolution{}, err
	}

	fmt.Println(im.GetIptcData().AllTags())

	return GetExifResolution(im)
}

func GetExifResolution(im *goexiv.Image) (types.Resolution, error) {
	unit, err := getExifResUnit(im)
	if err != nil {
		return types.Resolution{}, err
	}

	xRes, err := getExifResCount(im, resXCode)
	if err != nil {
		return types.Resolution{}, err
	}
	yRes, err := getExifResCount(im, resYCode)
	if err != nil {
		return types.Resolution{}, err
	}

	return types.Resolution{
		XResolution: types.PixelDensity{Count: xRes, Unit: unit},
		YResolution: types.PixelDensity{Count: yRes, Unit: unit},
	}, nil
}

func getExifResUnit(im *goexiv.Image) (*types.PixelDensityUnit, error) {
	unit, err := im.GetExifData().GetString(resUnitCode)
	if err != nil {
		return nil, err
	}

	switch unit {
	case "3":
		return types.UnitPixelsPerCentimetre, nil
	case "2", "1":
		return types.UnitPixelsPerInch, nil
	default:
		return nil, fmt.Errorf("unknown EXIF resolution unit")
	}
}

func getExifResCount(im *goexiv.Image, tag string) (*big.Rat, error) {
	val, err := im.GetExifData().GetString(tag)
	if err != nil {
		return nil, err
	}

	var a, b int64
	if _, err := fmt.Sscanf(val, "%d/%d", &a, &b); err != nil {
		return nil, fmt.Errorf("invalid rational number format: %w", err)
	}

	return big.NewRat(a, b), nil
}
