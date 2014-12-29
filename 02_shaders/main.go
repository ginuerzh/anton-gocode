package main

import (
	"../common"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"os"
)

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
	window, err := common.StartGL("02_shaders")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer glfw.Terminate()
	defer window.Destroy()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	buffer := createVbo()
	defer buffer.Delete()

	vao := createVao()
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

	colorLoc := program.GetUniformLocation("inputColour")
	if colorLoc < 0 {
		fmt.Fprintf(os.Stderr, "Can't find uniform %s, location %d\n",
			"inputColour", colorLoc)
		return
	}
	program.Use()
	colorLoc.Uniform4f(1.0, 0.0, 0.0, 1.0)

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
