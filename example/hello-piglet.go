// +build linux,arm

package main

import (
	"encoding/base64"
	"fmt"
	"github.com/FEEDFACE-COM/piglet"
	gl "github.com/FEEDFACE-COM/piglet/gles2"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/draw"
	_ "image/png"
	"math"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"
)

/* draw ***********************************************************************/

func UpdateScene(shaderProgram uint32, camera mgl32.Mat4, startTime time.Time) {

	clock := float32(time.Now().Sub(startTime).Seconds())

	angle := PI / 4. * Sin(clock)
	cam := camera.Mul4(mgl32.HomogRotate3DY(angle))

	gl.UseProgram(shaderProgram)
	ptr := gl.GetUniformLocation(shaderProgram, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(ptr, 1, false, &cam[0])

}

func DrawScene(shaderProgram, vertexBuffer, textureName uint32) {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(shaderProgram)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BindTexture(gl.TEXTURE_2D, textureName)

	gl.DrawArrays(gl.TRIANGLES, 0, 2*3)

	err := gl.GetError()
	if err != gl.NO_ERROR {
		Error("fail draw: %s", gl.ErrorString(err))
	}

}

/* init ***********************************************************************/

func InitScene(width, height int32) (uint32, uint32, uint32, mgl32.Mat4, time.Time) {
	Notice("init scene..")

	cameraMatrix := InitCameraMatrix(float32(width), float32(height))
	vertexBuffer := InitVertexBuffer()
	textureName := InitTextureName()

	vertexShader := InitShader(VertexSource, gl.VERTEX_SHADER)
	fragmentShader := InitShader(FragmentSource, gl.FRAGMENT_SHADER)
	shaderProgram := InitShaderProgram(vertexShader, fragmentShader)

	gl.Viewport(0, 0, width, height)
	gl.ClearColor(0., 0., 0., 0.)

	if gl.GetError() != gl.NO_ERROR {
		Error("fail init scene")
	}

	startTime := time.Now()

	return shaderProgram, vertexBuffer, textureName, cameraMatrix, startTime
}

func InitCameraMatrix(width, height float32) mgl32.Mat4 {

	fov, ratio := mgl32.DegToRad(45.), width/height
	near, far := float32(0.01), float32(10.0)
	projection := mgl32.Perspective(fov, ratio, near, far)

	pos := mgl32.Vec3{0, 0, 5}
	lookat := mgl32.Vec3{0, 0, 0}
	up := mgl32.Vec3{0, 1, 0}
	view := mgl32.LookAtV(pos, lookat, up)

	var cameraMatrix = mgl32.Ident4()
	cameraMatrix = cameraMatrix.Mul4(projection)
	cameraMatrix = cameraMatrix.Mul4(view)
	return cameraMatrix
}

func InitVertexBuffer() uint32 {
	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(VertexData)*4, gl.Ptr(VertexData), gl.STATIC_DRAW)

	if gl.GetError() != gl.NO_ERROR {
		Error("fail init buffer")
	}

	return vertexBuffer
}

func InitTextureName() uint32 {
	var textureName uint32

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(TextureData))
	img, _, _ := image.Decode(reader)
	textureRGBA := image.NewRGBA(img.Bounds())
	draw.Draw(textureRGBA, textureRGBA.Bounds(), img, image.Point{0, 0}, draw.Src)
	w, h := int32(textureRGBA.Rect.Size().X), int32(textureRGBA.Rect.Size().Y)
	ptr := gl.Ptr(textureRGBA.Pix)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.GenTextures(1, &textureName)
	gl.BindTexture(gl.TEXTURE_2D, textureName)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, ptr)

	if gl.GetError() != gl.NO_ERROR {
		Error("fail init texture")
	}
	return textureName
}

func InitShader(source string, xtype uint32) uint32 {
	var status int32
	shader := gl.CreateShader(xtype)
	src, free := gl.Strs(source + "\x00")
	defer free()
	gl.ShaderSource(shader, 1, src, nil)
	gl.CompileShader(shader)
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		Error("fail compile shader")
	}
	return shader
}

func InitShaderProgram(vertexShader, fragmentShader uint32) uint32 {
	var status int32
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		Error("fail link program")
	}

	ptr := gl.GetAttribLocation(shaderProgram, gl.Str("position\x00"))
	gl.EnableVertexAttribArray(uint32(ptr))
	gl.VertexAttribPointer(uint32(ptr), 3, gl.FLOAT, false, (3+2+4)*4, gl.PtrOffset(0*4))

	ptr = gl.GetAttribLocation(shaderProgram, gl.Str("texCoord\x00"))
	gl.EnableVertexAttribArray(uint32(ptr))
	gl.VertexAttribPointer(uint32(ptr), 2, gl.FLOAT, false, (3+2+4)*4, gl.PtrOffset(3*4))

	ptr = gl.GetAttribLocation(shaderProgram, gl.Str("color\x00"))
	gl.EnableVertexAttribArray(uint32(ptr))
	gl.VertexAttribPointer(uint32(ptr), 4, gl.FLOAT, false, (3+2+4)*4, gl.PtrOffset((3+2)*4))

	return shaderProgram
}

/* context ********************************************************************/

func ConfigureContext() (int32, int32) {

	Notice("create context..")
	err := piglet.CreateContext()
	if err != nil {
		Error("fail create context: %s", err)
	}
	width, height := piglet.GetDisplaySize()
	Notice("display: %dx%d", width, height)

	piglet.MakeCurrent()
	gl.InitWithProcAddrFunc(piglet.GetProcAddress)
	Notice("renderer: %s %s", gl.GoStr(gl.GetString((gl.VENDOR))), gl.GoStr(gl.GetString((gl.RENDERER))))
	Notice("version: %s / %s", gl.GoStr(gl.GetString((gl.VERSION))), gl.GoStr(gl.GetString((gl.SHADING_LANGUAGE_VERSION))))

	return width, height
}

func TerminateContext() {

	Notice("destroy context..")

	err := piglet.DestroyContext()
	if err != nil {
		Error("fail create context: %s", err)
	}

}

/* main ***********************************************************************/

const FrameRate = 60.

func main() {
	Notice("Hello, PiGLEt!!")

	tickChan := make(chan bool, 1)

	// handle user interrupt
	RegisterSignalHandler(tickChan)

	// lock main goroutine to current thread!
	runtime.LockOSThread()

	// configure opengl context
	width, height := ConfigureContext()

	// init scene
	program, buffer, texture, camera, startTime := InitScene(width, height)

	// start ticker
	go Ticker(tickChan)

	Notice("start draw..")
	ever := true
	for ever {

		UpdateScene(program, camera, startTime)
		DrawScene(program, buffer, texture)
		piglet.SwapBuffers()

		ever = WaitForTick(tickChan)

	}

	TerminateContext()

	os.Exit(0)
}

/* tick ***********************************************************************/

func Ticker(tickChan chan bool) {

	for { //ever
		tickChan <- true // send tick
		time.Sleep(time.Duration(int64(time.Second / FrameRate)))
	}

}

func WaitForTick(tickChan chan bool) bool {
	var tick bool
	tick = <-tickChan // read tick
	if !tick {
		return false
	}
	for { // clear all ticks
		select {
		case tick = <-tickChan:
			if !tick {
				return false
			}
		default:
			return true
		}
	}

}

/* math ***********************************************************************/

const PI = 3.1415926535897932384626433832795028841971693993751058209749445920

func Sin(x float32) float32 { return float32(math.Sin(float64(x))) }

/* util ***********************************************************************/

func RegisterSignalHandler(tickChan chan bool) {
	breakChan := make(chan os.Signal, 1)
	signal.Notify(breakChan, os.Interrupt)
	go func() {
		<-breakChan
		Notice("\rbreak.")
		tickChan <- false
	}()
}

func Notice(format string, args ...interface{}) { fmt.Fprintf(os.Stderr, format+"\n", args...) }
func Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(-1)
}

/* data ***********************************************************************/

var VertexData = []float32{

	// position       texCoord color
	// x,    y,   z,  s,  t,   r,   g,   b,  a,
	-1.0, -1.0, +0., 0., 1., 0.2, 0.2, 0.2, 1., // A    D     C
	+1.0, -1.0, +0., 1., 1., 0.2, 0.2, 0.2, 1., // B     +---+
	+1.0, +1.0, +0., 1., 0., 0.8, 0.8, 0.8, 1., // C     |  /|
	//                                                   | / |
	-1.0, -1.0, +0., 0., 1., 0.2, 0.2, 0.2, 1., // A     |/  |
	+1.0, +1.0, +0., 1., 0., 0.8, 0.8, 0.8, 1., // C     +---+
	-1.0, +1.0, +0., 0., 0., 0.8, 0.8, 0.8, 1., // D    A     B

}

const VertexSource = `
uniform mat4 camera;
attribute vec3 position;
attribute vec2 texCoord; varying vec2 vTexCoord;
attribute vec4 color;    varying vec4 vColor;


void main() {
    vTexCoord = texCoord;
    vColor = color; 
    gl_Position = camera * vec4(position, 1);
}
`

const FragmentSource = `
uniform sampler2D texture;
varying vec2 vTexCoord;
varying vec4 vColor;

void main() {
    vec4 texColor = texture2D(texture,vTexCoord);
    gl_FragColor = vec4(vColor.rgb, texColor.a);
}
`

const TextureData = "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"
