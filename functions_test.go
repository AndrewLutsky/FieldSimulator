package main

import (
	"bufio"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"testing"
)

type DistanceTest struct {
	c1, c2 Charge
	dist   float64
}

// Function that tests the distance method.
func TestDistance(t *testing.T) {

	//Reads in the distance tests from testing/Distance/
	tests := ReadDistanceTests("testing/Distance/")

	//Ranges through the tests.
	for _, test := range tests {

		//Sees what the method does to the input parameters.
		dist := Distance(test.c1.position, test.c2.position)
		//Checks to see if the result is the same.
		if dist != test.dist {
			//Spits out an error if they are not the same.
			t.Errorf("Distance(%v, %v) = %v, want %v", test.c1, test.c2, dist, test.dist)
		}
	}
}

// Read in Distance Tests from folder using directory as input.
func ReadDistanceTests(directory string) []DistanceTest {
	inputFiles := ReadDirectory(directory + "/input")
	numFiles := len(inputFiles)

	//Creates a slice of DistanceTest objects.
	tests := make([]DistanceTest, numFiles)
	for i, inputFile := range inputFiles {
		//Read in Boid1 and Boid2 indicating start and stop lines.
		//This starts reading at line 0 and stops at line 1.
		charges := ReadChargeFromFile(directory+"input/"+inputFile.Name(), 0, 1)

		tests[i].c1 = CopyCharge(charges[0])
		tests[i].c2 = CopyCharge(charges[1])

	}

	//Now, we read in the output files.
	outputFiles := ReadDirectory(directory + "/output")
	//Checks to see if there is the same number of input and output files.
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	//Ranges through the output files.
	for i, outputFile := range outputFiles {

		//Read in the test's result into the test.
		tests[i].dist = ReadFloatFromFile(directory + "output/" + outputFile.Name())
	}

	//Returns the slice of Distance tests.
	return tests
}

// ReadDirectory reads in a directory and returns a slice of fs.DirEntry
// objects containing file info for the directory.
func ReadDirectory(dir string) []fs.DirEntry {
	//read in all files in the given directory
	files, err := os.ReadDir(dir)
	//Checking if there is an error.
	if err != nil {
		panic(err)
	}
	//Return the files from the directory.
	return files
}

// Reads in a specific float value from a file containing ONLY one line.
// Otherwise it will read in the value from the last line. If in a line
// it contains more than one line it might panic.
// Inputs: A file that contains a string of a float.
// Outputs: A float value.
func ReadFloatFromFile(file string) float64 {
	//Reads in the file.
	f, err := os.Open(file)

	//Checks to see if there was an error in reading the file.
	if err != nil {
		panic(err)
	}

	defer f.Close()

	//Creates a new bufio scanner.
	scanner := bufio.NewScanner(f)

	//Creates a value for the string to read into.
	var value string

	//Scanner scans the file.
	for scanner.Scan() {
		//Reads in the line and stores it in the value variable.
		value = scanner.Text()
	}

	//Parses the float from the value string.
	result, err2 := strconv.ParseFloat(value, 64)

	//Checks to see if there was an error in parsing the float from the string.
	if err2 != nil {
		panic(err2)
	}

	//Returns the float parsed from the file.
	return result
}

// Function to read in charges from a file given a starting line and ending line.
// Input: A filename string.
// Returns a slice of charges.
func ReadChargeFromFile(file string, lineStart, lineEnd int) []Charge {
	//Create a slice of charges.
	sliceCharges := make([]Charge, 0)

	//Creates a counter variable to count what line the scanner is on.
	var count int

	//Reads in the file.
	f, err := os.Open(file)

	//Checks to see if there was an error reading the file.
	if err != nil {
		panic(err)
	}

	defer f.Close()

	//Creates a new scanner object.
	scanner := bufio.NewScanner(f)

	//Scan through the file using the bufio scanner object.
	for scanner.Scan() {
		//check to see when to start reading

		//checks to see if we are before the start and after the stop.
		if count < lineStart {
			//Continues to the next line if we are not in the desired range yet.
			count++
			continue
		} else if count > lineEnd {
			//makes sure we don't read in parameters from sky objects
			//as another boid object.
			break
		} else {
			//Read in the line.
			line := scanner.Text()

			//Split the line by spaces.
			values := strings.Split(line, " ")

			//Create a slice of floats to read in from the line.
			floatParameters := make([]float64, len(values))

			//Range through each value in the space-separated slice
			//of strings.
			for i, val := range values {
				//Parse float values and checks to see if there is an error parsing.
				floatVal, err2 := strconv.ParseFloat(val, 64)
				if err2 != nil {
					panic(err2)
				}

				//Read in the parsed float into the array.
				floatParameters[i] = floatVal
			}

			//Creates a new charge object and reads in the float parameters.
			var newCharge Charge
			newCharge.position.x = floatParameters[0]
			newCharge.position.y = floatParameters[1]
			newCharge.position.z = floatParameters[2]
			newCharge.velocity.x = floatParameters[3]
			newCharge.velocity.y = floatParameters[4]
			newCharge.velocity.z = floatParameters[5]
			newCharge.acceleration.x = floatParameters[6]
			newCharge.acceleration.y = floatParameters[7]
			newCharge.acceleration.z = floatParameters[8]
			newCharge.radius = floatParameters[9]
			newCharge.charge = floatParameters[10]

			//Appends the boid to the slice of boids.
			sliceCharges = append(sliceCharges, newCharge)

			//Increases the count to check what line we are on.
			count++
		}
	}

	return sliceCharges
}
