package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// CapsuleButton is a custom button with capsule shape
type CapsuleButton struct {
	widget.Button
	IsPrimary bool
}

// NewCapsuleButton creates a new capsule button
func NewCapsuleButton(label string, tapped func(), primary bool) *CapsuleButton {
	b := &CapsuleButton{IsPrimary: primary}
	b.Text = label
	b.OnTapped = tapped
	b.ExtendBaseWidget(b)
	return b
}

// CreateRenderer implements the widget interface
func (b *CapsuleButton) CreateRenderer() fyne.WidgetRenderer {
	r := &capsuleButtonRenderer{
		button: b,
		bg:     canvas.NewRectangle(color.Transparent),
		label:  canvas.NewText(b.Text, color.White),
	}
	r.Refresh()
	return r
}

type capsuleButtonRenderer struct {
	button *CapsuleButton
	bg     *canvas.Rectangle
	label  *canvas.Text
}

func (r *capsuleButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.label.Resize(size)
	r.label.Move(fyne.NewPos(0, 0))
}

func (r *capsuleButtonRenderer) MinSize() fyne.Size {
	return r.label.MinSize().Add(fyne.NewSize(20, 10))
}

func (r *capsuleButtonRenderer) Refresh() {
	// Colors
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	r.label.Text = r.button.Text
	r.label.TextSize = 14
	r.label.Alignment = fyne.TextAlignCenter
	
	// Capsule shape is simulated by rounded corners = height/2
	// But Fyne canvas.Rectangle corner radius is fixed by theme usually.
	// For "Capsule", we can use a high corner radius.
	r.bg.CornerRadius = 20 // Approximate capsule
	r.bg.StrokeWidth = 2
	r.bg.StrokeColor = white

	if r.button.IsPrimary {
		// White background, Black text
		r.bg.FillColor = white
		r.label.Color = black
	} else {
		// Black background, White border, White text
		r.bg.FillColor = black
		r.label.Color = white
	}
	
	r.bg.Refresh()
	r.label.Refresh()
}

func (r *capsuleButtonRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *capsuleButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.label}
}

func (r *capsuleButtonRenderer) Destroy() {}
