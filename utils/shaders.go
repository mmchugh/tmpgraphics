package utils

import(
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
)

func NewProgram(vertexShaderPath, fragmentShaderPath string) (uint32, error) {
	vertexShaderSource, err := ioutil.ReadFile(vertexShaderPath)
	if err != nil {
		return 0, err
	}

	fragmentShaderSource, err := ioutil.ReadFile(fragmentShaderPath)
	if err != nil {
		return 0, err
	}

    vertexShader, err := compileShader(string(vertexShaderSource) + "\x00", gl.VERTEX_SHADER)
    if err != nil {
        return 0, err
    }

    fragmentShader, err := compileShader(string(fragmentShaderSource) + "\x00", gl.FRAGMENT_SHADER)
    if err != nil {
        return 0, err
    }

    program := gl.CreateProgram()

    gl.AttachShader(program, vertexShader)
    gl.AttachShader(program, fragmentShader)
    gl.LinkProgram(program)

    var status int32
    gl.GetProgramiv(program, gl.LINK_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

        return 0, errors.New(fmt.Sprintf("failed to link program: %v", log))
    }

    gl.DeleteShader(vertexShader)
    gl.DeleteShader(fragmentShader)

    return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
    shader := gl.CreateShader(shaderType)

    csource := gl.Str(source)
    gl.ShaderSource(shader, 1, &csource, nil)
    gl.CompileShader(shader)

    var status int32
    gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

        return 0, fmt.Errorf("failed to compile %v: %v", source, log)
    }

    return shader, nil
}

