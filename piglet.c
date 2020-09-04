

#include <stdio.h>
#include <bcm_host.h>
#include <EGL/egl.h>
#include <GLES2/gl2.h>

#include "piglet.h"




#ifdef PIGLET_DEBUG
#define PIGLET_PRINT(...) do { fprintf(stderr,"PiGLEt "); fprintf(stderr,__VA_ARGS__); fprintf(stderr,"\n"); } while (0)
#define PIGLET_ERROR(...) do { fprintf(stderr,"PiGLEt ERROR: %s#%d ",__FILE__,__LINE__); fprintf(stderr,__VA_ARGS__); fprintf(stderr,"\n"); } while (0)
#define PIGLET_CHECK(...) do { EGLint err = eglGetError(); if ( err != EGL_SUCCESS ) { fprintf(stderr,"PiGLEt ERROR: %s#%d ",__FILE__,__LINE__); fprintf(stderr,__VA_ARGS__); fprintf(stderr," %s\n",eglGetErrorString(err)); } } while (0)

const char* 
eglGetErrorString(EGLint error) {
    switch (error) {
        case EGL_SUCCESS:               return "EGL_SUCCESS";
        case EGL_NOT_INITIALIZED:       return "EGL_NOT_INITIALIZED";
        case EGL_BAD_ACCESS:            return "EGL_BAD_ACCESS";
        case EGL_BAD_ALLOC:             return "EGL_BAD_ALLOC";
        case EGL_BAD_ATTRIBUTE:         return "EGL_BAD_ATTRIBUTE";
        case EGL_BAD_CONTEXT:           return "EGL_BAD_CONTEXT";
        case EGL_BAD_CONFIG:            return "EGL_BAD_CONFIG";
        case EGL_BAD_CURRENT_SURFACE:   return "EGL_BAD_CURRENT_SURFACE";
        case EGL_BAD_DISPLAY:           return "EGL_BAD_DISPLAY";
        case EGL_BAD_SURFACE:           return "EGL_BAD_SURFACE";
        case EGL_BAD_MATCH:             return "EGL_BAD_MATCH";
        case EGL_BAD_PARAMETER:         return "EGL_BAD_PARAMETER";
        case EGL_BAD_NATIVE_PIXMAP:     return "EGL_BAD_NATIVE_PIXMAP";
        case EGL_BAD_NATIVE_WINDOW:     return "EGL_BAD_NATIVE_WINDOW";
        case EGL_CONTEXT_LOST:          return "EGL_CONTEXT_LOST";
        default:                        return "UNKNOWN";
    }
}

#else
#define PIGLET_PRINT(...) do {} while (0)
#define PIGLET_ERROR(...) do {} while (0)
#define PIGLET_CHECK(...) do {} while (0)
#endif


static uint32_t width  = 0;
static uint32_t height = 0;


static EGLDisplay display;
static EGLContext context;
static EGLSurface surface;

static DISPMANX_ELEMENT_HANDLE_T dispman_element;
static DISPMANX_DISPLAY_HANDLE_T dispman_display;


int GetDisplayWidth()  { return (int) width;  }
int GetDisplayHeight() { return (int) height; }

int
CreateContext()
{
    
    // see /opt/vc/src/hello_pi/hello_triangle/triangle.c 
    
    
    static EGL_DISPMANX_WINDOW_T     native_window;
    static DISPMANX_UPDATE_HANDLE_T  dispman_update;

    
    VC_RECT_T src_rect;
    VC_RECT_T dst_rect;
    
    static const EGLint attribute_list[] = {
        EGL_RED_SIZE,                  8,
        EGL_GREEN_SIZE,                8,
        EGL_BLUE_SIZE,                 8,
        EGL_ALPHA_SIZE,                8,
        EGL_DEPTH_SIZE,               16,
        EGL_SURFACE_TYPE, EGL_WINDOW_BIT,
        EGL_NONE
    };
    
    static const EGLint context_attributes[] = {
        EGL_CONTEXT_CLIENT_VERSION,    2,
        EGL_NONE
    };
    
    EGLConfig config;
    EGLint    config_count;

    EGLBoolean res;
    int32_t err;


    PIGLET_PRINT("broadcom host init");
    bcm_host_init();

    
    display = eglGetDisplay(EGL_DEFAULT_DISPLAY);
    PIGLET_CHECK("eglGetDisplay");
    if (display == EGL_NO_DISPLAY) {
        PIGLET_ERROR("fail to get display!!");
        return -1;
    }


    res = eglInitialize(display, NULL, NULL);
    PIGLET_CHECK("eglInitialize");
    if (res == EGL_FALSE) {
        PIGLET_ERROR("fail to initialize eGL!!");
        return -1;
    }
    
    res = eglChooseConfig(display, attribute_list, &config, 1, &config_count);
    PIGLET_CHECK("eglChooseConfig");
    if (res == EGL_FALSE) {
        PIGLET_ERROR("fail to choose config!!");
        return -1;
    }
    
    
    res = eglBindAPI(EGL_OPENGL_ES_API);
    PIGLET_CHECK("eglBindAPI");
    if (res == EGL_FALSE) {
        PIGLET_ERROR("fail to bind API!!");
        return -1;
    }
    
    context = eglCreateContext(display, config, EGL_NO_CONTEXT, context_attributes);
    PIGLET_CHECK("eglCreateContext");
    if (context == EGL_NO_CONTEXT) {
        PIGLET_ERROR("fail to create context!!");
        return -1;
    }
        
    err = graphics_get_display_size(0, &width, &height);
    if ( err < 0 ) {
        PIGLET_ERROR("fail to get display size!!");
        return -1;
    }
    
    PIGLET_PRINT("display %dx%d",width,height);
    
    dst_rect.x = 0;
    dst_rect.y = 0;
    dst_rect.width =  width;
    dst_rect.height = height;
    
    src_rect.x = 0;
    src_rect.y = 0;
    src_rect.width =  width  << 16;
    src_rect.height = height << 16;
    
    
    dispman_display = vc_dispmanx_display_open( 0 );
    dispman_update = vc_dispmanx_update_start( 0 );
    
    dispman_element = vc_dispmanx_element_add( 
        dispman_update, 
        dispman_display, 
        0, &dst_rect,  //layer
        0, &src_rect,  //src
        DISPMANX_PROTECTION_NONE,
        0, //alpha
        0, //clamp
        0  //transform
    );

    vc_dispmanx_update_submit_sync( dispman_update );
    
    native_window.width =   width;
    native_window.height =  height;
    native_window.element = dispman_element;

    surface = eglCreateWindowSurface( display, config, &native_window, NULL );
    PIGLET_CHECK("eglCreateWindowSurface");
    if (surface == EGL_NO_SURFACE) {
        PIGLET_ERROR("fail to create window surface!!");
        return -1;
    }
        
    res = eglMakeCurrent(display, surface, surface, context);
    PIGLET_CHECK("eglMakeCurrent");
    if ( res == EGL_FALSE ) {
        PIGLET_ERROR("fail to make current!!");
        return -1;
    }
    
    PIGLET_PRINT("renderer %s %s",glGetString(GL_VENDOR),glGetString(GL_RENDERER));
    PIGLET_PRINT("version %s %s",glGetString(GL_VERSION),glGetString(GL_SHADING_LANGUAGE_VERSION));


//    PIGLET_PRINT("init: display %d at %p surface %d at %p context %d at %p",display,&display,surface,&surface,context,&context);

    return 0;    
    
}


void
MakeCurrent() 
{
//    PIGLET_PRINT("make: display %d at %p surface %d at %p context %d at %p",display,&display,surface,&surface,context,&context);
    eglMakeCurrent(display, surface, surface, context);
    PIGLET_CHECK("eglMakeCurrent");
}


void
SwapBuffers()
{
//    PIGLET_PRINT("swap: display %d at %p surface %d at %p context %d at %p",display,&display,surface,&surface,context,&context);
    eglSwapBuffers(display, surface);
    PIGLET_CHECK("eglSwapBuffers");
}


void*
GetProcAddress(const char *name)
{
    void *handle = dlopen(NULL, RTLD_LAZY);
    if (handle == NULL) {
        PIGLET_ERROR("fail to dlopen!!");
        return 0x0;
    }
    void *ret = dlsym(handle,name);
//    PIGLET_PRINT("GetProcAddress %s at %p",name,ret);
    dlclose(handle);
    return ret;
}


int
DestroyContext()
{
    EGLBoolean res;
    int32_t err;
    static DISPMANX_UPDATE_HANDLE_T  dispman_update;


    eglDestroySurface(display, surface);
    PIGLET_CHECK("eglDestroySurface");

    dispman_update = vc_dispmanx_update_start(0);
    err = vc_dispmanx_element_remove(dispman_update, dispman_element);
    if ( err != 0 ) {
        PIGLET_ERROR("fail to remove element from display!!");
        return -1;
    }

    vc_dispmanx_update_submit_sync(dispman_update);
    err = vc_dispmanx_display_close(dispman_display);
    if ( err != 0 ) {
        PIGLET_ERROR("fail to close display!!");
        return -1;
    }


    res = eglMakeCurrent(display, EGL_NO_SURFACE, EGL_NO_SURFACE, EGL_NO_CONTEXT);
    PIGLET_CHECK("eglMakeCurrent");
    if ( res == EGL_FALSE ) {
        PIGLET_ERROR("fail to make current!!");
        return -1;
    }

    eglDestroyContext(display, context);
    PIGLET_CHECK("eglMakeCurrent");

    eglTerminate(display);
    PIGLET_CHECK("eglTerminate");

    PIGLET_PRINT("terminated.");
    return 0;

}



