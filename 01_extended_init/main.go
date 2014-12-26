package main

import (
	"../common"
	"bytes"
	"flag"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"os"
	"runtime"
	"time"
)

var (
	major, minor  int
	width, height int
	title         string
	core          bool
	forward       bool
)

func init() {
	flag.IntVar(&major, "major", 3, "Major Version")
	flag.IntVar(&minor, "minor", 3, "Minor Version")
	flag.IntVar(&width, "w", 640, "Window Width")
	flag.IntVar(&height, "h", 480, "Window Height")
	flag.BoolVar(&core, "core", true, "Core Profile")
	flag.BoolVar(&forward, "forward", true, "Forward Compatible")
	flag.StringVar(&title, "title", "OpenGL Demo", "Widnow Title")
	flag.Parse()
}

const (
	glLogFile = "gl.log"
)

func restartGLLog() error {
	file, err := os.Create(glLogFile)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: could not open GL_LOG_FILE log file %s for writing\n",
			glLogFile)
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "GL_LOG_FILE log. local time %s\n", time.Now().String())
	fmt.Fprintf(file, "build version: %s\n\n", runtime.Version())

	return nil
}

func glLog(message string, a ...interface{}) error {
	file, err := os.OpenFile(glLogFile, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: could not open GL_LOG_FILE %s file for appending\n",
			glLogFile)
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, message, a...)

	return nil
}

func logGLParams() {
	params := []gl.GLenum{
		gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS,
		gl.MAX_CUBE_MAP_TEXTURE_SIZE,
		gl.MAX_DRAW_BUFFERS,
		gl.MAX_FRAGMENT_UNIFORM_COMPONENTS,
		gl.MAX_TEXTURE_IMAGE_UNITS,
		gl.MAX_TEXTURE_SIZE,
		gl.MAX_VARYING_FLOATS,
		gl.MAX_VERTEX_ATTRIBS,
		gl.MAX_VERTEX_TEXTURE_IMAGE_UNITS,
		gl.MAX_VERTEX_UNIFORM_COMPONENTS,
		gl.MAX_VIEWPORT_DIMS,
		gl.STEREO,
	}
	names := []string{
		"GL_MAX_COMBINED_TEXTURE_IMAGE_UNITS",
		"GL_MAX_CUBE_MAP_TEXTURE_SIZE",
		"GL_MAX_DRAW_BUFFERS",
		"GL_MAX_FRAGMENT_UNIFORM_COMPONENTS",
		"GL_MAX_TEXTURE_IMAGE_UNITS",
		"GL_MAX_TEXTURE_SIZE",
		"GL_MAX_VARYING_FLOATS",
		"GL_MAX_VERTEX_ATTRIBS",
		"GL_MAX_VERTEX_TEXTURE_IMAGE_UNITS",
		"GL_MAX_VERTEX_UNIFORM_COMPONENTS",
		"GL_MAX_VIEWPORT_DIMS",
		"GL_STEREO",
	}

	glLog("GL Context Params:\n")

	p := make([]int32, 2)
	for i, v := range params {
		if v == gl.STEREO {
			p := make([]bool, 1)
			gl.GetBooleanv(v, p)
			glLog("%s %v\n", names[i], p[0])
			continue
		}

		p[0] = 0
		p[1] = 0
		gl.GetIntegerv(v, p)
		if v == gl.MAX_VIEWPORT_DIMS {
			glLog("%s %d %d\n", names[i], p[0], p[1])
			continue
		}
		glLog("%s %d\n", names[i], p[0])
	}

	glLog("-----------------------------\n")
}

var prevSecs float64
var frameCount int

func updateFPSCounter(window *glfw.Window) {
	curSecs := glfw.GetTime()
	elapsedSecs := curSecs - prevSecs
	if elapsedSecs > 0.25 {
		prevSecs = curSecs
		fps := float64(frameCount) / elapsedSecs
		window.SetTitle(fmt.Sprintf("opengl @ fps: %.2f", fps))
		frameCount = 0
	}
	frameCount++
}

func createVbo() gl.Buffer {
	points := []gl.GLfloat{
		0.0, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
	}
	buffer := gl.GenBuffer()
	buffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, points, gl.STATIC_DRAW)

	return buffer
}

func createVao() gl.VertexArray {
	vao := gl.GenVertexArray()
	vao.Bind()
	var attrLoc gl.AttribLocation = 0
	attrLoc.EnableArray()
	attrLoc.AttribPointer(3, gl.FLOAT, false, 0, nil)

	return vao
}

func main() {
	restartGLLog()
	glLog("starting GLFW\n%s\n\n", glfw.GetVersionString())

	glfw.SetErrorCallback(func(err glfw.ErrorCode, desc string) {
		glLog("GLFW ERROR: code %d msg: %s\n", err, desc)
	})

	if !glfw.Init() {
		panic("ERROR: could not start GLFW3")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, major)
	glfw.WindowHint(glfw.ContextVersionMinor, minor)
	if forward && major >= 3 {
		glfw.WindowHint(glfw.OpenglForwardCompatible, 1)
	}
	if core && major >= 3 && minor >= 2 {
		glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	}
	glfw.WindowHint(glfw.Samples, 4)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()
	if r := gl.Init(); r > 0 {
		fmt.Println("init opengl:", r)
	}

	glLog("Vendor: %s\n", gl.GetString(gl.VENDOR))
	glLog("Renderer: %s\n", gl.GetString(gl.RENDERER))
	glLog("Version: %s\n", gl.GetString(gl.VERSION))
	glLog("Shading language version: %s\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	glLog("Extensions: %s\n\n", gl.GetString(gl.EXTENSIONS))

	logGLParams()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	buffer := createVbo()
	defer buffer.Delete()

	vao := createVao()
	defer vao.Delete()

	vertex_shader := `
	#version 330
	in vec3 vp;
	void main () {
		gl_Position = vec4 (vp, 1.0);
	}
`

	fragment_shader := `
	#version 330
	out vec4 frag_colour;
	void main () {
		frag_colour = vec4 (0.5, 0.0, 0.5, 1.0);
	}
`
	vs, err := common.CreateShader(gl.VERTEX_SHADER,
		bytes.NewReader([]byte(vertex_shader)))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer vs.Delete()

	fs, err := common.CreateShader(gl.FRAGMENT_SHADER,
		bytes.NewReader([]byte(fragment_shader)))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fs.Delete()

	program, err := common.CreateProgram(vs, fs)
	if err != nil {
		fmt.Println(err)
	}
	defer program.Delete()

	for !window.ShouldClose() {
		updateFPSCounter(window)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		program.Use()
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}