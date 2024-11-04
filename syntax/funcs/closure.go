package main

import "fmt"

func Closure(name string) func() string {
	return func() string {
		return "go," + name
	}
}

func Closure1() func() string {
	name := "大坤坤"
	age := 18
	return func() string {
		return fmt.Sprintf("Hello,%s,%d\n", name, age)
	}
}

func Closure2() func() int {
	age := 0
	fmt.Printf("out: %p\n", &age)
	return func() int {
		fmt.Printf("before: %p\n", &age)
		age++
		fmt.Printf("after: %p\n", &age)
		return age
	}
}
