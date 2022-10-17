package make3d

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"html/template"
	"image"
	"image/jpeg"
	"io"

	_ "golang.org/x/image/webp"

	"github.com/dotpostcard/postcards-go"
	"github.com/dotpostcard/postcards-go/internal/types"
)

var (
	//go:embed obj.tmpl
	objTmplFile string
	objTmpl     = template.Must(template.New("obj").Funcs(template.FuncMap{"vtFlip": vtFlip}).Parse(objTmplFile))

	//go:embed mtl.tmpl
	mtlTmplFile string
	mtlTmpl     = template.Must(template.New("mtl").Parse(mtlTmplFile))
)

type tmplVars struct {
	Version      string
	Flip         types.Flip
	WidthM       float64
	HeightM      float64
	ThickM       float64
	FileFront    string
	FileBack     string
	FileObject   string
	FileMaterial string
}

func WriteObjZip(pc *types.Postcard, w io.Writer, jpegOptions *jpeg.Options) error {
	frontImg, _, err := image.Decode(bytes.NewReader(pc.Front))
	if err != nil {
		return err
	}
	backImg, _, err := image.Decode(bytes.NewReader(pc.Back))
	if err != nil {
		return err
	}

	cmW, _ := pc.Meta.FrontDimensions.CmWidth.Float64()
	cmH, _ := pc.Meta.FrontDimensions.CmHeight.Float64()

	vars := tmplVars{
		Version:      postcards.Version.String(),
		Flip:         pc.Meta.Flip,
		WidthM:       cmW / 100,
		HeightM:      cmH / 100,
		ThickM:       0.0004,
		FileFront:    "front.jpg",
		FileBack:     "back.jpg",
		FileObject:   "postcard.obj",
		FileMaterial: "postcard.mtl",
	}

	ar := zip.NewWriter(w)
	defer ar.Close()

	// Object
	objFile, err := ar.Create(vars.FileObject)
	if err != nil {
		return err
	}
	if err := objTmpl.Execute(objFile, vars); err != nil {
		return err
	}

	// Material
	mtlFile, err := ar.Create(vars.FileMaterial)
	if err != nil {
		return err
	}
	if err := mtlTmpl.Execute(mtlFile, vars); err != nil {
		return err
	}

	frontFile, err := ar.Create(vars.FileFront)
	if err != nil {
		return err
	}
	if err := jpeg.Encode(frontFile, frontImg, jpegOptions); err != nil {
		return err
	}

	backFile, err := ar.Create(vars.FileBack)
	if err != nil {
		return err
	}
	if err := jpeg.Encode(backFile, backImg, jpegOptions); err != nil {
		return err
	}

	return nil
}

// vtFlip shifts the 'vt' declarations by an offset which orients the back of the postcard correctly according to the Flip parameter
func vtFlip(idx int, flip types.Flip) int {
	var offset int
	switch flip {
	case types.FlipBook:
		offset = 0
	case types.FlipRightHand:
		offset = 1
	case types.FlipCalendar:
		offset = 2
	case types.FlipLeftHand:
		offset = 3
	}

	return 4 - ((idx + offset) % 4)
}
