package types

type Point struct {
	X float64
	Y float64
}

func (p1 Point) Div(p2 Point) Point {
	p1.X /= p2.X
	p1.Y /= p2.Y
	return p1
}

func (p1 Point) Mul(p2 Point) Point {
	p1.X *= p2.X
	p1.Y *= p2.Y
	return p1
}

func (p1 Point) Add(p2 Point) Point {
	p1.X += p2.X
	p1.Y += p2.Y
	return p1
}

func (p1 Point) Sub(p2 Point) Point {
	p1.X -= p2.X
	p1.Y -= p2.Y
	return p1
}
