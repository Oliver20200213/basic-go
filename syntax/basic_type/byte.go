package main

import "fmt"

func Byte() {
	var a byte = 'a'
	println(a) //输出的是a对应的asic吗值
	println(fmt.Printf("%c\n", a))

	var str string = "this is string"
	var bs []byte = []byte(str)
	println(str, bs)

	/*

		!(a&&b) 等价于  !a || !b
		!(a||b) 等价于  !a && !b
	*/

}
