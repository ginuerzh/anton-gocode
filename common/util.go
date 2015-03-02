package common

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
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

func StartGL(title string) (window *glfw.Window, err error) {
	restartGLLog()

	GLog("starting GLFW\n%s\n\n", glfw.GetVersionString())
	/*
		glfw.SetErrorCallback(func(err glfw.ErrorCode, desc string) {
			GLogErr("GLFW ERROR: code %d msg: %s\n", err, desc)
		})
	*/

	if err := glfw.Init(); err != nil {
		GLogErr("ERROR: could not init GLFW3: %s\n", err.Error())
		return nil, err
	}

	glfw.WindowHint(glfw.ContextVersionMajor, config.Major)
	glfw.WindowHint(glfw.ContextVersionMinor, config.Minor)
	if config.Forward && config.Major >= 3 {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)
	}
	if config.Core && config.Major >= 3 && config.Minor >= 2 {
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
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
		gl.Viewport(0, 0, int32(config.Width), int32(config.Height))
	})

	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		GLogErr("ERROR: could not init OpenGL: %s\n", err.Error())
		glfw.Terminate()
		return
	}

	GLog("Vendor: %s\n", gl.GoStr(gl.GetString(gl.VENDOR)))
	GLog("Renderer: %s\n", gl.GoStr(gl.GetString(gl.RENDERER)))
	GLog("Version: %s\n", gl.GoStr(gl.GetString(gl.VERSION)))
	GLog("Shading language version: %s\n", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
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
	params := []uint32{
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
			p := false
			gl.GetBooleanv(v, &p)
			GLog("%s %v\n", names[i], p)
			continue
		}

		p[0] = 0
		p[1] = 0
		gl.GetIntegerv(v, &p[0])
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

	monitor = glfw.GetPrimaryMonitor()
	vm := monitor.GetVideoMode()

	GLog("Primary monitor: %s (%d*%d, %dHZ)\n\n",
		monitor.GetName(), vm.Width, vm.Height, vm.RefreshRate)

	return vm.Width, vm.Height, monitor, nil
}

func CreateShaderFile(shaderType uint32, filename string) (uint32, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		GLogErr("ERROR: opening shader file %s: %s\n", filename, err)
		return 0, err
	}

	return CreateShader(shaderType, src)
}

func getShanderInfoLog(shader uint32) string {
	var length int32
	infoLog := make([]byte, 1)
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
	if length > 0 {
		infoLog = make([]byte, length)
		gl.GetShaderInfoLog(shader, length, nil, &infoLog[0])
	}

	return gl.GoStr(&infoLog[0])
}

func shaderTypeStr(stype uint32) string {
	switch stype {
	case gl.FRAGMENT_SHADER:
		return "fragment shader"
	case gl.VERTEX_SHADER:
		return "vertex shader"
	default:
		return "unknown shader" // TODO: add more shaders
	}
}

func getShaderSource(shader uint32) string {
	var length int32
	ss := make([]byte, 1)
	gl.GetShaderiv(shader, gl.SHADER_SOURCE_LENGTH, &length)
	if length > 0 {
		ss = make([]byte, length)
		gl.GetShaderSource(shader, length, nil, &ss[0])
	}

	return gl.GoStr(&ss[0])
}

func CreateShader(shaderType uint32, src []byte) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	xstring := &src[0]
	gl.ShaderSource(shader, 1, &xstring, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	infoLog := getShanderInfoLog(shader)
	if len(infoLog) > 0 {
		GLogErr("shader info log for GL index %d:\n%s\n", shader, infoLog)
	}

	if status == gl.FALSE {
		gl.DeleteShader(shader)
		return shader, errors.New("Compile: " + infoLog)
	}

	GLog("%s:\n%s\n", shaderTypeStr(shaderType), getShaderSource(shader))

	return shader, nil
}

func getProgramInfoLog(program uint32) string {
	var length int32
	infoLog := make([]byte, 1)
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
	if length > 0 {
		infoLog = make([]byte, length)
		gl.GetProgramInfoLog(program, length, nil, &infoLog[0])
	}

	return gl.GoStr(&infoLog[0])
}

func CreateProgram(shaders ...uint32) (uint32, error) {
	program := gl.CreateProgram()

	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}

	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)

	infoLog := getProgramInfoLog(program)
	if len(infoLog) > 0 {
		GLogErr("Program link info log for GL index %d:\n%s\n", program, infoLog)
	}

	if status == gl.FALSE {
		gl.DeleteProgram(program)
		return program, errors.New("Link: " + infoLog)
	}

	gl.ValidateProgram(program)
	gl.GetProgramiv(program, gl.VALIDATE_STATUS, &status)
	infoLog = getProgramInfoLog(program)
	if len(infoLog) > 0 {
		GLogErr("Program validate info log for GL index %d:\n%s\n", program, infoLog)
	}
	if status == gl.FALSE {
		gl.DeleteProgram(program)
		return program, errors.New("Validate: " + infoLog)
	}

	return program, nil
}

func PrintAll(program uint32) {
	var v int32

	fmt.Printf("--------------------\nshader programme %d info:\n", program)

	gl.GetProgramiv(program, gl.LINK_STATUS, &v)
	fmt.Printf("GL_LINK_STATUS = %d\n", v)

	gl.GetProgramiv(program, gl.VALIDATE_STATUS, &v)
	fmt.Printf("GL_VALIDATE_STATUS = %d\n", v)

	gl.GetProgramiv(program, gl.ATTACHED_SHADERS, &v)
	fmt.Printf("GL_ATTACHED_SHADERS = %d\n", v)

	var activeAttrs int32
	gl.GetProgramiv(program, gl.ACTIVE_ATTRIBUTES, &activeAttrs)
	fmt.Printf("GL_ACTIVE_ATTRIBUTES = %d\n", activeAttrs)

	var length int32
	gl.GetProgramiv(program, gl.ACTIVE_ATTRIBUTE_MAX_LENGTH, &length)

	var i int32
	for i = 0; i < activeAttrs; i++ {
		var size int32
		var typ uint32
		name := make([]byte, length)
		gl.GetActiveAttrib(program, uint32(i), length, nil, &size, &typ, &name[0])
		if size > 1 {
			var j int32
			for j = 0; j < size; j++ {
				longName := fmt.Sprintf("%s[%d]", name, j)
				name := []byte(longName)
				fmt.Printf(" %d) type:%s name:%s location:%d\n",
					i, glType2String(typ), name, gl.GetAttribLocation(program, &name[0]))
			}
		} else {
			fmt.Printf(" %d) type:%s name:%s location:%d\n",
				i, glType2String(typ), name, gl.GetAttribLocation(program, &name[0]))
		}
	}

	var activeUniforms int32
	gl.GetProgramiv(program, gl.ACTIVE_UNIFORMS, &activeUniforms)
	fmt.Printf("GL_ACTIVE_UNIFORMS = %d\n", activeUniforms)

	length = 0
	gl.GetProgramiv(program, gl.ACTIVE_UNIFORM_MAX_LENGTH, &length)

	for i = 0; i < activeUniforms; i++ {
		var size int32
		var typ uint32
		name := make([]byte, length)
		gl.GetActiveUniform(program, uint32(i), length, nil, &size, &typ, &name[0])
		if size > 1 {
			var j int32
			for j = 0; j < size; j++ {
				longName := fmt.Sprintf("%s[%d]", name, j)
				name = []byte(longName)
				fmt.Printf(" %d) type:%s name:%s location:%d\n",
					i, glType2String(typ), name, gl.GetUniformLocation(program, &name[0]))
			}
		} else {
			fmt.Printf(" %d) type:%s name:%s location:%d\n",
				i, glType2String(typ), name, gl.GetUniformLocation(program, &name[0]))
		}
	}

	if infoLog := getProgramInfoLog(program); len(infoLog) > 0 {
		fmt.Printf("Program info log for GL index %d:\n%s", program, infoLog)
	}
}

func glType2String(typ uint32) string {
	switch typ {
	case gl.BOOL:
		return "bool"
	case gl.INT:
		return "int"
	case gl.FLOAT:
		return "float"
	case gl.FLOAT_VEC2:
		return "vec2"
	case gl.FLOAT_VEC3:
		return "vec3"
	case gl.FLOAT_VEC4:
		return "vec4"
	case gl.FLOAT_MAT2:
		return "mat2"
	case gl.FLOAT_MAT3:
		return "mat3"
	case gl.FLOAT_MAT4:
		return "mat4"
	case gl.SAMPLER_2D:
		return "sampler2D"
	case gl.SAMPLER_3D:
		return "sampler3D"
	case gl.SAMPLER_CUBE:
		return "samplerCube"
	case gl.SAMPLER_2D_SHADOW:
		return "sampler2DShadow"
	default:
	}

	return "other"
}
