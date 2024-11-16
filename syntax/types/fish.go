package main

import (
	"fmt"
	"time"
)

/*
 基本语法：type TypeA TypeB
我一般在想使用第三方库，又没有办法修改源码，又想
扩展这个库的结构体的方法的情况下，就会用这个。
记住核心：衍生类型是一个全新的类型。
衍生类型可以互相转换，使用 （） 进行转换。
注意， TypeB 实现了某个接口，不等于 TypeA 也实现
了某个接口。
*/

type Integer int //Integer是int的衍生类型

func UserInt() {
	i1 := 10
	i2 := Integer(i1)
	var i3 Integer = 11
	println(i2)
	println(i3)
}

type Fish struct {
	Name string
}

func (f Fish) Swim() {
	println("fish 在游")
}

type FakeFish Fish

func UseFish() {
	f1 := Fish{}
	f2 := FakeFish(f1) //将f1转换成FakeFish这个类型并赋值给f2
	//衍生类型是一个全新的类型 是可以互相转换的
	//f2.Swim() //f2是没有Swim方法的但是可以访问字段
	f2.Name = "Jerry" //修改的是自己的 无法修改f1的
	println(f1.Name)

	var y Yu //和Fish是一模一样的  只是名字不同
	fmt.Println(y.Name)
}

// 衍生类常见用法是：使用第三方库 但是又没法修改源码自己定义额外的方法，就借用衍生类来实现
// 例如：
type MyTime time.Time

func (m MyTime) MyFunc() {

}

/*
类型别名：
 基本语法：type TypeA = TypeB
 • 别名，除了换了一个名字，没有任何区别
• 它和衍生类型的区别，就是用了 =
类型别名一般用在导出类型、兼容性修改里面，也不常
见。
*/

// 向后兼容和Fish类型一模一样只是名字不同
type Yu = Fish
