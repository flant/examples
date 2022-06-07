package main

import (
	"fmt"
	"strconv"
)

// Функция, вычисляющая площадь прямоугольника.
func getArea(x, y int) (res int) {
	return x * y
}

func main() {
	fmt.Println("Площадь прямоугольника: " + strconv.Itoa(getArea(10, 10)))
}
