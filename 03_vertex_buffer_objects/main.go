package main

import (
	"fmt"
	"github.com/ginuerzh/anton-gocode/common"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"os"
)

func createVbo() (buffers []uint32) {
	points := []float32{
		0.0, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
	}
	colours := []float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
	}

	buffers = make([]uint32, 2)
	gl.GenBuffers(2, &buffers[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, gl.Ptr(points), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[1])
	gl.BufferData(gl.ARRAY_BUFFER, len(colours)*4, gl.Ptr(colours), gl.STATIC_DRAW)

	return
}

func createVao(buffers []uint32) (vao uint32) {
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[0])
	var attrLoc1 uint32 = 0
	gl.EnableVertexAttribArray(attrLoc1)
	gl.VertexAttribPointer(attrLoc1, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	var attrLoc2 uint32 = 1
	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[1])
	gl.EnableVertexAttribArray(attrLoc2)
	gl.VertexAttribPointer(attrLoc2, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	return
}

func main() {
	window, err := common.StartGL("03 - Vertex Buffer Objects")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer glfw.Terminate()
	defer window.Destroy()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	buffers := createVbo()
	defer gl.DeleteBuffers(2, &buffers[0])

	vao := createVao(buffers)
	defer gl.DeleteVertexArrays(1, &vao)

	vs, err := common.CreateShaderFile(gl.VERTEX_SHADER, "vs.glsl")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer gl.DeleteShader(vs)

	fs, err := common.CreateShaderFile(gl.FRAGMENT_SHADER, "fs.glsl")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer gl.DeleteShader(fs)

	program, err := common.CreateProgram(vs, fs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer gl.DeleteProgram(program)

	common.PrintAll(program)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)

	gl.UseProgram(program)

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
