package main

import (
	"fmt"
)

/*
结构体的嵌套
可以多层嵌套
可以是指针

• 当 A 组合了 B 之后：
• 可以直接在 A 上调用 B 的方法。
• B 实现的所有接口，都认为 A 已经实现了。
• A 组合 B 之后，在初始化 A 的时候将 B 看做普通
字段来初始化。

注意：
• 组合不是继承，没有多态。

*/

type Inner struct {
}

func (i Inner) DoSomething() {
	println("这是inner")
}
func (i Inner) Name() string {
	return "inner"
}
func (i Inner) SayHello() {
	println("hello,", i.Name())
}

type Outer struct {
	Inner //常用该方法
}

func (o Outer) Name() string {
	return "outer"
}

type OuterV1 struct {
	Inner
}

func (o OuterV1) DoSomething() {
	println("这是outerv1")
}

type OuterPtr struct {
	*Inner //不建议使用，看看就好，除非是调用第三方他返回的就是指针
}

type OOOOuter struct {
	Outer
}

func UseInner() {
	var o Outer
	o.DoSomething()
	o.Inner.DoSomething() //这两个调用是等价的

	var op *OuterPtr
	op.DoSomething()

	//初始化
	o1 := Outer{
		Inner: Inner{},
	}
	op1 := OuterPtr{
		Inner: &Inner{},
	}
	fmt.Println(o1)
	fmt.Println(op1)
}

func main() {
	var o1 OuterV1
	o1.DoSomething()       // 这是outerv1 先找自己的，自己没有的去组合的
	o1.Inner.DoSomething() // 这是inner

	var o Outer
	//输出什么？
	// hello, inner
	// hello, outer
	o.SayHello() // hello.inner  如果是多态的话就会输出hell,outer

}
