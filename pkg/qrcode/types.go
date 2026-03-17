package qrcode

type IconSize string

const (
	IconSizeSmall  IconSize = "sm"
	IconSizeMedium IconSize = "md"
	IconSizeLarge  IconSize = "lg"
)

const (
	DefaultCornerSizeModules = 7
	DefaultInnerGapModules   = 7
)

type Options struct {
	X, Y, Size        float64
	Content           string
	IconSize          IconSize
	ShowLogo          bool
	CornerSizeModules int
	InnerGapModules   int
}
