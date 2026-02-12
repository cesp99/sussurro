package ui

import (
	"image/color"
	"time"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
)

type OverlayState int

const (
	StateIdle OverlayState = iota
	StateListening
	StateTranscribing
	StateLoading
)

type OverlayWindow struct {
	window fyne.Window
	state  OverlayState
	
	// Widgets
	container   *fyne.Container
	statusText  *canvas.Text
	waveBars    []*canvas.Rectangle
	shimmerRect *canvas.Rectangle
	
	// Animation channels
	stopAnim chan struct{}
}

func NewOverlayWindow(a fyne.App) *OverlayWindow {
	// Use driver to create a splash window which is borderless by default
	var w fyne.Window
	if drv, ok := a.Driver().(desktop.Driver); ok {
		w = drv.CreateSplashWindow()
		w.SetTitle("Sussurro Overlay")
	} else {
		w = a.NewWindow("Sussurro Overlay")
		// w.SetDecorated(false) // Fyne v2.7.2 should support this, but linter complains. Ignoring for fallback.
	}
	
	w.SetPadded(false)
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(200, 60))

	o := &OverlayWindow{
		window:   w,
		state:    StateIdle,
		stopAnim: make(chan struct{}),
	}

	o.buildUI()
	o.centerOnScreen()
	
	return o
}

func (o *OverlayWindow) Show() {
	o.window.Show()
	o.centerOnScreen()
}

func (o *OverlayWindow) Hide() {
	o.window.Hide()
}

func (o *OverlayWindow) SetState(s OverlayState) {
	o.state = s
	o.updateUI()
}

func (o *OverlayWindow) centerOnScreen() {
	// Simple centering logic - usually needs screen size
	// For now, we put it at bottom center
	// Fyne doesn't give screen size easily without a canvas.
	// We'll rely on OS positioning or set a default location.
	o.window.CenterOnScreen()
}

func (o *OverlayWindow) buildUI() {
	// Background Capsule
	bg := canvas.NewRectangle(color.Black)
	bg.CornerRadius = 30 // Fully rounded for height 60
	bg.FillColor = color.Black
	
	// Content Container
	o.container = fyne.NewContainer() // Manual layout or stack
	
	// Initialize Wave Bars (Hidden by default)
	o.waveBars = make([]*canvas.Rectangle, 5)
	for i := 0; i < 5; i++ {
		rect := canvas.NewRectangle(color.White)
		rect.CornerRadius = 2
		rect.Resize(fyne.NewSize(6, 20))
		rect.Hide()
		o.waveBars[i] = rect
		o.container.Add(rect)
	}

	// Initialize Status Text
	o.statusText = canvas.NewText("transcribing...", color.White)
	o.statusText.TextSize = 18
	o.statusText.Alignment = fyne.TextAlignCenter
	o.statusText.Hide()
	o.container.Add(o.statusText)
	
	// Shimmer/Flare Rect
	o.shimmerRect = canvas.NewRectangle(color.RGBA{255, 255, 255, 50}) // Semi-transparent white
	o.shimmerRect.Resize(fyne.NewSize(40, 60))
	o.shimmerRect.Hide()
	o.container.Add(o.shimmerRect)

	// Combine Background and Content
	// We use a custom renderer or container to layer them
	// Stack: BG -> Content
	// Use NewContainer to avoid "invalid memory address" when passing nil layout to NewContainerWithLayout
	stack := fyne.NewContainer() 
	
	// Custom Layout to ensure BG fills window
	stack.Layout = &overlayLayout{bg: bg, content: o.container, overlay: o}
	stack.Add(bg)
	stack.Add(o.container)
	
	o.window.SetContent(stack)
	o.updateUI()
}

type overlayLayout struct {
	bg      *canvas.Rectangle
	content *fyne.Container
	overlay *OverlayWindow
}

func (l *overlayLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	l.bg.Resize(size)
	
	// Center content
	// Layout wave bars manually
	center := size.Width / 2
	centerY := size.Height / 2
	
	if l.overlay.state == StateListening {
		// Position bars
		gap := float32(10)
		totalWidth := float32(5*6) + float32(4*gap)
		startX := center - totalWidth/2
		
		for i, bar := range l.overlay.waveBars {
			bar.Move(fyne.NewPos(startX + float32(i)*(6+gap), centerY - bar.Size().Height/2))
		}
	} else if l.overlay.state == StateTranscribing || l.overlay.state == StateLoading {
		// Center Text
		textMin := l.overlay.statusText.MinSize()
		l.overlay.statusText.Move(fyne.NewPos(center - textMin.Width/2, centerY - textMin.Height/2))
	}
}

func (l *overlayLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(200, 60)
}

func (o *OverlayWindow) updateUI() {
	// Stop existing animations
	select {
	case o.stopAnim <- struct{}{}:
	default:
	}

	// Reset visibility
	for _, bar := range o.waveBars {
		bar.Hide()
	}
	o.statusText.Hide()
	o.shimmerRect.Hide()

	switch o.state {
	case StateIdle:
		// Maybe show a small dot or logo? Or hidden?
		// User said "always show a capsule overlay"
		// Let's show a static waveform or just the capsule
		for _, bar := range o.waveBars {
			bar.Resize(fyne.NewSize(6, 10)) // Small static bars
			bar.Show()
		}
		
	case StateListening:
		for _, bar := range o.waveBars {
			bar.Show()
		}
		go o.animateWaves()

	case StateTranscribing:
		o.statusText.Text = "Transcribing..."
		o.statusText.Refresh()
		o.statusText.Show()
		o.shimmerRect.Show()
		go o.animateShimmer()

	case StateLoading:
		o.statusText.Text = "Loading..."
		o.statusText.Refresh()
		o.statusText.Show()
		// Optional: Pulse animation or static
	}
	
	o.container.Refresh()
}

func (o *OverlayWindow) animateWaves() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-o.stopAnim:
			return
		case <-ticker.C:
			// Randomize bar heights
			for _, bar := range o.waveBars {
				h := float32(10 + rand.Intn(30)) // 10 to 40
				bar.Resize(fyne.NewSize(6, h))
				// Re-center Y
				// This requires triggering relayout or moving manually here
				// Since we are inside the layout logic, simple move works if we know center
				centerY := float32(30) // 60/2
				bar.Move(fyne.NewPos(bar.Position().X, centerY - h/2))
			}
			o.container.Refresh()
		}
	}
}

func (o *OverlayWindow) animateShimmer() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60fps
	defer ticker.Stop()

	startX := float32(-50)
	endX := float32(250)
	currentX := startX

	for {
		select {
		case <-o.stopAnim:
			return
		case <-ticker.C:
			currentX += 4
			if currentX > endX {
				currentX = startX
			}
			o.shimmerRect.Move(fyne.NewPos(currentX, 0))
			o.container.Refresh()
		}
	}
}
