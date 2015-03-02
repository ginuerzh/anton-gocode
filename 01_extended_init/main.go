package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
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
	fullscreen    bool
)

func init() {
	flag.IntVar(&major, "major", 3, "Major Version")
	flag.IntVar(&minor, "minor", 3, "Minor Version")
	flag.IntVar(&width, "w", 640, "Window Width")
	flag.IntVar(&height, "h", 480, "Window Height")
	flag.BoolVar(&fullscreen, "full", false, "Fullscreen")
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

func GLog(message string, a ...interface{}) error {
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
	params := []uint32{
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

	GLog("GL Context Params:\n")

	p := make([]int32, 2)
	for i, v := range params {
		if v == gl.STEREO {
			p := false
			gl.GetBooleanv(v, &p)
			GLog("%s %v\n", names[i], p)
			continue
		}

		p[0] = 0
		p[1] = 0
		gl.GetIntegerv(v, &p[0])
		if v == gl.MAX_VIEWPORT_DIMS {
			GLog("%s %d %d\n", names[i], p[0], p[1])
			continue
		}
		GLog("%s %d\n", names[i], p[0])
	}

	GLog("-----------------------------\n")
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

/* we can run a full-screen window here */
func fullscr() (width int, height int, monitor *glfw.Monitor, err error) {
	GLog("Full Screen Mode\n")

	monitor = glfw.GetPrimaryMonitor()
	vm := monitor.GetVideoMode()

	GLog("Primary monitor: %s (%d*%d, %dHZ)\n\n",
		monitor.GetName(), vm.Width, vm.Height, vm.RefreshRate)

	return vm.Width, vm.Height, monitor, nil
}

func createVbo() (buffer uint32) {
	points := []float32{
		0.0, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
	}
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, gl.Ptr(points), gl.STATIC_DRAW)

	return
}

func createVao() (array uint32) {
	gl.GenVertexArrays(1, &array)
	gl.BindVertexArray(array)
	var index uint32 = 0
	gl.EnableVertexAttribArray(index)
	gl.VertexAttribPointer(index, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	return
}

func getShanderInfoLog(shader uint32) string {
	var length int32
	var infoLog []byte
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
	if length > 0 {
		infoLog = make([]byte, length)
		gl.GetShaderInfoLog(shader, length, nil, &infoLog[0])
	}

	return string(infoLog)
}

func CreateShader(shaderType uint32, src []byte) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	xstring := &src[0]
	gl.ShaderSource(shader, 1, &xstring, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	infoLog := getShanderInfoLog(shader)
	if len(infoLog) > 0 {
		fmt.Printf("shader info log for GL index %d:\n%s\n", shader, infoLog)
	}

	if status == gl.FALSE {
		gl.DeleteShader(shader)
		return shader, errors.New("Compile: " + infoLog)
	}

	return shader, nil
}

func getProgramInfoLog(program uint32) string {
	var length int32
	var infoLog []byte
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
	if length > 0 {
		infoLog = make([]byte, length)
		gl.GetProgramInfoLog(program, length, nil, &infoLog[0])
	}

	return string(infoLog)
}

func CreateProgram(shaders ...uint32) (uint32, error) {
	program := gl.CreateProgram()

	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}

	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)

	infoLog := getProgramInfoLog(program)
	if len(infoLog) > 0 {
		fmt.Printf("Program link info log for GL index %d:\n%s\n", program, infoLog)
	}

	if status == gl.FALSE {
		gl.DeleteProgram(program)
		return program, errors.New("Link: " + infoLog)
	}

	gl.ValidateProgram(program)
	gl.GetProgramiv(program, gl.VALIDATE_STATUS, &status)
	infoLog = getProgramInfoLog(program)
	if len(infoLog) > 0 {
		fmt.Printf("Program validate info log for GL index %d:\n%s\n", program, infoLog)
	}
	if status == gl.FALSE {
		gl.DeleteProgram(program)
		return program, errors.New("Validate: " + infoLog)
	}

	return program, nil
}

func main() {
	restartGLLog()
	GLog("starting GLFW\n%s\n\n", glfw.GetVersionString())

	/*
		glfw.SetErrorCallback(func(err glfw.ErrorCode, desc string) {
			glLog("GLFW ERROR: code %d msg: %s\n", err, desc)
		})
	*/

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, major)
	glfw.WindowHint(glfw.ContextVersionMinor, minor)
	if forward && major >= 3 {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)
	}
	if core && major >= 3 && minor >= 2 {
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	}
	glfw.WindowHint(glfw.Samples, 16)

	var monitor *glfw.Monitor
	if fullscreen {
		width, height, monitor, _ = fullscr()
	}

	window, err := glfw.CreateWindow(width, height, title, monitor, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.SetSizeCallback(func(win *glfw.Window, w, h int) {
		width = w
		height = h
		//fmt.Printf("width %d height %d\n", width, height)
	})

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	GLog("Vendor: %s\n", gl.GoStr(gl.GetString(gl.VENDOR)))
	GLog("Renderer: %s\n", gl.GoStr(gl.GetString(gl.RENDERER)))
	GLog("Version: %s\n", gl.GoStr(gl.GetString(gl.VERSION)))
	GLog("Shading language version: %s\n", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	GLog("Extensions: %s\n\n", gl.GoStr(gl.GetString(gl.EXTENSIONS)))

	logGLParams()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	buffer := createVbo()
	defer gl.DeleteBuffers(1, &buffer)

	vao := createVao()
	defer gl.DeleteVertexArrays(1, &vao)

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
	vs, err := CreateShader(gl.VERTEX_SHADER, []byte(vertex_shader))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gl.DeleteShader(vs)

	fs, err := CreateShader(gl.FRAGMENT_SHADER, []byte(fragment_shader))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gl.DeleteShader(fs)

	program, err := CreateProgram(vs, fs)
	if err != nil {
		fmt.Println(err)
	}
	defer gl.DeleteProgram(program)

	for !window.ShouldClose() {
		updateFPSCounter(window)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Viewport(0, 0, int32(width), int32(height))

		gl.UseProgram(program)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		glfw.PollEvents()
		window.SwapBuffers()

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}
