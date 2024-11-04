package main

import "fmt"

//defer机制：允许在你返回的前一刻执行一段逻辑

func Defer() {
	defer func() {
		println("第一个defer")
	}()
	defer func() {
		println("第二个defer")
	}()
	defer func() {
		println("第三个个defer")
	}()
}

/*
defer类似于栈，也就是后进先出，也就是先定义的后执行，后定义的先执行。
*/

func DeferClosure() {
	i := 0
	defer func() {
		println(i) //1  先执行i:=0 然后是i=1 然后执行defer
	}()
	i = 1
}
func DeferClosureV1() {
	i := 0
	defer func(i int) {
		println(i)
	}(i) //0  已经把值给到defer了（值传递） 后边再改i的值和defer没关系
	i = 1
}

func DeferReturn() int { //a的地址是同一个
	a := 0
	fmt.Printf("out：%p\n", &a)
	defer func() {
		a = 1
		fmt.Printf("inner：%p\n", &a)
	}()
	return a //0
}

func DeferReturnV1() (a int) {
	a = 0
	fmt.Printf("out：%p\n", &a)
	defer func() { //a的地址是同一个
		fmt.Printf("inner：%p\n", &a)
		a = 1
	}()
	return a //1
}

func DeferClosureLoopV1() {
	for i := 0; i < 10; i++ {
		defer func() {
			println(i)
		}()
	}
}

//每循环迭代一次，执行一个defer语句 i=10时中断循环最后i的值为10

func DeferClosureLoopV2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			println(val)
		}(i)
	}
}

//后入先出 i=10时中断 然后执行10个defer

func DeferClosureLoopV3() {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			println(j)
		}()
	}
}
