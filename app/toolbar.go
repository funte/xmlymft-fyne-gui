package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"xmlymft-fyne-gui/app/mytheme"
)

// Custom toolbar button with a text label.
type ToolbarAction struct {
	Icon        fyne.Resource
	Label       string
	OnActivated func()
}

func (t *ToolbarAction) ToolbarObject() fyne.CanvasObject {
	button := widget.NewButtonWithIcon(t.Label, t.Icon, t.OnActivated)
	// 无边框
	button.Importance = widget.LowImportance
	return button
}

// Custom select entry with a fixed width.
type SelectEntryWithFixedWidth struct {
	widget.SelectEntry

	FixedWidth float32
}

func (e *SelectEntryWithFixedWidth) MinSize() fyne.Size {
	e.ExtendBaseWidget(e)
	return fyne.NewSize(e.FixedWidth, e.Entry.MinSize().Height)
}

// Custom toolbar select entry with an icon.
type ToolbarSelectEntry struct {
	Options  []string
	Entry    *SelectEntryWithFixedWidth
	OnSearch func(keyword string)
}

func (t *ToolbarSelectEntry) ToolbarObject() fyne.CanvasObject {
	t.Entry = &SelectEntryWithFixedWidth{FixedWidth: 180.0 * mytheme.Factor}

	t.Entry.SetPlaceHolder("输入要搜索的专辑")
	t.Entry.OnSubmitted = func(s string) {
		isKeywordStored := false
		for _, keyword := range t.Options {
			if s == keyword {
				isKeywordStored = true
				break
			}
		}
		if isKeywordStored == false {
			t.Options = append(t.Options, s)
			t.Entry.SetOptions(t.Options)
			t.Entry.Refresh()
		}

		t.Entry.SetText(s)
		if t.OnSearch != nil {
			t.OnSearch(s)
		}
	}
	return t.Entry
}

func newToolbar(
	window fyne.Window,
	onOpenFavorite func(),
	onSearch func(keyword string),
) *widget.Toolbar {
	favoriteBtn := &ToolbarAction{theme.StorageIcon(), "收藏", func() {
		if onOpenFavorite != nil {
			onOpenFavorite()
		}
		dialog := dialog.NewInformation("提示", "还没弄好...", window)
		dialog.Show()
	}}
	searchEntry := &ToolbarSelectEntry{
		OnSearch: onSearch,
	}
	// Create toolbar.
	return widget.NewToolbar(
		favoriteBtn,
		widget.NewToolbarSpacer(),
		searchEntry,
	)
}
