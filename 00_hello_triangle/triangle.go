package main

import (
	"../common"
	"bytes"
	"flag"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
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
	glfw.SetErrorCallback(func(err glfw.ErrorCode, desc string) {
		fmt.Printf("[error] %v: %v\n", err, desc)
	})

	if !glfw.Init() {
		panic("Can't init glfw!")
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
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()
	if r := gl.Init(); r > 0 {
		fmt.Println("init opengl:", r)
	}

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
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		program.Use()
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
