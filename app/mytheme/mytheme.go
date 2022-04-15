package mytheme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"xmlymft-fyne-gui/resources"
)

// Global scale factor.
const Factor = 1.0

type Theme struct{}

var _ fyne.Theme = (*Theme)(nil)

func (t Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t Theme) Font(style fyne.TextStyle) fyne.Resource {
	return resources.Msyh
}

func (t Theme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t Theme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name) * Factor
}
