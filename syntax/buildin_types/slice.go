package main

import "fmt"

func Slice() {
	s1 := []int{1, 2, 3}
	fmt.Printf("s1: %v, len=%d, cap=%d \n", s1, len(s1), cap(s1))

	s2 := make([]int, 3, 4)
	fmt.Printf("s2: %v, len=%d, cap=%d \n", s2, len(s2), cap(s2))

	s3 := make([]int, 4) // {0,0,0,0}  len=cap=4
	//正常写法 make([]int,0,4) 不指定长度
	fmt.Printf("s3: %v, len=%d, cap=%d \n", s3, len(s3), cap(s3))

	//习惯性写法
	s4 := make([]int, 0, 4)
	s4 = append(s4, 1)
	fmt.Printf("s4: %v, len=%d, cap=%d \n", s4, len(s4), cap(s4))

	//在初始化切片的时候要预估容量

}
