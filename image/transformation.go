package image

import (
	"errors"
	"fmt"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	_ "log"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type RegionInstruction struct {
	X      int
	Y      int
	Height int
	Width  int
}

func (rgi RegionInstruction) String() string {
	return fmt.Sprintf("[region] from %d, %d by %d, %d pixels to %d, %d", rgi.X, rgi.Y, rgi.Width, rgi.Height, rgi.X+rgi.Width, rgi.Y+rgi.Height)
}

type SizeInstruction struct {
	Height  int
	Width   int
	Force   bool
	Enlarge bool
}

type RotationInstruction struct {
	Flip  bool
	Angle float64
}

func (ri RotationInstruction) String() string {
	return fmt.Sprintf("[rotation] by %.3f, flip: %t", ri.Angle, ri.Flip)
}

type FormatInstruction struct {
	Format string
}

type Transformation struct {
	level    iiiflevel.Level
	Region   string
	Size     string
	Rotation string
	Quality  string
	Format   string
}

func NewTransformation(level iiiflevel.Level, region string, size string, rotation string, quality string, format string) (*Transformation, error) {

	var ok bool
	var err error

	ok, err = level.Compliance().IsValidImageRegion(region)

	if !ok {
		return nil, err
	}

	ok, err = level.Compliance().IsValidImageSize(size)

	if !ok {
		return nil, err
	}

	ok, err = level.Compliance().IsValidImageRotation(rotation)

	if !ok {
		return nil, err
	}

	ok, err = level.Compliance().IsValidImageQuality(quality)

	if !ok {
		return nil, err
	}

	ok, err = level.Compliance().IsValidImageFormat(format)

	if !ok {
		return nil, err
	}

	// http://iiif.io/api/image/2.1/#canonical-uri-syntax (sigh...)

	if quality == "default" {

		quality, err = level.Compliance().DefaultQuality()

		if err != nil {
			return nil, err
		}
	}

	t := Transformation{
		level:    level,
		Region:   region,
		Size:     size,
		Rotation: rotation,
		Quality:  quality,
		Format:   format,
	}

	return &t, nil
}

func (t *Transformation) ToURI(id string) (string, error) {

	nodes := []string{
		id,
		t.Region,
		t.Size,
		t.Rotation,
		t.Quality,
	}

	for i, v := range nodes {

		// https://github.com/mrap/stringutil/blob/master/urlencode.go

		u, err := url.Parse(v)

		if err != nil {
			return "", err
		}

		nodes[i] = u.String()
	}

	uri := fmt.Sprintf("%s.%s", strings.Join(nodes, "/"), t.Format)
	return uri, nil
}

func (t *Transformation) HasTransformation() bool {

	if t.Region != "full" {
		return true
	}

	if t.Size != "full" {
		return true
	}

	if t.Rotation != "0" {
		return true
	}

	if t.Quality != "default" {
		return true
	}

	return false
}

func (t *Transformation) RegionInstructions(im Image) (*RegionInstruction, error) {

	dims, err := im.Dimensions()

	if err != nil {
		return nil, err
	}

	width := dims.Width()
	height := dims.Height()

	if t.Region == "square" {

		var x int
		var y int

		min := math.Min(float64(width), float64(height))

		if width > height {
			x = (width - height) / 2.
			y = 0

		} else {
			x = 0
			y = (height - width) / 2.
		}

		instruction := RegionInstruction{
			X:      x,
			Y:      y,
			Width:  int(min),
			Height: int(min),
		}

		return &instruction, nil
	}

	arr := strings.Split(t.Region, ":")

	if len(arr) == 1 {

		sizes := strings.Split(arr[0], ",")

		if len(sizes) != 4 {
			message := fmt.Sprintf("Invalid region")
			return nil, errors.New(message)
		}

		x, err := strconv.ParseInt(sizes[0], 10, 64)

		if err != nil {
			return nil, err
		}

		y, err := strconv.ParseInt(sizes[1], 10, 64)

		if err != nil {
			return nil, err
		}

		w, err := strconv.ParseInt(sizes[2], 10, 64)

		if err != nil {
			return nil, err
		}

		h, err := strconv.ParseInt(sizes[3], 10, 64)

		if err != nil {
			return nil, err
		}

		instruction := RegionInstruction{
			Width:  int(w),
			Height: int(h),
			X:      int(x),
			Y:      int(y),
		}

		/*

			Because otherwise you end up with stuff like this:

			./bin/iiif-tile-seed -config config.json -scale-factor 4 184512_5f7f47e5b3c66207_x.jpg
			2016/09/11 07:35:18 184512_5f7f47e5b3c66207_x.jpg
			2016/09/11 07:35:18 184512_5f7f47e5b3c66207_x.jpg/3072,2048,1024,1024/full/0/default.jpg 5.125935ms extract_area: bad extract area
			2016/09/11 07:35:18 184512_5f7f47e5b3c66207_x.jpg/3072,0,1024,1024/full/0/default.jpg 2.667272ms extract_area: bad extract area
			2016/09/11 07:35:18 184512_5f7f47e5b3c66207_x.jpg/3072,1024,1024,1024/full/0/default.jpg 393.638µs extract_area: bad extract area

			It's possible this is best moved in to the specific packages (like vips.go) where this is
			actually a problem... (20160911/thisisaaronland)

		*/

		if instruction.X+instruction.Width > width {
			instruction.Width = width - instruction.X
		}

		if instruction.Y+instruction.Height > height {
			instruction.Height = height - instruction.Y
		}

		return &instruction, nil

	}

	if arr[0] == "pct" {

		sizes := strings.Split(arr[1], ",")

		if len(sizes) != 4 {
			message := fmt.Sprintf("Invalid region", t.Region)
			return nil, errors.New(message)
		}

		px, err := strconv.ParseFloat(sizes[0], 64)

		if err != nil {
			return nil, err
		}

		py, err := strconv.ParseFloat(sizes[1], 64)

		if err != nil {
			return nil, err
		}

		pw, err := strconv.ParseFloat(sizes[2], 64)

		if err != nil {
			return nil, err
		}

		ph, err := strconv.ParseFloat(sizes[3], 64)

		if err != nil {
			return nil, err
		}

		x := int(math.Ceil((float64(width) * px) / 100.))
		y := int(math.Ceil((float64(height) * py) / 100.))

		w := int(math.Ceil((float64(width) * pw) / 100.))
		h := int(math.Ceil((float64(height) * ph) / 100.))

		instruction := RegionInstruction{
			Width:  w,
			Height: h,
			X:      x,
			Y:      y,
		}

		if instruction.X+instruction.Width > width {
			instruction.Width = width - instruction.X
		}

		if instruction.Y+instruction.Height > height {
			instruction.Height = height - instruction.Y
		}

		return &instruction, nil

	} else {
	}

	message := fmt.Sprintf("Unrecognized region")
	return nil, errors.New(message)

}

func (t *Transformation) SizeInstructions(im Image) (*SizeInstruction, error) {

	dims, err := im.Dimensions()

	if err != nil {
		return nil, err
	}

	sizeError := "IIIF 2.1 `size` argument is not recognized: %#v"

	w := 0
	h := 0
	force := false
	enlarge := false

	arr := strings.Split(t.Size, ":")

	if len(arr) == 1 {

		best := strings.HasPrefix(t.Size, "!")
		sizes := strings.Split(strings.Trim(arr[0], "!"), ",")

		if len(sizes) != 2 {
			message := fmt.Sprintf(sizeError, t.Size)
			return nil, errors.New(message)
		}

		wi, err_w := strconv.ParseInt(sizes[0], 10, 64)
		hi, err_h := strconv.ParseInt(sizes[1], 10, 64)

		if err_w != nil && err_h != nil {
			message := fmt.Sprintf(sizeError, t.Size)
			return nil, errors.New(message)

		} else if err_w == nil && err_h == nil {

			w = int(wi)
			h = int(hi)

			if best {
				enlarge = true
			} else {
				force = true
			}

		} else if err_h != nil {

			width := dims.Width()
			height := dims.Height()

			w = int(wi)
			h = height * w / width

		} else {

			width := dims.Width()
			height := dims.Height()

			h = int(hi)
			w = width * h / height

		}

		instruction := SizeInstruction{
			Height:  h,
			Width:   w,
			Enlarge: enlarge,
			Force:   force,
		}

		return &instruction, nil

	} else if arr[0] == "pct" {

		pct, err := strconv.ParseFloat(arr[1], 64)

		if err != nil {
			err := errors.New("invalid size")
			return nil, err
		}

		dims, err := im.Dimensions()

		if err != nil {
			return nil, err
		}

		width := dims.Width()
		height := dims.Height()

		w = int(math.Ceil(pct / 100 * float64(width)))
		h = int(math.Ceil(pct / 100 * float64(height)))

	} else {

		message := fmt.Sprintf(sizeError, t.Size)
		return nil, errors.New(message)
	}

	instruction := SizeInstruction{
		Height:  h,
		Width:   w,
		Enlarge: enlarge,
		Force:   force,
	}

	return &instruction, nil

}

func (t *Transformation) RotationInstructions(im Image) (*RotationInstruction, error) {

	rotationError := "IIIF 2.1 `rotation` argument is not recognized: %#v"

	flip := strings.HasPrefix(t.Rotation, "!")
	rotation := strings.Trim(t.Rotation, "!")

	angle, err := strconv.ParseFloat(rotation, 64)

	if err != nil {
		message := fmt.Sprintf(rotationError, t.Rotation)
		return nil, errors.New(message)

	}

	instruction := RotationInstruction{
		Flip:  flip,
		Angle: angle,
	}

	return &instruction, nil
}

func (t *Transformation) FormatInstructions(im Image) (*FormatInstruction, error) {

	fmt := ""

	compliance := t.level.Compliance()
	spec := compliance.Spec()

	for name, details := range spec.Image.Format {

		re, err := regexp.Compile(details.Match)

		if err != nil {
			return nil, err
		}

		if re.MatchString(t.Format) {
			fmt = name
			break
		}
	}

	if fmt == "" {
		return nil, errors.New("failed to determine format")
	}

	instruction := FormatInstruction{
		Format: fmt,
	}

	return &instruction, nil
}
