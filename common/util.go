package common

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"io/ioutil"
	"os"
	"runtime"
	"time"
)

const (
	glLogFile = "gl.log"
)

type Config struct {
	Major, Minor  int
	Width, Height int
	Title         string
	Core          bool
	Forward       bool
	Fullscreen    bool
	FPS           bool
	Log           bool
}

var (
	config Config
)

func init() {
	flag.IntVar(&config.Major, "major", 3, "Major Version")
	flag.IntVar(&config.Minor, "minor", 3, "Minor Version")
	flag.IntVar(&config.Width, "w", 640, "Window Width")
	flag.IntVar(&config.Height, "h", 480, "Window Height")
	flag.BoolVar(&config.Fullscreen, "full", false, "Fullscreen")
	flag.BoolVar(&config.FPS, "fps", true, "Show FPS")
	flag.BoolVar(&config.Core, "core", true, "Core Profile")
	flag.BoolVar(&config.Forward, "forward", true, "Forward Compatible")
	flag.BoolVar(&config.Log, "log", true, "Enable log")
	flag.Parse()
}

func GetConfig() Config {
	return config
}

func StartGL(title string) (window *glfw.Window, err error) {
	restartGLLog()

	GLog("starting GLFW\n%s\n\n", glfw.GetVersionString())

	glfw.SetErrorCallback(func(err glfw.ErrorCode, desc string) {
		GLogErr("GLFW ERROR: code %d msg: %s\n", err, desc)
	})

	if !glfw.Init() {
		GLogErr("ERROR: could not init GLFW3\n")
		return nil, errors.New("ERROR: could not init GLFW3")
	}

	glfw.WindowHint(glfw.ContextVersionMajor, config.Major)
	glfw.WindowHint(glfw.ContextVersionMinor, config.Minor)
	if config.Forward && config.Major >= 3 {
		glfw.WindowHint(glfw.OpenglForwardCompatible, 1)
	}
	if config.Core && config.Major >= 3 && config.Minor >= 2 {
		glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	}
	glfw.WindowHint(glfw.Samples, 16)

	var monitor *glfw.Monitor
	if config.Fullscreen {
		config.Width, config.Height, monitor, _ = fullscreen()
	}

	window, err = glfw.CreateWindow(config.Width, config.Height,
		title, monitor, nil)
	if err != nil {
		glfw.Terminate()
		return
	}
	config.Title = title

	window.SetSizeCallback(func(win *glfw.Window, w, h int) {
		config.Width = w
		config.Height = h
		//fmt.Printf("width %d height %d\n", width, height)
		gl.Viewport(0, 0, config.Width, config.Height)
	})

	window.MakeContextCurrent()

	if r := gl.Init(); r > 0 {
		GLogErr("ERROR: could not init OpenGL\n")
		err = errors.New("ERROR: could not init OpenGL")
		glfw.Terminate()
		return
	}

	GLog("Vendor: %s\n", gl.GetString(gl.VENDOR))
	GLog("Renderer: %s\n", gl.GetString(gl.RENDERER))
	GLog("Version: %s\n", gl.GetString(gl.VERSION))
	GLog("Shading language version: %s\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	//GLog("Extensions: %s\n\n", gl.GetString(gl.EXTENSIONS))

	logGLParams()

	return
}

func restartGLLog() error {
	if !config.Log {
		return nil
	}

	file, err := os.Create(glLogFile)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: could not open GL_LOG_FILE log file %s for writing\n",
			glLogFile)
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "GL_LOG_FILE log. local time %s\n", time.Now().String())
	fmt.Fprintf(file, "build version: %s\n\n", runtime.Version())

	return nil
}

func GLog(message string, a ...interface{}) error {
	if !config.Log {
		return nil
	}

	file, err := os.OpenFile(glLogFile, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: could not open GL_LOG_FILE %s file for appending\n",
			glLogFile)
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, message, a...)

	return nil
}

func GLogErr(message string, a ...interface{}) error {
	fmt.Fprintf(os.Stderr, message, a...)
	return GLog(message, a...)
}

func logGLParams() {
	params := []gl.GLenum{
		gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS,
		gl.MAX_CUBE_MAP_TEXTURE_SIZE,
		gl.MAX_DRAW_BUFFERS,
		gl.MAX_FRAGMENT_UNIFORM_COMPONENTS,
		gl.MAX_TEXTURE_IMAGE_UNITS,
		gl.MAX_TEXTURE_SIZE,
		gl.MAX_VARYING_FLOATS,
		gl.MAX_VERTEX_ATTRIBS,
		gl.MAX_VERTEX_TEXTURE_IMAGE_UNITS,
		gl.MAX_VERTEX_UNIFORM_COMPONENTS,
		gl.MAX_VIEWPORT_DIMS,
		gl.STEREO,
	}
	names := []string{
		"GL_MAX_COMBINED_TEXTURE_IMAGE_UNITS",
		"GL_MAX_CUBE_MAP_TEXTURE_SIZE",
		"GL_MAX_DRAW_BUFFERS",
		"GL_MAX_FRAGMENT_UNIFORM_COMPONENTS",
		"GL_MAX_TEXTURE_IMAGE_UNITS",
		"GL_MAX_TEXTURE_SIZE",
		"GL_MAX_VARYING_FLOATS",
		"GL_MAX_VERTEX_ATTRIBS",
		"GL_MAX_VERTEX_TEXTURE_IMAGE_UNITS",
		"GL_MAX_VERTEX_UNIFORM_COMPONENTS",
		"GL_MAX_VIEWPORT_DIMS",
		"GL_STEREO",
	}

	GLog("GL Context Params:\n")

	p := make([]int32, 2)
	for i, v := range params {
		if v == gl.STEREO {
			p := make([]bool, 1)
			gl.GetBooleanv(v, p)
			GLog("%s %v\n", names[i], p[0])
			continue
		}

		p[0] = 0
		p[1] = 0
		gl.GetIntegerv(v, p)
		if v == gl.MAX_VIEWPORT_DIMS {
			GLog("%s %d %d\n", names[i], p[0], p[1])
			continue
		}
		GLog("%s %d\n", names[i], p[0])
	}

	GLog("-----------------------------\n")
}

var (
	prevSecs   float64
	frameCount int
	fps        float64
)

func ShowFPS(window *glfw.Window) float64 {
	if !config.FPS {
		return fps
	}

	curSecs := glfw.GetTime()
	elapsedSecs := curSecs - prevSecs
	if elapsedSecs > 0.25 {
		prevSecs = curSecs
		fps = float64(frameCount) / elapsedSecs
		if window != nil {
			window.SetTitle(config.Title + fmt.Sprintf(" @fps: %.2f", fps))
		}
		frameCount = 0
	}
	frameCount++

	return fps
}

/* we can run a full-screen window here */
func fullscreen() (width int, height int, monitor *glfw.Monitor, err error) {
	GLog("Full Screen Mode\n")

	monitor, err = glfw.GetPrimaryMonitor()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get primary monitor")
		return
	}
	vm, err := monitor.GetVideoMode()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get video mode")
		return
	}

	vms, err := monitor.GetVideoModes()
	if err == nil {
		for _, vm := range vms {
			GLog("(%d*%d, %dHZ)\n", vm.Width, vm.Height, vm.RefreshRate)
		}
	}

	name, _ := monitor.GetName()
	GLog("Primary monitor: %s (%d*%d, %dHZ)\n\n",
		name, vm.Width, vm.Height, vm.RefreshRate)

	return vm.Width, vm.Height, monitor, nil
}

func CreateShaderFile(shaderType gl.GLenum, filename string) (gl.Shader, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		GLogErr("ERROR: opening shader file %s: %s\n", filename, err)
		return 0, err
	}

	return CreateShader(shaderType, src)
}

func CreateShader(shaderType gl.GLenum, src []byte) (gl.Shader, error) {
	shader := gl.CreateShader(shaderType)
	shader.Source(string(src))
	shader.Compile()

	if shader.Get(gl.COMPILE_STATUS) == int(gl.FALSE) {
		infoLog := shader.GetInfoLog()
		GLogErr("shader info log for GL index %d:\n%s\n", shader, infoLog)
		shader.Delete()
		return shader, errors.New("Compile: " + infoLog)
	}

	return shader, nil
}

func CreateProgram(shaders ...gl.Shader) (gl.Program, error) {
	program := gl.CreateProgram()

	for _, shader := range shaders {
		program.AttachShader(shader)
	}

	program.Link()
	if program.Get(gl.LINK_STATUS) == int(gl.FALSE) {
		infoLog := program.GetInfoLog()
		GLogErr("Program link info log for GL index %d:\n%s\n", program, infoLog)
		program.Delete()
		return program, errors.New("Link: " + infoLog)
	}

	program.Validate()
	if program.Get(gl.VALIDATE_STATUS) == int(gl.FALSE) {
		infoLog := program.GetInfoLog()
		GLogErr("Program validate info log for GL index %d:\n%s\n", program, infoLog)
		program.Delete()
		return program, errors.New("Validate: " + infoLog)
	}

	return program, nil
}

func ShaderDetail(program gl.Program) {
	fmt.Printf("--------------------\nshader programme %d info:\n", program)
	fmt.Printf("GL_LINK_STATUS = %d\n", program.Get(gl.LINK_STATUS))
	fmt.Printf("GL_ATTACHED_SHADERS = %d\n", program.Get(gl.ATTACHED_SHADERS))
	fmt.Printf("GL_ACTIVE_ATTRIBUTES = %d\n", program.Get(gl.ACTIVE_ATTRIBUTES))

}
