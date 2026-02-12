package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// SussurroTheme implements a high-contrast Black & White theme
type SussurroTheme struct{}

var (
	ColorBlack = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	ColorWhite = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	ColorGray  = color.RGBA{R: 40, G: 40, B: 40, A: 255}
)

func (t *SussurroTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground, theme.ColorNameOverlayBackground, theme.ColorNameMenuBackground:
		return ColorBlack
	case theme.ColorNameForeground:
		return ColorWhite
	case theme.ColorNameButton, theme.ColorNameInputBackground:
		return ColorBlack
	case theme.ColorNamePrimary:
		return ColorWhite
	case theme.ColorNameHover:
		return ColorGray
	case theme.ColorNameFocus:
		return ColorWhite
	case theme.ColorNameShadow:
		return color.Transparent
	case theme.ColorNameScrollBar:
		return ColorWhite
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t *SussurroTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *SussurroTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *SussurroTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 12
	case theme.SizeNameScrollBarSmall:
		return 4
	case theme.SizeNameText:
		return 14
	case theme.SizeNameInputRadius:
		return 8
	default:
		return theme.DefaultTheme().Size(name)
	}
}
