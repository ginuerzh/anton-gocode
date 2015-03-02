package main

import (
	"../common"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"os"
	"runtime"
)

func init() {
	runtime.LockOSThread()
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

func createVao() (vao uint32) {
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	var index uint32 = 0
	gl.EnableVertexAttribArray(index)
	gl.VertexAttribPointer(index, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	return
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
	defer gl.DeleteBuffers(1, &buffer)

	vao := createVao()
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

	name := []byte("inputColour")
	colorLoc := gl.GetUniformLocation(program, &name[0])
	if colorLoc < 0 {
		fmt.Fprintf(os.Stderr, "Can't find uniform %s, location %d\n",
			name, colorLoc)
		return
	}
	gl.UseProgram(program)
	gl.Uniform4f(colorLoc, 1.0, 0.0, 0.0, 1.0)

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
