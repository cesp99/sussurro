//go:build darwin

package ui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework QuartzCore -framework CoreVideo

extern void* overlay_create_macos(void);
extern void  overlay_set_state_macos(int state);
extern void  overlay_push_rms_macos(float rms);
extern void  overlay_show_macos(void);
extern void  overlay_hide_macos(void);
*/
import "C"

type darwinOverlay struct{}

func newOverlay() Overlay {
	C.overlay_create_macos()
	return &darwinOverlay{}
}

func (o *darwinOverlay) Show() {
	C.overlay_show_macos()
}

func (o *darwinOverlay) Hide() {
	C.overlay_hide_macos()
}

func (o *darwinOverlay) SetState(state AppState) {
	C.overlay_set_state_macos(C.int(state))
}

func (o *darwinOverlay) PushRMS(rms float32) {
	C.overlay_push_rms_macos(C.float(rms))
}

func (o *darwinOverlay) Close() {
	o.Hide()
}
