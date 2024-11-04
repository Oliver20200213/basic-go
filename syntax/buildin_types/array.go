package main

import "fmt"

func Array() {
	a1 := [3]int{1, 2, 3}
	fmt.Printf("a1: %v, len=%d, cap=%d \n", a1, len(a1), cap(a1))

	a2 := [3]int{9, 8} //可以只声明两个大括号里面只能少不能多
	fmt.Printf("a2: %v, len=%d, cap=%d \n", a2, len(a2), cap(a2))

	var a3 [3]int
	fmt.Printf("a2: %v, len=%d, cap=%d \n", a3, len(a3), cap(a3))

}
