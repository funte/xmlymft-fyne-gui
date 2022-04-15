package utils

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func AbortOnError(err error, window fyne.Window) {
	dlg := dialog.NewError(err, window)
	dlg.SetOnClosed(func() { os.Exit(1) })
	dlg.Show()
}
