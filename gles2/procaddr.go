
package gles2

import "unsafe"

// this is just to satisfy gles2.Init()! instead call:
// gles2.InitWithProcAddrFunc( piglet.GetProcAddress )

func getProcAddress(string) unsafe.Pointer { 
    panic("DO NOT CALL gles2.getProcAddress!! IT WILL CRASH!!")
    return unsafe.Pointer(nil) 
}

