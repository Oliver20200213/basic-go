package main

func Functional4() {
	println("this is Functional")
}

func UserFunctional4() {
	myFunc := Functional4
	myFunc()
}

// 匿名方法
func functional() {
	//新定义个一个方法 赋值给了fn
	fn := func() string {
		return "go"
	}
	fn()
}

// 匿名方法 立刻发起调用
func functional2() {
	fn := func() string {
		return "go"
	}()
	println(fn)
}

// 返回一个返回string的无参数的方法
func functional3() func() string {
	return func() string {
		return "this is functional"
	}
}
