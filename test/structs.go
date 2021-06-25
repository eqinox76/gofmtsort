package test

type Money int

type Car struct {
}

func (c *Car) accelerate() {}

func (c *Car) Drive() {}

func BuyCar(amount Money) *Car {
	return &Car{}
}
