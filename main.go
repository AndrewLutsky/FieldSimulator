package main

import (
	"gifhelper"
)

func main() {

	numSteps := 1000
	//numChargesFixed, numChargesFree, width, height, depth int,
	//lowerBoundCharge, upperBoundCharge, lowerBoundMass, upperBoundMass float64
	charges := InitializeRandomCharges(2, 1000, 3000, 3000, 3000,
		-1, 1, 1, 1)

	simulation := make([][]Charge, numSteps)
	simulation[0] = charges

	SimulateFields(simulation, 1, 3000)
	images := AnimateSystem(simulation, 3000, 1, 3000)
	gifhelper.ImagesToGIF(images, "field")
}
