package common

import (
	"errors"
	"github.com/go-gl/gl"
	"io"
	"io/ioutil"
)

func CreateShader(shaderType gl.GLenum, r io.Reader) (gl.Shader, error) {
	s := gl.CreateShader(shaderType)

	src, err := ioutil.ReadAll(r)
	if err != nil {
		return s, err
	}

	s.Source(string(src))
	s.Compile()

	if s.Get(gl.COMPILE_STATUS) == int(gl.FALSE) {
		return s, errors.New("(compile) - " + s.GetInfoLog())
	}

	return s, nil
}

func CreateProgram(shaders ...gl.Shader) (gl.Program, error) {
	program := gl.CreateProgram()

	for _, shader := range shaders {
		program.AttachShader(shader)
	}

	program.Link()
	if program.Get(gl.LINK_STATUS) == int(gl.FALSE) {
		return program, errors.New("(link) - " + program.GetInfoLog())
	}

	program.Validate()
	if program.Get(gl.VALIDATE_STATUS) == int(gl.FALSE) {
		return program, errors.New("(validate) - " + program.GetInfoLog())
	}

	return program, nil
}
