package draw

import (
	"fmt"
	"image"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"github.com/nomad-software/meme/font"
)

const (
	fontBorderRadius  = 3.0  // px
	fontLeading       = 1.4  // percentage
	maxFontSize       = 85.0 // pts
	topTextDivisor    = 5.0  // divisor
	bottomTextDivisor = 3.75 // divisor
	imageMargin       = 18.0 // px
)

// NewContext creates a new context for the passed image
func NewContext(img image.Image) *gg.Context {
	return gg.NewContextForImage(img)
}

// TopBanner draws the top text onto the meme.
func TopBanner(ctx *gg.Context, text string) {
	x := float64(ctx.Width()) / 2
	y := imageMargin
	drawText(ctx, text, x, y, 0.5, 0.0, topTextDivisor)
}

// BottomBanner draws the bottom text onto the meme.
func BottomBanner(ctx *gg.Context, text string) {
	x := float64(ctx.Width()) / 2
	y := float64(ctx.Height()) - imageMargin
	drawText(ctx, text, x, y, 0.5, 1.0, bottomTextDivisor)
}

// Draw text onto the meme.
func drawText(ctx *gg.Context, text string, x float64, y float64, ax float64, ay float64, divisor float64) {
	text = strings.ToUpper(text)
	width := float64(ctx.Width()) - (imageMargin * 2)
	height := float64(ctx.Height()) / divisor
	calculateFontSize(ctx, text, width, height)

	// Draw the text border.
	ctx.SetHexColor("#000")
	for angle := 0.0; angle < (2 * math.Pi); angle += 0.35 {
		bx := x + (math.Sin(angle) * fontBorderRadius)
		by := y + (math.Cos(angle) * fontBorderRadius)
		ctx.DrawStringWrapped(text, bx, by, ax, ay, width, fontLeading, gg.AlignCenter)
	}

	// Draw the text itself.
	ctx.SetHexColor("#FFF")
	ctx.DrawStringWrapped(text, x, y, ax, ay, width, fontLeading, gg.AlignCenter)
}

// Dynamically calculate the correct size needed for text.
func calculateFontSize(ctx *gg.Context, text string, width float64, height float64) {
	for size := maxFontSize; size > 20; size-- {
		var rWidth, rHeight float64
		var lWidth, lHeight float64

		// font.Path is validated once, either by font's init() (the
		// embedded font written to a temp file) or by font.SetPath (a
		// user-supplied path/fontconfig match, checked with os.Stat before
		// being assigned). A failure here means that file has become
		// unreadable after that check, which is an environment fault, not
		// something caused by the meme's text or image content — so it's
		// not worth threading an error return through every text-drawing
		// call for.
		if err := ctx.LoadFontFace(font.Path, size); err != nil {
			panic(fmt.Errorf("could not load font file %q: %w", font.Path, err))
		}
		lines := ctx.WordWrap(text, width)

		for _, line := range lines {
			lWidth, lHeight = ctx.MeasureString(line)
			if lWidth > rWidth {
				rWidth = lWidth
			}
		}

		rHeight = (lHeight * fontLeading) * float64(len(lines))

		if rWidth <= width && rHeight <= height {
			break
		}
	}
}
