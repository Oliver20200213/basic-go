package main

import (
	"fmt"
)

func main() {
	var a int = 345
	var b int = 123
	fmt.Println(a + b)
	fmt.Println(a - b)
	fmt.Println(a * b)
	fmt.Println(a / b)
	if b != 0 {
		fmt.Println(a / b)
		fmt.Println(a % b)
	}

	var c float64 = 12
	//fmt.Println(a+c)
	fmt.Println(a + int(c))

	//math.Ceil()  数学计算包

	Str()
	Byte()
}
