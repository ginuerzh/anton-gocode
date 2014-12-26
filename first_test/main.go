package main

import (
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

	vendor := gl.GetString(gl.VENDOR)
	version := gl.GetString(gl.VERSION)
	renderer := gl.GetString(gl.RENDERER)
	shaderVer := gl.GetString(gl.SHADING_LANGUAGE_VERSION)
	extensions := gl.GetString(gl.EXTENSIONS)

	fmt.Println("Vendor:", vendor)
	fmt.Println("Renderer:", renderer)
	fmt.Println("OpenGL version supported", version)
	fmt.Println("Shading language version string:", shaderVer)
	fmt.Println("Extensions:", extensions)
}
