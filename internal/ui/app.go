package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"mkhtpasswd/internal/htpasswd"
	"mkhtpasswd/internal/password"
)

func Run() error {
	a := app.New()
	a.Settings().SetTheme(theme.DefaultTheme())

	w := a.NewWindow("mkhtpasswd")
	w.Resize(fyne.NewSize(420, 440))

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("Username")

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("Password")

	resultLabel := widget.NewLabel("")
	resultLabel.Selectable = true
	resultLabel.Wrapping = fyne.TextWrapBreak

	lengthEntry := widget.NewEntry()
	lengthEntry.SetText("20")

	lowerCheck := widget.NewCheck("a–z", nil)
	lowerCheck.SetChecked(true)
	upperCheck := widget.NewCheck("A–Z", nil)
	upperCheck.SetChecked(true)
	numberCheck := widget.NewCheck("0–9", nil)
	numberCheck.SetChecked(true)
	symbolCheck := widget.NewCheck("!@#$…", nil)

	randomOpts := container.NewVBox(
		widget.NewLabel("Random password options"),
		container.NewBorder(nil, nil, widget.NewLabel("Length:"), nil, lengthEntry),
		lowerCheck,
		upperCheck,
		numberCheck,
		symbolCheck,
	)
	randomOpts.Hide()

	toggleOptsBtn := widget.NewButton("Random password options ▸", nil)
	genPassBtn := widget.NewButton("Generate random password", nil)
	genEntryBtn := widget.NewButton("Generate .htpasswd entry", nil)
	copyBtn := widget.NewButton("Copy", nil)
	genEntryBtn.Disable()

	toggleOptsBtn.OnTapped = func() {
		if randomOpts.Visible() {
			randomOpts.Hide()
			toggleOptsBtn.SetText("Random password options ▸")
		} else {
			randomOpts.Show()
			toggleOptsBtn.SetText("Random password options ▾")
		}
	}

	genPassBtn.OnTapped = func() {
		length := 20
		if l, err := parseLength(lengthEntry.Text); err == nil && l > 0 {
			length = l
		}

		opts := password.Options{
			Length:    length,
			Lowercase: lowerCheck.Checked,
			Uppercase: upperCheck.Checked,
			Numbers:   numberCheck.Checked,
			Symbols:   symbolCheck.Checked,
		}

		pwd, err := password.GenerateRandomPassword(opts)
		if err != nil {
			passEntry.SetText("")
			return
		}
		passEntry.SetText(pwd)
	}

	genEntryBtn.OnTapped = func() {
		line, err := htpasswd.Generate(userEntry.Text, passEntry.Text)
		if err != nil {
			resultLabel.SetText(fmt.Sprintf("Error: %v", err))
			return
		}
		resultLabel.SetText(line)
	}

	copyBtn.OnTapped = func() {
		text := resultLabel.Text
		if text == "" {
			return
		}
		cb := w.Clipboard()
		cb.SetContent(text)

		copyBtn.SetText("Copied!")
		go func() {
			time.Sleep(2 * time.Second)
			fyne.Do(func() {
				copyBtn.SetText("Copy")
			})
		}()
	}

	updateGenerateState := func() {
		if userEntry.Text == "" || passEntry.Text == "" {
			genEntryBtn.Disable()
		} else {
			genEntryBtn.Enable()
		}
	}
	userEntry.OnChanged = func(string) { updateGenerateState() }
	passEntry.OnChanged = func(string) { updateGenerateState() }

	logo := canvas.NewImageFromResource(LogoPng)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(150, 150))

	w.SetIcon(LogoPng)

	resultBox := container.NewVBox(
		widget.NewLabelWithStyle("Result", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		resultLabel,
	)

	form := container.New(layout.NewFormLayout(),
		widget.NewLabel("Username"), userEntry,
		widget.NewLabel("Password"), passEntry,
	)

	content := container.NewVBox(
		logo,
		form,
		genPassBtn,
		toggleOptsBtn,
		randomOpts,
		widget.NewSeparator(),
		genEntryBtn,
		resultBox,
		copyBtn,
	)

	w.SetContent(content)
	w.ShowAndRun()
	return nil
}

func parseLength(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return 0, err
	}
	if n > password.MaxPasswordLength {
		n = password.MaxPasswordLength
	}
	if n < 1 {
		return 0, fmt.Errorf("length must be positive")
	}
	return n, nil
}