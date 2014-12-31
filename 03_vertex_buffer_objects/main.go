package main

import (
	"../common"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"os"
)

func createVbo() (buffers []gl.Buffer) {
	points := []gl.GLfloat{
		0.0, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
	}
	colours := []gl.GLfloat{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
	}

	buffers = make([]gl.Buffer, 2)
	gl.GenBuffers(buffers)
	buffers[0].Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, points, gl.STATIC_DRAW)

	buffers[1].Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(colours)*4, colours, gl.STATIC_DRAW)

	return
}

func createVao(buffers []gl.Buffer) gl.VertexArray {
	vao := gl.GenVertexArray()
	vao.Bind()

	var attrLoc gl.AttribLocation = 0
	buffers[0].Bind(gl.ARRAY_BUFFER)
	attrLoc.AttribPointer(3, gl.FLOAT, false, 0, nil)
	attrLoc.EnableArray()

	var attrLoc2 gl.AttribLocation = 1
	buffers[1].Bind(gl.ARRAY_BUFFER)
	attrLoc2.AttribPointer(3, gl.FLOAT, false, 0, nil)
	attrLoc2.EnableArray()

	return vao
}

func main() {
	window, err := common.StartGL("02_shaders")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer glfw.Terminate()
	defer window.Destroy()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	buffers := createVbo()
	defer buffers[0].Delete()
	defer buffers[1].Delete()

	vao := createVao(buffers)
	defer vao.Delete()

	vs, err := common.CreateShaderFile(gl.VERTEX_SHADER, "vs.glsl")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer vs.Delete()

	fs, err := common.CreateShaderFile(gl.FRAGMENT_SHADER, "fs.glsl")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer fs.Delete()

	program, err := common.CreateProgram(vs, fs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer program.Delete()

	common.PrintAll(program)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)

	program.Use()

	for !window.ShouldClose() {
		common.ShowFPS(window)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(0.6, 0.6, 0.8, 1.0)

		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		glfw.PollEvents()
		window.SwapBuffers()

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}
