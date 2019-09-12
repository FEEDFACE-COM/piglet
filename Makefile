
BUILD_NAME      = piglet
BUILD_VERSION  ?= $(shell git describe --tags)
BUILD_PLATFORM ?= $(shell go env GOOS )-$(shell go env GOARCH)

GLOW ?= ${GOPATH}/bin/glow



OPENGL_API      ?= gles2
OPENGL_VERSION  ?= 2.0
OPENGL_REMEXT   ?= GL_EXT_win32_|GL_APPLE_|GL_QCOM_|GL_INTEL_|GL_NVX_|GL_NV_|GL_AMD_|GL_KHR_|GL_EXT_|
OPENGL_ADDEXT   ?= GL_EXT_discard_framebuffer


OPENGL_FILES = $(addprefix ${OPENGL_API}/, conversions.go package.go procaddr.go error_string.go)
PIGLET_FILES = $(wildcard *.go *.c *.h)

GLOW_FILES = $(addprefix ${OPENGL_API}/, conversions.go  package.go )




help:
	@echo "### Usage ###"
	@echo " make ${OPENGL_API}    # build bindings"
	@echo " make glow     # fetch glow tool"
	@echo " make specs    # fetch opengl specs"
	@echo " make info     # show build info"
	@echo " make clean    # clean up"


info: 
	@echo "### Version Info ###"
	@echo " name       ${BUILD_NAME}"
	@echo " version    ${BUILD_VERSION}"
	@echo " platform   ${BUILD_PLATFORM}"
	@echo " opengl     ${OPENGL_API} ${OPENGL_VERSION}"
	@echo "### Package Info ###"
	@echo " piglet     ${PIGLET_FILES}"
	@echo " opengl     ${OPENGL_FILES}"
	


${OPENGL_API}: ${OPENGL_FILES}



PACKAGE_CFLAGS=// \#cgo linux,arm  CFLAGS: -I/opt/vc/include
PACKAGE_LDFLAGS=// \#cgo linux,arm LDFLAGS: -L/opt/vc/lib


${OPENGL_API}/package.go: tmp/package.go
# strip all cgo directives, and add our own
	sed -e '/package gles2/,\%#include <KHR/khrplatform.h>% { s|// #cgo.*|//|; s|^$$|\n${PACKAGE_CFLAGS}\n${PACKAGE_LDFLAGS}|; }' $^ >| $@


${OPENGL_API}/conversions.go: tmp/conversions.go
	cp -f $^ $@

	
tmp/%.go:
	mkdir -p tmp
	${GLOW} generate -out tmp -api=${OPENGL_API} -version=${OPENGL_VERSION} -remext="${OPENGL_REMEXT}" -addext="${OPENGL_ADDEXT}"


specs:
	${GLOW} download

glow:
	go get -v github.com/go-gl/glow


	
clean:
	go clean -v
	rm -rf ./tmp
	rm -f ${GLOW_FILES}





.PHONY: help info ${OPENGL_API} specs glow clean

