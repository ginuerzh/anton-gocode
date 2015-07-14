package main

import (
	"flag"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
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

func main() {
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

	if err = gl.Init(); err != nil {
		panic(err)
	}

	fmt.Println("Vendor:", gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println("Renderer:", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println("OpenGL version supported", gl.GoStr(gl.GetString(gl.RENDERER)))
	fmt.Println("Shading language version string:", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	fmt.Println("Extensions:", gl.GoStr(gl.GetString(gl.EXTENSIONS)))
}
