package store

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type AlbumViewItem struct {
	widget.BaseWidget

	cover fyne.Resource
	title string
	// TODO: collect to favorite icon.
	collectIcon widget.Icon
}
