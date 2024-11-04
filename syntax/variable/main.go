package main

var Global = "全局变量"
var internal = "包内变量，私有变量"

//尽量少用包变量

func main() {
	var a int = 123
	println(a)

	var b = 345
	println(b)

	var c uint = 33
	println(c)

	//println(a == c)
}
