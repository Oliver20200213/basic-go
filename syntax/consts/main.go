package main

const External = "包外"
const internal = "包内"

const (
	StatusA = iota
	StatusB
	StatusC
	StatusD

	DayA = iota << 3
	/*
		左移当前数乘以2的位移次幂   iota<<3  4*2的3次方= 4*8=32
		右移则是当前数除以2的位移次幂
	*/
	DayB = 100
	DayC
	DayD
	DayE
)
const (
	Daa = iota*3 + 1
)

func main() {
	const a = 123

}
