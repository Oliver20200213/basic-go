package main

import (
	"fmt"
)

func main() {
	s := []int{1, 2, 3, 4, 5, 6}
	ret, _ := DeleteAt(s, 3)
	fmt.Println(ret)
}
