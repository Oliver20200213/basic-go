package main

//func Recursive() { //会产生臭名昭著的栈溢出 stack overflow
//	Recursive()
//}

// 递归使用不当就有可能出现stack overflow（gorountine的栈）
// 需要加上退出机制
func Recursive(n int) {
	if n > 10 {
		return
	}
	Recursive(n + 1)

}

// 生产中递归stack voerflow产生的常见方式
func A() {
	B()
}

func B() {
	C()
}

func C() {
	A()
}
