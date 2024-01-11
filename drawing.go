package main

import (
	"canvas"
	"fmt"
	"image"
)

func AnimateSystem(simulation [][]Charge, canvasWidth, drawingFrequency int, simWidth float64) []image.Image {
	images := make([]image.Image, 0)

	for i := range simulation {
		// if i is divisible by drawingFrequency and append
		if i%drawingFrequency == 0 {
			images = append(images, DrawToCanvas(simulation[i], canvasWidth, simWidth))
		}
	}

	fmt.Println("Created images!")
	return images
}

// DrawToCanvas generates the image corresponding to a canvas after drawing a Universe
// object's bodies on a square canvas that is canvasWidth pixels x canvasWidth pixels
func DrawToCanvas(step []Charge, canvasWidth int, simWidth float64) image.Image {
	// set a new square canvas
	c := canvas.CreateNewCanvas(canvasWidth, canvasWidth)

	// create a black background
	c.SetFillColor(canvas.MakeColor(0, 0, 0))
	c.ClearRect(0, 0, canvasWidth, canvasWidth)
	c.Fill()

	// range over all the bodies and draw them.
	for _, charge := range step {
		if charge.fixed {
			c.SetFillColor(canvas.MakeColor(200, 25, 25))
		} else {
			c.SetFillColor(canvas.MakeColor(125, 200, 125))
		}

		if OutOfRange(charge.position, simWidth) {
			panic("Simulation Crash")
		}

		centerX := (charge.position.x / simWidth) * float64(canvasWidth)
		centerY := (charge.position.y / simWidth) * float64(canvasWidth)
		//centerZ := (charge.position.z / simWidth) * float64(canvasWidth)
		r := (charge.radius / simWidth) * float64(canvasWidth)
		c.Circle(centerX, centerY, r*5)
		c.Fill()
	}
	// we want to return an image!
	return c.GetImage()
}
