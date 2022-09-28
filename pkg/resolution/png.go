package resolution

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/dotpostcard/postcards-go/internal/types"
	pngstructure "github.com/dsoprea/go-png-image-structure"
)

const (
	pngHeader = "\x89\x50\x4E\x47\x0D\x0A\x1A\x0A"
)

func decodePng(data []byte) (types.Resolution, error) {
	pmp := pngstructure.NewPngMediaParser()

	intfc, err := pmp.ParseBytes(data)
	if err != nil {
		return types.Resolution{}, err
	}

	cs := intfc.(*pngstructure.ChunkSlice)
	index := cs.Index()
	phys, ok := index["pHYs"]
	if !ok {
		// No physical dimension information
		return types.Resolution{}, nil
	}
	b := phys[0].Data
	if len(b) < 9 {
		return types.Resolution{}, fmt.Errorf("incomplete PNG pHYs header")
	}

	unit := b[8]
	if unit != 1 {
		return types.Resolution{}, fmt.Errorf("invalid PNG resolution unit (%d)", unit)
	}

	pdX := binary.BigEndian.Uint32(b[0:4])
	pdY := binary.BigEndian.Uint32(b[4:8])

	return types.Resolution{
		XResolution: types.PixelDensity{Count: big.NewRat(int64(pdX), 1), Unit: types.UnitPixelsPerMetre},
		YResolution: types.PixelDensity{Count: big.NewRat(int64(pdY), 1), Unit: types.UnitPixelsPerMetre},
	}, nil
}
