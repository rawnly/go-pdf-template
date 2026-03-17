package qrcode

import (
	"errors"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/signintech/gopdf"
)

const (
	defaultRadiusPointNum = 6
)

func DrawStyledQRCode(pdf *gopdf.GoPdf, opts Options) error {
	if pdf == nil {
		return errors.New("pdf is nil")
	}
	if opts.Content == "" {
		return errors.New("content is empty")
	}
	if opts.Size <= 0 {
		return errors.New("size must be positive")
	}
	if opts.ShowLogo && opts.IconSize == "" {
		opts.IconSize = IconSizeMedium
	}

	cornerSize := opts.CornerSizeModules
	if cornerSize <= 0 {
		cornerSize = DefaultCornerSizeModules
	}
	innerGap := opts.InnerGapModules
	if innerGap <= 0 {
		innerGap = DefaultInnerGapModules
	}

	code, err := qr.Encode(opts.Content, qr.H, qr.Auto)
	if err != nil {
		return err
	}

	bounds := code.Bounds()
	qrSize := bounds.Dx()
	if qrSize <= 0 {
		return errors.New("invalid qr size")
	}

	moduleSize := opts.Size / float64(qrSize)
	if moduleSize <= 0 {
		return errors.New("invalid module size")
	}

	setFill(pdf, colorWhite)
	drawRect(pdf, opts.X, opts.Y, opts.Size, opts.Size, 0)

	cornerRadius := moduleSize * 2 / 5
	innerCornerRadius := moduleSize * 3 / 10
	innerMostRadius := moduleSize * 1 / 5

	drawCorner := func(x int, y int) {
		px := opts.X + float64(x)*moduleSize
		py := opts.Y + float64(y)*moduleSize

		setFill(pdf, colorDark)
		drawRect(pdf, px, py, moduleSize*float64(cornerSize), moduleSize*float64(cornerSize), cornerRadius)

		setFill(pdf, colorWhite)
		drawRect(pdf, px+moduleSize, py+moduleSize, moduleSize*float64(cornerSize-2), moduleSize*float64(cornerSize-2), innerCornerRadius)

		setFill(pdf, colorDark)
		drawRect(pdf, px+moduleSize*2, py+moduleSize*2, moduleSize*float64(cornerSize-4), moduleSize*float64(cornerSize-4), innerMostRadius)
	}

	drawCorner(0, 0)
	drawCorner(0, qrSize-cornerSize)
	drawCorner(qrSize-cornerSize, 0)

	setFill(pdf, colorDark)
	for x := 0; x < qrSize; x++ {
		for y := 0; y < qrSize; y++ {
			if isInCorner(x, y, qrSize, cornerSize) {
				continue
			}
			if isInInnerGap(x, y, qrSize, cornerSize, innerGap) {
				continue
			}
			if !isDark(code, x, y) {
				continue
			}

			px := opts.X + float64(x)*moduleSize
			py := opts.Y + float64(y)*moduleSize
			drawRect(pdf, px, py, moduleSize, moduleSize, moduleSize/5)
		}
	}

	if opts.ShowLogo {
		drawLogo(pdf, opts)
	}

	return nil
}

func isDark(code barcode.Barcode, x, y int) bool {
	color := code.At(x, y)
	r, g, b, a := color.RGBA()
	if a == 0 {
		return false
	}
	return (r+g+b)/3 < 0x8000
}

func isInCorner(x, y, size, corner int) bool {
	if x < corner && y < corner {
		return true
	}
	if x < corner && y >= size-corner {
		return true
	}
	if x >= size-corner && y < corner {
		return true
	}
	return false
}

func isInInnerGap(x, y, size, corner, gap int) bool {
	if gap <= 0 {
		return false
	}
	innerStart := corner
	innerEnd := corner + gap

	if x >= innerStart && x < innerEnd && y >= innerStart && y < innerEnd {
		return true
	}
	if x >= innerStart && x < innerEnd && y >= size-innerEnd && y < size-innerStart {
		return true
	}
	if x >= size-innerEnd && x < size-innerStart && y >= innerStart && y < innerEnd {
		return true
	}
	return false
}

func drawRect(pdf *gopdf.GoPdf, x, y, w, h, radius float64) {
	if radius <= 0 {
		pdf.RectFromUpperLeftWithStyle(x, y, w, h, "F")
		return
	}
	_ = pdf.Rectangle(x, y, x+w, y+h, "F", radius, defaultRadiusPointNum)
}

func setFill(pdf *gopdf.GoPdf, c gopdf.RGBColor) {
	pdf.SetFillColor(c.R, c.G, c.B)
}

func drawLogo(pdf *gopdf.GoPdf, opts Options) {
	logoSize, offsetX, offsetY := logoPlacement(opts.Size, opts.IconSize)
	if logoSize <= 0 {
		return
	}

	logoX := opts.X + offsetX
	logoY := opts.Y + offsetY
	viewBox := 100.0
	squareSize := logoSize * 14.0 / viewBox
	squareRadius := logoSize * 2.0 / viewBox

	setFill(pdf, colorRed)
	drawRect(pdf, logoX, logoY, squareSize, squareSize, squareRadius)

	setFill(pdf, colorWhite)
	polygons := satispayLogoPolygons()
	for _, polygon := range polygons {
		points := make([]gopdf.Point, 0, len(polygon))
		for _, p := range polygon {
			points = append(points, gopdf.Point{
				X: logoX + (p.x/viewBox)*logoSize,
				Y: logoY + (p.y/viewBox)*logoSize,
			})
		}
		if len(points) > 1 {
			pdf.Polygon(points, "F")
		}
	}
}

func logoPlacement(size float64, iconSize IconSize) (logoSize, offsetX, offsetY float64) {
	if size <= 0 {
		return 0, 0, 0
	}

	switch iconSize {
	case IconSizeSmall:
		logoSize = size * 2
		offsetX = size/2 - size/7
		offsetY = size/2 - size/7
	case IconSizeLarge:
		logoSize = size * 1.5
		offsetX = size/2 - size/7 + 2.7
		offsetY = size/2 - size/7 + 2.7
	case IconSizeMedium:
		fallthrough
	default:
		logoSize = size * 1.7
		offsetX = size/2 - size/7 + 2.4
		offsetY = size/2 - size/7 + 2.4
	}

	return logoSize, offsetX, offsetY
}
