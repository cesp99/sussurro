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

// NewCapsuleButtonWithIcon creates a new capsule button with an icon
func NewCapsuleButtonWithIcon(label string, icon fyne.Resource, tapped func(), primary bool) *CapsuleButton {
	b := &CapsuleButton{IsPrimary: primary}
	b.Text = label
	b.Icon = icon
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
	
	if b.Icon != nil {
		r.icon = canvas.NewImageFromResource(b.Icon)
	}

	r.Refresh()
	return r
}

type capsuleButtonRenderer struct {
	button *CapsuleButton
	bg     *canvas.Rectangle
	label  *canvas.Text
	icon   *canvas.Image
}

func (r *capsuleButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	
	iconSize := float32(16)
	spacing := float32(8)
	
	textMin := r.label.MinSize()
	
	// Calculate total width of content
	contentWidth := textMin.Width
	if r.icon != nil {
		contentWidth += iconSize + spacing
	}
	
	// Center content
	startX := (size.Width - contentWidth) / 2
	centerY := size.Height / 2
	
	if r.icon != nil {
		r.icon.Resize(fyne.NewSize(iconSize, iconSize))
		r.icon.Move(fyne.NewPos(startX, centerY - iconSize/2))
		r.label.Move(fyne.NewPos(startX + iconSize + spacing, centerY - textMin.Height/2))
	} else {
		r.label.Move(fyne.NewPos(startX, centerY - textMin.Height/2))
	}
}

func (r *capsuleButtonRenderer) MinSize() fyne.Size {
	textMin := r.label.MinSize()
	width := textMin.Width + 40 // Padding
	height := textMin.Height + 20
	
	if r.icon != nil {
		width += 16 + 8 // Icon + Spacing
	}
	
	return fyne.NewSize(width, height)
}

func (r *capsuleButtonRenderer) Refresh() {
	// Colors
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	r.label.Text = r.button.Text
	r.label.TextSize = 14
	r.label.Alignment = fyne.TextAlignCenter
	
	// Capsule shape
	r.bg.CornerRadius = 16 
	r.bg.StrokeWidth = 2
	r.bg.StrokeColor = white

	if r.button.IsPrimary {
		// White background, Black text
		r.bg.FillColor = white
		r.label.Color = black
		// Need to invert icon color? Fyne icons are usually black/white adaptable
		// But canvas.Image from resource is just the image.
		// If icon is white, we might need a black version for primary button?
		// For now, assume icon works or ignore color shift.
	} else {
		// Black background, White border, White text
		r.bg.FillColor = black
		r.label.Color = white
	}
	
	r.bg.Refresh()
	r.label.Refresh()
	if r.icon != nil {
		r.icon.Refresh()
	}
}

func (r *capsuleButtonRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *capsuleButtonRenderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{r.bg, r.label}
	if r.icon != nil {
		objs = append(objs, r.icon)
	}
	return objs
}

func (r *capsuleButtonRenderer) Destroy() {}
