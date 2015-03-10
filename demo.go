// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3.1 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/mmchugh/tmpgraphics/utils"
)

const WindowWidth = 800
const WindowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := utils.NewProgram("vertex_shader.src", "fragment_shader.src")
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(WindowWidth)/WindowHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	texture, err := utils.NewTexture("square.png")
	if err != nil {
		panic(err)
	}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Render
		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.BindVertexArray(vao)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}
