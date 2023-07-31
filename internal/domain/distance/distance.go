package distance

type Distance struct {
	kilometers float64
}

func NewDistance(kilometers float64) Distance {
	return Distance{kilometers: kilometers}
}

func (d Distance) Kilometers() float64 {
	return d.kilometers
}
