package main

import (
	"errors"
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

var (
	srcDec, targetDec              []byte
	sourceFilePath, targetFilePath string
	fileFilter                     = storage.NewExtensionFileFilter([]string{".dat"})
	fileBrowserSize                = fyne.NewSize(1000, 800)
)

// Starts a blocking event loop of the GUI.
func RunGUI() {
	a := app.New()
	w := a.NewWindow("Tasukeru")
	w.Resize(fyne.NewSize(800, 600))

	confirmButton := widget.NewButtonWithIcon("Import", theme.ConfirmIcon(), func() {
		if srcDec == nil || targetDec == nil {
			dialog.ShowError(errors.New("missing data"), w)
			return
		}
		targetDec = holocure.MergeSaveDataDecoded(srcDec, targetDec)
		err := holocure.WriteSaveFile(targetFilePath, targetDec)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		dialog.ShowInformation("Result", "Save imported successfully!", w)
	})

	confirmButton.Hide()

	openSourceFileButtonLabel := widget.NewLabel("Select save file to import")
	openSourceFileButtonLabel.Alignment = fyne.TextAlignCenter
	openSourceFileButton := widget.NewButtonWithIcon("Browse source file", theme.FileIcon(), func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}

			u, _ := url.ParseRequestURI(reader.URI().String())
			sourceFilePath = u.Path

			srcDec, err = holocure.Decode(loadFile(reader))
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			openSourceFileButtonLabel.SetText("Selected source file: " + sourceFilePath)

			confirmButton.Show()
		}, w)

		fd.SetFilter(fileFilter)
		fd.Resize(fileBrowserSize)
		fd.Show()
	})

	openTargetFileButtonLabel := widget.NewLabel("Select target file to merge into")
	openTargetFileButtonLabel.Alignment = fyne.TextAlignCenter
	openTargetFileButton := widget.NewButtonWithIcon("Browse target file", theme.FileIcon(), func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}

			u, _ := url.ParseRequestURI(reader.URI().String())
			targetFilePath = u.Path

			targetDec, err = holocure.Decode(loadFile(reader))
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			openTargetFileButtonLabel.SetText("Selected target file: " + targetFilePath)
		}, w)

		fd.SetFilter(fileFilter)
		fd.Resize(fileBrowserSize)
		fd.Show()
	})

	if holocure.DefaultSaveFileExists() {
		targetFilePath = holocure.SaveFilePath()
		targetDec, _ = holocure.DecodeSaveFile(targetFilePath)
		openTargetFileButtonLabel.SetText("Override default save file: " + targetFilePath)
	}

	aboutUrl, err := url.Parse(HomePage)
	check(err)
	aboutHyperLink := widget.NewHyperlink("About", aboutUrl)
	aboutHyperLink.Alignment = fyne.TextAlignTrailing
	versionLabel := widget.NewLabelWithStyle("Version "+Version, fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true})

	box := container.NewVBox(
		openSourceFileButtonLabel,
		openSourceFileButton,
		openTargetFileButtonLabel,
		openTargetFileButton,
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
