// +build linux,arm

package piglet

// #cgo CFLAGS:  -I/opt/vc/include 
// #cgo LDFLAGS: -L/opt/vc/lib -ldl -lbcm_host -lbrcmGLESv2 -lbrcmEGL
// #include <stdlib.h>
// #include "piglet.h"
import "C"
import "unsafe"
import "errors"



func CreateContext() error { 
    err := int(C.CreateContext())
    if err != 0 {
        return errors.New("fail to create context!!")
    }
    return nil
}

func MakeCurrent() error { 
    C.MakeCurrent()
    return nil
}

func SwapBuffers() error { 
    C.SwapBuffers()
    return nil
}


func GetDisplaySize() (int32,int32) { 
    return int32(C.GetDisplayWidth()), int32(C.GetDisplayHeight()) 
}


func GetProcAddress(name string) unsafe.Pointer {
    cname := C.CString(name)
    defer C.free(unsafe.Pointer(cname))
    return C.GetProcAddress(cname)    
}
