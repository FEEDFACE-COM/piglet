// +build linux,arm

// PiGLEt provides functions to create, use and destroy an EGL rendering
// context. OpenGLES programs drawing within in this context will be rendered
// fullscreen to the display attached to the HDMI port of the Raspberry Pi 3.
package piglet

// #cgo CFLAGS:  -I/opt/vc/include
// #cgo LDFLAGS: -L/opt/vc/lib -ldl -lbcm_host -lbrcmEGL -lbrcmGLESv2
// #include <stdlib.h>
// #include "piglet.h"
import "C"
import "unsafe"
import "errors"
import "github.com/FEEDFACE-COM/piglet/gles2"






// Create a new EGL rendering context
func CreateContext() error {
	err := int(C.CreateContext())
	if err != 0 {
		return errors.New("fail to create context!!")
	}
	return nil
}

// Attach the EGL rendering context to the EGL surface
func MakeCurrent() error {
	C.MakeCurrent()
	return nil
}

// Post the EGL surface color buffer to the native display
func SwapBuffers() error {
	C.SwapBuffers()
	return nil
}

// Return the size of the native display, in pixels
func GetDisplaySize() (int32, int32) {
	return int32(C.GetDisplayWidth()), int32(C.GetDisplayHeight())
}

// Return a GL or an EGL extension function
func GetProcAddress(name string) unsafe.Pointer {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return C.GetProcAddress(cname)
}

// Destroy an EGL rendering context
func DestroyContext() error {
	err := int(C.DestroyContext())
	if err != 0 {
		return errors.New("fail to destroy context!!")
	}
	return nil
}


func Loop() bool {
	return true
}

func ErrorString(e uint32) string {
	return gles2.ErrorString(e)
}


