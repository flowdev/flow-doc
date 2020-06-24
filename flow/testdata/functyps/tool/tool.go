package tool

import "fmt"

// Data is just data.
type Data struct {
	Foo string
}

// DoIt just does it.
func DoIt() {
	fmt.Println("Hello world!")
}

// GiveIt returns a number.
func GiveIt() (int, string) {
	return 7, "foo"
}
