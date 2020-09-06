
# PiGLEt - Golang EGL bindings for the Rasperry Pi native OpenGL driver 

The _PiGLEt_ Go module provides bindings to the native OpenGLES and EGL libraries available for Raspberry Pi. 

Specifically, PiGLEt provides functions to create an OpenGL context using the native Broadcom drivers `libbcm_host.so` and `libbrcmEGL.so` in `/opt/vc/lib/`. It also provides Go function wrappers for the `libbrcmGLESv2.so` OpenGLES2 library. These libraries make it possible to run OpenGL programs from the command line without requiring the X Window System. OpenGLES programs drawing within a PiGLEt context wll be rendered to the display attached to the HDMI port of the Raspberry Pi.


## Requirements

PiGLEt runs on 

* Raspberry Pi 2 Model B
* Raspberry Pi 3 Model B

PiGLEt does **not** run on

* Raspberry Pi 4 Model B

as the native OpenGL driver is not available for this platform :(





## General Usage



    // +build linux,arm
    
    import (
        "github.com/FEEDFACE-COM/piglet"
        gl "github.com/FEEDFACE-COM/piglet/gles2"
    )

	func main() {
	    piglet.CreateContext()
	    width,height := piglet.GetDisplaySize()
	    piglet.MakeCurrent()
	    gl.InitWithProcAddrFunc( piglet.GetProcAddress )
	
	    gl.Viewport(0, 0, width, height)
	    gl.ClearColor(0., 0., 0., 0.)
	    gl.Clear( gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT )
	    piglet.SwapBuffers()
	    
	    piglet.DestroyContext()
	}





See the `examples/hello-piglet.go` code for a more thorough example.

## Author
_PiGLEt_ is designed and implemented by <folkert@feedface.com>. The GLES2 bindings are generated using glow <https://github.com/go-gl/glow>.


