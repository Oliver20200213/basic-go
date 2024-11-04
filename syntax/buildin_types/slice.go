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

func Subslice() {
	s1 := []int{2, 4, 6, 8, 10}
	s2 := s1[1:3]
	fmt.Printf("s2: %v, len=%d, cap=%d \n", s2, len(s2), cap(s2))

	s3 := s1[2:]
	fmt.Printf("s3: %v, len=%d, cap=%d \n", s3, len(s3), cap(s3))

	s4 := s1[:3]
	fmt.Printf("s4: %v, len=%d, cap=%d \n", s4, len(s4), cap(s4))

}

func ShareSlice() {
	//s1 := []int{1, 2, 3, 4}
	//s2 := s1[2:] // [3,4] 2,2
	//fmt.Printf("s2: %v, len=%d, cap=%d \n", s2, len(s2), cap(s2))
	//s2[0] = 99 // [99,4]
	//fmt.Printf("s2: %v, len=%d, cap=%d \n", s2, len(s2), cap(s2))
	//fmt.Printf("s1: %v, len=%d, cap=%d \n", s1, len(s1), cap(s1))
	//
	//s2 = append(s2, 199) // [99,4,199]  3, 4
	//fmt.Printf("s2: %v, len=%d, cap=%d \n", s2, len(s2), cap(s2))
	//fmt.Printf("s1: %v, len=%d, cap=%d \n", s1, len(s1), cap(s1))

	s1 := []int{1, 2, 3, 4, 5}
	s2 := s1[2:3] // [3] 1 3
	fmt.Printf("s2: %v, len=%d, cap=%d \n", s2, len(s2), cap(s2))
	s2[0] = 99 //[99]
	fmt.Printf("s2: %v, len=%d, cap=%d \n", s2, len(s2), cap(s2))
	fmt.Printf("s1: %v, len=%d, cap=%d \n", s1, len(s1), cap(s1))

	s2 = append(s2, 199) //[99,199]
	fmt.Printf("s2: %v, len=%d, cap=%d \n", s2, len(s2), cap(s2))
	fmt.Printf("s1: %v, len=%d, cap=%d \n", s1, len(s1), cap(s1))

	//扩容会重新创建底层数组

}
