package main

import (
	"fmt"
	"math"
	"math/rand"
)

// Initializes the random fixed point charges and the number of free charges using various parameters.
// Input: Number of fixed and free charges, width, height, and depth of the simulation, the boundaries of
// charge and boundaries of mass.
// Output:A slice of fixed and free charges
func InitializeRandomCharges(numChargesFixed, numChargesFree, width, height, depth int,
	lowerBoundCharge, upperBoundCharge, lowerBoundMass, upperBoundMass float64) []Charge {

	//Create a slice of charges.
	charges := make([]Charge, numChargesFixed+numChargesFree)

	//Range through the amount of fixed charges.
	for i := 0; i < numChargesFixed; i++ {
		var newCharge Charge
		newCharge.position.x = rand.Float64() * float64(width)
		newCharge.position.y = rand.Float64() * float64(height)
		newCharge.position.z = rand.Float64() * float64(depth)
		newCharge.fixed = true

		//Get the random charge value for the charge.
		newCharge.charge = RandomFloatFromInterval(lowerBoundCharge, upperBoundCharge)

		//Get the random mass value for the charge.
		newCharge.mass = RandomFloatFromInterval(lowerBoundMass, upperBoundMass)

		newCharge.radius = 1
		//Add the new charge to the slice.
		charges[i] = newCharge

	}

	//Range through the number of fixed charges.
	for j := numChargesFixed; j < numChargesFixed+numChargesFree; j++ {
		//Create random positions and set it's fixed parameter as false as these are free charges.
		var newCharge Charge
		newCharge.position.x = rand.Float64() * float64(width)
		newCharge.position.y = rand.Float64() * float64(height)
		newCharge.position.z = rand.Float64() * float64(depth)
		newCharge.fixed = false

		//Get the random charge value for the charge.
		newCharge.charge = RandomFloatFromInterval(lowerBoundCharge, upperBoundCharge)

		//Get the random mass value for the charge.
		newCharge.mass = RandomFloatFromInterval(lowerBoundMass, upperBoundMass)

		//Set the radius.
		newCharge.radius = 1

		//Add the new charge to the slice.
		charges[j] = newCharge
	}

	//Return the total slice of charges.
	return charges
}

// Function that generates a random number given a specific interval.
// Input: A lower and upper bound float64.
// Output: Returns a random float64 in that interval.
func RandomFloatFromInterval(lowerBound, upperBound float64) float64 {
	//Generate a random number between 0 and 1.
	randFloat := rand.Float64()

	//Multiply the float by the size of the interval and shift by lower bound.
	randFloat *= upperBound - lowerBound
	randFloat += lowerBound

	//Return the random number in the interval.
	return randFloat
}

// Function to simulate the simulation.
// Takes as input as slice of slice of charges, a timestep, and a simulation width. This returns
// a slice of slice of charges. That is a number of timesteps of a snapshot for each time step.
func SimulateFields(simulation [][]Charge, timeStep, simulationWidth float64) [][]Charge {
	//Set the previous step as the first step of the simulation.
	//Range through the length of the simulation and simulate.
	for i := 1; i < len(simulation); i++ {
		//Print out the step.
		fmt.Print("Step ")
		fmt.Print(i)
		fmt.Println()
		//Find the next step of the simulation.
		simulation[i] = UpdateCharges(simulation[i-1], timeStep, simulationWidth)

	}

	//Returns the simulation.
	return simulation
}

// Function that updates the overall positions, velocities, and accelerations of the charges.
// Takes as input a slice of the simulation, a timestep, and a simulation width.
// Returns the next step in the simulation.
func UpdateCharges(previousStep []Charge, timeStep, simulationWidth float64) []Charge {
	//Copy the previous step.
	nextStep := CopyStep(previousStep)
	//Set the maximum speed of a charge.
	maxSpeed := 10.0
	//Range through each charge, check if its fixed and update its newtonion triples(position, velocity
	//acceleration, etc.)
	for i := range nextStep {
		//Check if the charge is fixed.
		if !nextStep[i].fixed {
			//If not fixed update its newtownian triples.
			//fmt.Println(UpdateAcceleration(previousStep[i], previousStep))
			nextStep[i].acceleration = UpdateAcceleration(previousStep[i], previousStep)
			//fmt.Println(previousStep[i].acceleration, nextStep[i].acceleration)
			nextStep[i].velocity = UpdateVelocity(nextStep[i], timeStep, maxSpeed)
			nextStep[i].position = UpdatePosition(nextStep[i], timeStep, simulationWidth)
		}
	}

	//Return the next step in the simulation.
	return nextStep
}

// Function to update the simulation.
func UpdateAcceleration(chargeToCalculate Charge, totalCharges []Charge) OrderedTriple {
	Force := ComputeNetForce(chargeToCalculate, totalCharges)
	var acceleration OrderedTriple
	acceleration.x = Force.x / chargeToCalculate.mass
	acceleration.y = Force.y / chargeToCalculate.mass
	acceleration.x = Force.z / chargeToCalculate.mass

	return acceleration
}

func UpdateVelocity(chargeToUpdate Charge, timeStep, maxSpeed float64) OrderedTriple {
	var velocity OrderedTriple

	velocity.x = chargeToUpdate.velocity.x + chargeToUpdate.acceleration.x*timeStep
	velocity.y = chargeToUpdate.velocity.y + chargeToUpdate.acceleration.y*timeStep
	velocity.z = chargeToUpdate.velocity.z + chargeToUpdate.acceleration.z*timeStep

	speed := math.Sqrt(velocity.x*velocity.x + velocity.y*velocity.y + velocity.z*velocity.z)
	if speed >= maxSpeed {
		//find ratio of velocity to speed and normalize the velocities in each direction
		velocity.x *= maxSpeed / speed
		velocity.y *= maxSpeed / speed
		velocity.z *= maxSpeed / speed
	}

	return velocity
}

func UpdatePosition(chargeToUpdate Charge, timeStep, simulationWidth float64) OrderedTriple {
	var position OrderedTriple

	time2 := timeStep * timeStep
	position.x = chargeToUpdate.position.x + chargeToUpdate.velocity.x*timeStep + 0.5*chargeToUpdate.acceleration.x*time2
	position.y = chargeToUpdate.position.y + chargeToUpdate.velocity.y*timeStep + 0.5*chargeToUpdate.acceleration.y*time2
	position.z = chargeToUpdate.position.z + chargeToUpdate.velocity.z*timeStep + 0.5*chargeToUpdate.acceleration.z*time2

	//position checker/wrapper
	//Periodic boundary conditions.
	if position.x < 0 {
		position.x += simulationWidth
	} else if position.x > simulationWidth {
		position.x -= simulationWidth
	}
	if position.y < 0 {
		position.y += simulationWidth
	} else if position.y > simulationWidth {
		position.y -= simulationWidth
	}
	if position.z < 0 {
		position.z += simulationWidth
	} else if position.z > simulationWidth {
		position.z -= simulationWidth
	}
	return position

}

func ComputeNetForce(c1 Charge, system []Charge) OrderedTriple {
	var force OrderedTriple
	closePenalty := 10.0
	for i := range system {
		if !CompareCharges(c1, system[i]) {
			//fmt.Println(c1.position, system[i].position)
			newForce := ComputeForce(c1, system[i])

			force.x += newForce.x
			force.y += newForce.y
			force.z += newForce.z

			//Meant to mitigate pbc jumps.
			dist := Distance(c1.position, system[i].position)
			if dist < 5 {
				force.x += closePenalty * (1.0 / dist)
				force.y += closePenalty * (1.0 / dist)
				force.z += closePenalty * (1.0 / dist)
			}
		}

	}
	//Cap out maximum force, otherwise simulation jumps out of PBC.
	maxForce := 20.0

	netForce := math.Sqrt(force.x*force.x + force.y*force.y + force.z*force.z)
	force.x *= maxForce / netForce
	force.y *= maxForce / netForce
	force.z *= maxForce / netForce

	return force
}

func ComputeForce(c1, c2 Charge) OrderedTriple {
	var OP OrderedTriple

	//Find distance between two charges.
	dist := Distance(c1.position, c2.position)

	//Find magnitude of the force.

	k := 9000.0
	F := (k * c1.charge * c2.charge) / (dist * dist)
	//split this into components
	deltaX := c1.position.x - c2.position.x
	deltaY := c1.position.y - c2.position.y
	deltaZ := c1.position.z - c2.position.z

	//Find components of the forces.
	OP.x = F * deltaX / dist
	OP.y = F * deltaY / dist
	OP.z = F * deltaZ / dist
	//Return the pointer of the ordered pair force.
	return OP
}

// Function to copy previous step of the simulation.
// Input: A previous step.
// Output: A deep copy of the same step.(No pointer replication)
func CopyStep(toCopy []Charge) []Charge {
	copyFinal := make([]Charge, len(toCopy))
	//Not safe...
	for i := range toCopy {
		copyFinal[i] = CopyCharge(toCopy[i])
	}

	return copyFinal
}

// Distance takes two position ordered triples and it returns the distance between these two points in 3-D space.
func Distance(p1, p2 OrderedTriple) float64 {
	//Find deltas and take the square root of those squares.
	deltaX := p1.x - p2.x
	deltaY := p1.y - p2.y
	deltaZ := p1.z - p2.z
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY + deltaZ*deltaZ)
}

func CompareCharges(c1, c2 Charge) bool {
	if c1.position != c2.position {
		return false
	} else if c1.acceleration != c2.acceleration {
		return false
	} else if c1.velocity != c2.velocity {
		return false
	} else if c1.charge != c2.charge {
		return false
	} else if c1.mass != c2.mass {
		return false
	} else if c1.fixed != c2.fixed {
		return false
	}

	return true
}

func WithinRadius(c1, c2 Charge) bool {
	radiusC1 := c1.radius
	radiusC2 := c2.radius

	dist := Distance(c1.position, c2.position)

	return radiusC1+radiusC2 > dist
}

// Function to copy the charge.
func CopyCharge(c1 Charge) Charge {
	var newCharge Charge
	newCharge.position = CopyOrderedTriple(c1.position)
	newCharge.velocity = CopyOrderedTriple(c1.velocity)
	newCharge.acceleration = CopyOrderedTriple(c1.acceleration)
	newCharge.mass = c1.mass
	newCharge.radius = c1.radius
	newCharge.fixed = c1.fixed
	newCharge.charge = c1.charge
	return newCharge
}

func OutOfRange(position OrderedTriple, simWidth float64) bool {
	if position.x < 0 || position.x > simWidth {
		return true
	} else if position.y < 0 || position.y > simWidth {
		return true
	} else if position.z < 0 || position.z > simWidth {
		return true
	}

	return false
}

func CopyOrderedTriple(OP OrderedTriple) OrderedTriple {
	var newOP OrderedTriple
	newOP.x = OP.x
	newOP.y = OP.y
	newOP.z = OP.z
	return newOP
}
