package main

import (
	"../common"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"math"
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

	program.Use()

	matrix := [16]float32{
		1.0, 0.0, 0.0, 0.0, // first column
		0.0, 1.0, 0.0, 0.0, // second column
		0.0, 0.0, 1.0, 0.0, // third column
		0.5, 0.0, 0.0, 1.0, // fourth column
	}

	matLoc := program.GetUniformLocation("matrix")
	matLoc.UniformMatrix4fv(false, matrix)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)

	speed := 1.0
	lastPos := 0.0
	prevSecs := glfw.GetTime()

	for !window.ShouldClose() {
		curSecs := glfw.GetTime()
		elapsedSecs := curSecs - prevSecs
		prevSecs = curSecs

		common.ShowFPS(window)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(0.6, 0.6, 0.8, 1.0)

		lastPos = elapsedSecs*speed + lastPos
		matrix[12] = float32(lastPos)
		if math.Abs(lastPos) > 1.0 {
			speed = -speed
		}
		matLoc.UniformMatrix4fv(false, matrix)

		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		glfw.PollEvents()
		window.SwapBuffers()

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}
