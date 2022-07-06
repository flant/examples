package main

import (
	"fmt"
	"strconv"
)

// A function for calculating the area of a rectangle.
func getArea(x, y int) (res int) {
	return x * y
}

func main() {
	fmt.Println("Rectangle area: " + strconv.Itoa(getArea(10, 10)))
}
