package main

import (
	"io"
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/DaniruKun/tasukeru/holocure"
)

// Starts a blocking event loop of the GUI.
func RunGUI() {
	a := app.New()
	w := a.NewWindow("Tasukeru")
	w.Resize(fyne.NewSize(800, 600))

	var srcDec []byte

	targetSaveFilePath := holocure.SaveFilePath()

	confirmButton := widget.NewButtonWithIcon("Import", theme.ConfirmIcon(), func() {
		targetDec := mergeFiles(srcDec, targetSaveFilePath)
		err := holocure.WriteSaveFile(targetSaveFilePath, targetDec)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		dialog.NewInformation("Result", "Save imported successfully!", w).Show()
	})
	confirmButton.Hide()

	openFileButtonLabel := widget.NewLabel("Select save file to import")
	openFileButtonLabel.Alignment = fyne.TextAlignCenter
	openFileButton := widget.NewButton("Open file", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}

			srcDec, err = holocure.Decode(loadFile(reader))
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			confirmButton.Show()
		}, w)

		fd.SetFilter(storage.NewExtensionFileFilter([]string{".dat"}))
		fd.Resize(fyne.NewSize(1000, 800))
		fd.Show()
	})
	openFileButton.Icon = theme.FileIcon()

	aboutUrl, err := url.Parse(HomePage)
	check(err)
	aboutHyperLink := widget.NewHyperlink("About", aboutUrl)
	aboutHyperLink.Alignment = fyne.TextAlignTrailing
	versionLabel := widget.NewLabelWithStyle("Version "+Version, fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true})

	box := container.NewVBox(
		openFileButtonLabel,
		openFileButton,
		confirmButton,
		aboutHyperLink,
		versionLabel,
	)

	w.SetContent(box)
	w.ShowAndRun()
}

func loadFile(f fyne.URIReadCloser) []byte {
	data, err := io.ReadAll(f)
	if err != nil {
		fyne.LogError("Failed to load file data", err)
		return nil
	}
	return data
}

func mergeFiles(srcDec []byte, targetSaveFilePath string) []byte {
	var start, end int

	start, end = holocure.FindSaveBlockStartEnd(&srcDec)
	srcSaveBlock := srcDec[start : end+1]

	targetDec, err := holocure.DecodeSaveFile(targetSaveFilePath)
	check(err)

	start, _ = holocure.FindSaveBlockStartEnd(&targetDec)

	for i, char := range srcSaveBlock {
		targetOffset := start + i

		if targetOffset >= len(targetDec) {
			targetDec = append(targetDec, char)
		} else {
			targetDec[targetOffset] = char
		}
	}

	return targetDec
}
