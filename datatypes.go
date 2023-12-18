package main

type OrderedTriple struct {
	x, y, z float64
}

type Charge struct {
	position, velocity, acceleration OrderedTriple
	charge                           float64
	fixed                            bool
	mass                             float64
	radius                           float64
}

type FieldArrow struct {
	Position  OrderedTriple
	Direction OrderedTriple
	Magnitude float64
}

type Field struct {
	FieldArrows []FieldArrow
}
