package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"runtime"
	"unsafe"
)

var (
	major, minor  int
	width, height int
	title         string
	core          bool
	forward       bool
)

func init() {
	runtime.LockOSThread()

	flag.IntVar(&major, "major", 3, "Major Version")
	flag.IntVar(&minor, "minor", 3, "Minor Version")
	flag.IntVar(&width, "w", 640, "Window Width")
	flag.IntVar(&height, "h", 480, "Window Height")
	flag.BoolVar(&core, "core", true, "Core Profile")
	flag.BoolVar(&forward, "forward", true, "Forward Compatible")
	flag.StringVar(&title, "title", "OpenGL Demo", "Widnow Title")
	flag.Parse()
}

func createVbo() (buffer uint32) {
	points := []float32{
		0.0, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
	}
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, unsafe.Pointer(&points[0]), gl.STATIC_DRAW)

	return
}

func createVao() (array uint32) {
	gl.GenVertexArrays(1, &array)
	gl.BindVertexArray(array)
	var index uint32 = 0
	gl.EnableVertexArrayAttrib(array, index)
	gl.VertexAttribPointer(index, 3, gl.FLOAT, false, 0, nil)

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
	/*
		glfw.SetErrorCallback(func(err glfw.ErrorCode, desc string) {
			fmt.Printf("[error] %v: %v\n", err, desc)
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
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}

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

	gl.UseProgram(program)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
