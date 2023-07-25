package widget

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/internal/cache"
	col "fyne.io/fyne/v2/internal/color"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

type progressRenderer struct {
	widget.BaseRenderer
	//background, bar *canvas.Rectangle
	background             *canvas.Rectangle
	bar                    *canvas.LinearGradient
	barHighlight           *canvas.LinearGradient
	barShadow              *canvas.LinearGradient
	label                  *canvas.Text
	progress               *ProgressBar
	labelBackground        *canvas.Rectangle // new field for the label background
	leftTopCornerImage     *canvas.Image
	leftBottomCornerImage  *canvas.Image
	rightTopCornerImage    *canvas.Image
	rightBottomCornerImage *canvas.Image
	rightTopCornerFix      *canvas.Image
	rightBottomCornerFix   *canvas.Image
}

// MinSize calculates the minimum size of a progress bar.
// This is simply the "100%" label size plus padding.
func (p *progressRenderer) MinSize() fyne.Size {
	var tsize fyne.Size
	if text := p.progress.TextFormatter; text != nil {
		tsize = fyne.MeasureText(text(), p.label.TextSize, p.label.TextStyle)
	} else {
		tsize = fyne.MeasureText("100%", p.label.TextSize, p.label.TextStyle)
	}

	return fyne.NewSize(tsize.Width+theme.InnerPadding()*2, tsize.Height+theme.InnerPadding()*2)
}

func (p *progressRenderer) updateBar() {
	if p.progress.Value < p.progress.Min {
		p.progress.Value = p.progress.Min
	}
	if p.progress.Value > p.progress.Max {
		p.progress.Value = p.progress.Max
	}

	delta := float32(p.progress.Max - p.progress.Min)
	ratio := float32(p.progress.Value-p.progress.Min) / delta

	if text := p.progress.TextFormatter; text != nil {
		p.label.Text = text()
	} else {
		p.label.Text = strconv.Itoa(int(ratio*100)) + "%"
	}

	size := p.progress.Size()
	p.bar.Resize(fyne.NewSize(size.Width*ratio, size.Height))
	p.barHighlight.Resize(fyne.NewSize(size.Width*ratio, size.Height/5))
	p.barShadow.Resize(fyne.NewSize(size.Width*ratio, size.Height/5))
	p.barShadow.Move(fyne.NewPos(0, size.Height-size.Height/5))
	p.rightTopCornerImage.Move(fyne.NewPos(p.bar.Size().Width-3, 0))
	p.rightBottomCornerImage.Move(fyne.NewPos(p.bar.Size().Width-3, size.Height-3))
}

// Layout the components of the check widget
func (p *progressRenderer) Layout(size fyne.Size) {
	p.background.Resize(size)
	p.labelBackground.Resize(fyne.NewSize(55, 25)) // resize the label background to be slightly larger
	// center the label background
	p.labelBackground.Move(fyne.NewPos((size.Width-p.labelBackground.Size().Width)/2, (size.Height-p.labelBackground.Size().Height)/2))
	p.label.Resize(size)
	p.updateBar()

	imageSize := fyne.NewSize(3, 3) // Adjust the size of the left corner image as needed

	// Position the left corner image at the top-left corner of the progress bar
	p.leftTopCornerImage.Resize(imageSize)
	p.leftTopCornerImage.Move(fyne.NewPos(0, 0))
	p.leftBottomCornerImage.Resize(imageSize)
	p.leftBottomCornerImage.Move(fyne.NewPos(0, size.Height-imageSize.Height))

	if p.progress.Value == 0 {
		p.rightTopCornerImage.Hide()
		p.rightBottomCornerImage.Hide()
	} else {
		p.rightTopCornerImage.Show()
		p.rightBottomCornerImage.Show()
	}

	// Position the right corner image at the top-right corner of the progress bar
	p.rightTopCornerImage.Resize(imageSize)
	p.rightTopCornerImage.Move(fyne.NewPos(p.bar.Size().Width-3, 0))
	p.rightBottomCornerImage.Resize(imageSize)
	p.rightBottomCornerImage.Move(fyne.NewPos(p.bar.Size().Width-3, size.Height-imageSize.Height))

	p.rightTopCornerFix.Resize(imageSize)
	p.rightTopCornerFix.Move(fyne.NewPos(size.Width-3, 0))
	p.rightBottomCornerFix.Resize(imageSize)
	p.rightBottomCornerFix.Move(fyne.NewPos(size.Width-3, size.Height-imageSize.Height))
}

// applyTheme updates the progress bar to match the current theme
func (p *progressRenderer) applyTheme() {
	p.background.FillColor = theme.DisabledButtonColor()
	p.background.CornerRadius = theme.InputRadiusSize()
	//p.bar.FillColor = theme.PrimaryColor()
	//p.bar.CornerRadius = theme.InputRadiusSize()
	p.labelBackground.FillColor = translucentBackgroundColor() // set the label background color
	p.labelBackground.CornerRadius = theme.InputRadiusSize()   // set the label background corner radius
	p.label.Color = theme.ForegroundColor()
	p.label.TextSize = theme.TextSize()
	p.label.TextStyle.Bold = true
	if p.leftTopCornerImage != nil {
		p.leftTopCornerImage.FillMode = canvas.ImageFillOriginal
		//p.leftCornerImage.SetMinSize(fyne.NewSize(theme.InputRadiusSize(), p.progress.Size().Height))
	}
	if p.leftBottomCornerImage != nil {
		p.leftBottomCornerImage.FillMode = canvas.ImageFillOriginal
		//p.leftCornerImage.SetMinSize(fyne.NewSize(theme.InputRadiusSize(), p.progress.Size().Height))
	}

	if p.rightTopCornerImage != nil {
		p.rightTopCornerImage.FillMode = canvas.ImageFillOriginal
		//p.rightCornerImage.SetMinSize(fyne.NewSize(theme.InputRadiusSize(), p.progress.Size().Height))
	}
	if p.rightBottomCornerImage != nil {
		p.rightBottomCornerImage.FillMode = canvas.ImageFillOriginal
		//p.rightCornerImage.SetMinSize(fyne.NewSize(theme.InputRadiusSize(), p.progress.Size().Height))
	}

	if p.rightTopCornerFix != nil {
		p.rightTopCornerFix.FillMode = canvas.ImageFillOriginal
		//p.rightCornerImage.SetMinSize(fyne.NewSize(theme.InputRadiusSize(), p.progress.Size().Height))
	}
	if p.rightBottomCornerFix != nil {
		p.rightBottomCornerFix.FillMode = canvas.ImageFillOriginal
		//p.rightCornerImage.SetMinSize(fyne.NewSize(theme.InputRadiusSize(), p.progress.Size().Height))
	}
	if p.progress.Value == 0 {
		p.rightTopCornerImage.Hide()
		p.rightBottomCornerImage.Hide()
	} else {
		p.rightTopCornerImage.Show()
		p.rightBottomCornerImage.Show()
	}
}

func (p *progressRenderer) Refresh() {
	p.applyTheme()
	p.updateBar()
	p.background.Refresh()
	p.bar.Refresh()
	p.labelBackground.Refresh() // refresh the label background
	p.leftTopCornerImage.Refresh()
	p.leftBottomCornerImage.Refresh()
	p.rightBottomCornerImage.Refresh()
	p.rightTopCornerImage.Refresh()
	p.rightTopCornerFix.Refresh()
	p.rightBottomCornerFix.Refresh()
	if p.progress.Value == 0 {
		p.rightTopCornerImage.Hide()
		p.rightBottomCornerImage.Hide()
	} else {
		p.rightTopCornerImage.Show()
		p.rightBottomCornerImage.Show()
	}
	canvas.Refresh(p.progress.super())
}

// ProgressBar widget creates a horizontal panel that indicates progress
type ProgressBar struct {
	BaseWidget

	Min, Max, Value float64

	// TextFormatter can be used to have a custom format of progress text.
	// If set, it overrides the percentage readout and runs each time the value updates.
	//
	// Since: 1.4
	TextFormatter func() string `json:"-"`

	binder basicBinder
}

// Bind connects the specified data source to this ProgressBar.
// The current value will be displayed and any changes in the data will cause the widget to update.
//
// Since: 2.0
func (p *ProgressBar) Bind(data binding.Float) {
	p.binder.SetCallback(p.updateFromData)
	p.binder.Bind(data)
}

// SetValue changes the current value of this progress bar (from p.Min to p.Max).
// The widget will be refreshed to indicate the change.
func (p *ProgressBar) SetValue(v float64) {
	p.Value = v
	p.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (p *ProgressBar) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (p *ProgressBar) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)
	if p.Min == 0 && p.Max == 0 {
		p.Max = 1.0
	}

	//background := canvas.NewVerticalGradient(progressBackgroundColor(), &color.NRGBA{R: uint8(200), G: uint8(0), B: uint8(200), A: uint8(50)})
	//background.Angle = 45
	background := canvas.NewRectangle(theme.DisabledButtonColor())
	background.CornerRadius = theme.InputRadiusSize()
	bar := canvas.NewVerticalGradient(theme.PrimaryColor(), &color.NRGBA{R: uint8(100), G: uint8(0), B: uint8(200), A: uint8(255)}) // or theme.PrimaryColor() and theme.PrimaryColorNamed("purple") | &color.NRGBA{R: uint8(200), G: uint8(0), B: uint8(200), A: uint8(255)}
	bar.Angle = 45
	highlight := 50 // 30
	shadow := 100   // 120
	barHighlight := canvas.NewVerticalGradient(&color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: uint8(highlight)}, &color.NRGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(0)})
	barShadow := canvas.NewVerticalGradient(&color.NRGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(0)}, &color.NRGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(shadow)})
	//bar := canvas.NewRectangle(theme.PrimaryColor())
	//bar.CornerRadius = theme.InputRadiusSize()
	labelBackground := canvas.NewRectangle(translucentBackgroundColor()) // create the label background
	labelBackground.CornerRadius = theme.InputRadiusSize()
	label := canvas.NewText("0%", theme.BackgroundColor())
	label.Alignment = fyne.TextAlignCenter

	leftTopCornerImage := canvas.NewImageFromFile(`C:\Users\john\Documents\Python Scripts\- go\encoder\corners\CornerRadiusTopLeft.png`)          // Replace with the file path for your left corner image
	leftBottomCornerImage := canvas.NewImageFromFile(`C:\Users\john\Documents\Python Scripts\- go\encoder\corners\CornerRadiusBottomLeft.png`)    // Replace with the file path for your right corner image
	rightTopCornerImage := canvas.NewImageFromFile(`C:\Users\john\Documents\Python Scripts\- go\encoder\corners\CornerRadiusTopRight.png`)        // Replace with the file path for your left corner image
	rightBottomCornerImage := canvas.NewImageFromFile(`C:\Users\john\Documents\Python Scripts\- go\encoder\corners\CornerRadiusBottomRight.png`)  // Replace with the file path for your right corner image
	rightTopCornerFix := canvas.NewImageFromFile(`C:\Users\john\Documents\Python Scripts\- go\encoder\corners\CornerRadiusTopRightFix.png`)       // Replace with the file path for your left corner image
	rightBottomCornerFix := canvas.NewImageFromFile(`C:\Users\john\Documents\Python Scripts\- go\encoder\corners\CornerRadiusBottomRightFix.png`) // Replace with the file path for your right corner image

	return &progressRenderer{widget.NewBaseRenderer([]fyne.CanvasObject{background, bar, barHighlight, barShadow, labelBackground, label, leftTopCornerImage, leftBottomCornerImage, rightTopCornerImage, rightBottomCornerImage, rightTopCornerFix, rightBottomCornerFix}), background, bar, barHighlight, barShadow, label, p, labelBackground, leftTopCornerImage, leftBottomCornerImage, rightTopCornerImage, rightBottomCornerImage, rightTopCornerFix, rightBottomCornerFix}
}

func translucentBackgroundColor() color.Color {
	r, g, b, a := col.ToNRGBA(theme.DisabledButtonColor())
	faded := uint8(a) / 3
	return &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: faded}
}

func darkenColor() color.Color {
	r, g, b, a := col.ToNRGBA(theme.PrimaryColor())
	faded := uint8(a) - 0x80
	return &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: faded}
}

// Unbind disconnects any configured data source from this ProgressBar.
// The current value will remain at the last value of the data source.
//
// Since: 2.0
func (p *ProgressBar) Unbind() {
	p.binder.Unbind()
}

// NewProgressBar creates a new progress bar widget.
// The default Min is 0 and Max is 1, Values set should be between those numbers.
// The display will convert this to a percentage.
func NewProgressBar() *ProgressBar {
	p := &ProgressBar{Min: 0, Max: 1}

	cache.Renderer(p).Layout(p.MinSize())
	return p
}

// NewProgressBarWithData returns a progress bar connected with the specified data source.
//
// Since: 2.0
func NewProgressBarWithData(data binding.Float) *ProgressBar {
	p := NewProgressBar()
	p.Bind(data)

	return p
}

func progressBackgroundColor() color.Color {
	r, g, b, a := col.ToNRGBA(theme.PrimaryColor())
	faded := uint8(a) / 5
	return &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: faded}
}

func (p *ProgressBar) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	floatSource, ok := data.(binding.Float)
	if !ok {
		return
	}

	val, err := floatSource.Get()
	if err != nil {
		fyne.LogError("Error getting current data value", err)
		return
	}
	p.SetValue(val)
}
