package store

import (
	"fyne.io/fyne/v2/widget"
)

type TrackViewItem struct {
	widget.BaseWidget

	title string
	// TODO: play selected track icon.
	playIcon     widget.Icon
	downloadIcon widget.Icon
}
