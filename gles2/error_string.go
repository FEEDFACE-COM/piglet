
package gles2

func ErrorString(error uint32) string {
    switch (error) {
        case NO_ERROR:                      return "NO_ERROR"
        case INVALID_ENUM:                  return "INVALID_ENUM"
        case INVALID_VALUE:                 return "INVALID_VALUE"
        case INVALID_OPERATION:             return "INVALID_OPERATION"
        case INVALID_FRAMEBUFFER_OPERATION: return "INVALID_FRAMEBUFFER_OPERATION"
    }
    return "UNKNOWN"
}

