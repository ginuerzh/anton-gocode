package main

import (
	"fmt"
	"github.com/ginuerzh/anton-gocode/common"
	"github.com/ginuerzh/math3d/m32"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	//"math"
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
	window, err := common.StartGL("05 - Virtual Camera")
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

	/* create PROJECTION MATRIX */
	w, h := common.WindowSize()
	aspect := float32(w) / float32(h)
	pm := m32.Perspective(67.0, aspect, 0.1, 100.0)

	/* create VIEW MATRIX */
	speed := 1.0               // 1 unit per second
	yawSpeed := 10.0           // 10 degrees per second
	yawXDeg := float32(0.0)    // x-rotation in degrees
	yawYDeg := float32(0.0)    // y-rotation in degrees
	pos := m32.Vec3{0, 0, 2.0} // don't start at zero, or we will be too close
	t := m32.Ident4().Translate(pos.Negate())
	r := m32.Ident4().RotateY(-yawYDeg)
	r = r.RotateX(-yawXDeg)
	vm := r.Mul4(t)

	viewMatLoc := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projMatLoc := gl.GetUniformLocation(program, gl.Str("proj\x00"))

	gl.UseProgram(program)
	gl.UniformMatrix4fv(viewMatLoc, 1, false, &vm[0])
	gl.UniformMatrix4fv(projMatLoc, 1, false, &pm[0])

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)

	prevSecs := glfw.GetTime()

	for !window.ShouldClose() {
		curSecs := glfw.GetTime()
		elapsedSecs := curSecs - prevSecs
		prevSecs = curSecs

		common.ShowFPS(window)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(0.6, 0.6, 0.8, 1.0)

		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		glfw.PollEvents()

		moved := false

		if window.GetKey(glfw.KeyA) != glfw.Release {
			pos[0] -= float32(speed * elapsedSecs)
			moved = true
		}
		if window.GetKey(glfw.KeyD) != glfw.Release {
			pos[0] += float32(speed * elapsedSecs)
			moved = true
		}
		if window.GetKey(glfw.KeyLeftShift) == glfw.Release &&
			window.GetKey(glfw.KeySpace) != glfw.Release {
			pos[1] += float32(speed * elapsedSecs)
			moved = true
		}
		if window.GetKey(glfw.KeyLeftShift) != glfw.Release &&
			window.GetKey(glfw.KeySpace) != glfw.Release {
			pos[1] -= float32(speed * elapsedSecs)
			moved = true
		}
		if window.GetKey(glfw.KeyW) != glfw.Release {
			pos[2] -= float32(speed * elapsedSecs)
			moved = true
		}
		if window.GetKey(glfw.KeyS) != glfw.Release {
			pos[2] += float32(speed * elapsedSecs)
			moved = true
		}
		if window.GetKey(glfw.KeyLeft) != glfw.Release {
			yawYDeg += float32(yawSpeed * elapsedSecs)
			moved = true
		}
		if window.GetKey(glfw.KeyRight) != glfw.Release {
			yawYDeg -= float32(yawSpeed * elapsedSecs)
			moved = true
		}
		if window.GetKey(glfw.KeyUp) != glfw.Release {
			yawXDeg += float32(yawSpeed * elapsedSecs)
			moved = true
		}
		if window.GetKey(glfw.KeyDown) != glfw.Release {
			yawXDeg -= float32(yawSpeed * elapsedSecs)
			moved = true
		}
		/* update view matrix */
		if moved {
			t = m32.Ident4().Translate(pos.Negate())
			r = m32.Ident4().RotateY(-yawYDeg)
			r = r.RotateX(-yawXDeg)
			vm = r.Mul4(t)
			gl.UniformMatrix4fv(viewMatLoc, 1, false, &vm[0])
		}

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		window.SwapBuffers()
	}
}
