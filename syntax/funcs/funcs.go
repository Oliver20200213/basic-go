package main

import "strings"

// 没有参数
func Func1() { //名字 + 参数列表 + 返回值，这一行也叫做方法签名

}

// 一个参数
func Func2(x int) { //名字 + 参数列表 + 返回值，这一行也叫做方法签名

}

// 多个参数
func Func3(x int, y int) { //名字 + 参数列表 + 返回值，这一行也叫做方法签名

}

// 多个参数，同一类型
func Func4(x, y int) { //名字 + 参数列表 + 返回值，这一行也叫做方法签名

}

//go中同一个包不能重名，不能重载
//func Func4(x, y int) {
//
//}

// 一个返回值
func Func5(a, b string) string {
	return "oliver"
}

// 返回多个返回值
func Func6(a, b string) (string, string) {
	return "oliver", "muzi"
}

// 返回值有名字
func Func7() (name string, age int) {
	return "oliver", 18
}

func Func8() (name string, age int) {
	name = "oliver"
	age = 18
	return
} //和上面等效

func Func9() (name string, age int) {
	//等价于"",0
	//对应类型的零值
	return
}

func func10(abc string) (string, int) {
	s := strings.Split(abc, "")
	return s[0], len(s[0])
}
func func11(abc string) (first string, length int) {
	s := strings.Split(abc, "")
	first = s[0]
	length = len(s[0])
	return
} //扩大了返回值的作用域  两种方式没有啥大的区别  风格问题
