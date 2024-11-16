package main

import (
	"fmt"
)

/*
• Go 没有构造函数！！
• 初始化语法：Struct{}
 • 获取指针：
• &Struct{}
 • new(Struct)
 new 可以理解为 Go 会为你的变量分配内存，并且把内存都置为0。

*/

/*
Go 的指针和别的语言的概念一样，本质上都是一个
内存地址。
• 和 C、C++ 一样，*表示指针， &取地址。
• 如果声明了一个指针，但是没有赋值，那么它是
nil。
*/

func NewUser() {
	//初始化结构体
	u := User{}
	fmt.Printf("u %v \n", u)
	fmt.Printf("u %+v \n", u) //带字段打印
	//println(u) //没法直接打印结构体 print只能打印基本类型
	var u1 User
	println(u1.Name)

	//up是一个指针
	up := &User{}
	fmt.Printf("up %+v \n", up)
	up2 := new(User)
	println(up2.FirstName)
	fmt.Printf("up2 %+v \n", up2)
	//两者的效果相同 &{Name: Age:0}

	u4 := User{Name: "oliver", Age: 18}
	u5 := User{Name: "oliver", Age: 18}

	u4.Name = "Jerry"
	u5.Age = 22

	var up3 *User
	// nil上访问字段，或者方法会报错
	//println(up3.FirstName)
	println(up3)
}

type User struct {
	Name      string
	FirstName string
	Age       int
}

/*
方法接收器：结构体接收器、指针接收器
如果是基本类型、结构体，那么就相当于复制了一份
如果是指针，那么就复制了一份指针（也就是地址），但是指向的结构体还是同一个
内置类型和节本类型相似也是相当于复制了一份
*/

func (u User) ChangeName(name string) { // 可以理解为func ChangeName(u User,name string){}
	fmt.Printf("ChangeName中u的地址：%p \n", &u)
	u.Name = name
}

func (u *User) ChangeAge(age int) {
	fmt.Printf("ChangeAge中u的地址：%p \n", u)
	u.Age = age
}

func ChangeUser() {
	u1 := User{Name: "Tom", Age: 18}
	fmt.Printf("u1的地址：%p \n", &u1)
	// 这一步执行的时候，其实相当于复制了一个u1，改的是复制体
	// 所以u1原封不动（值传递）
	u1.ChangeName("Jerry")
	u1.ChangeAge(22) //==>(&u1).ChangeAge(22)
	fmt.Printf("%+v \n", u1)

	up1 := &User{}
	up1.ChangeName("Jerry") //修改name仍然不会生效==》(*up1).ChangeName("Jerry)
	up1.ChangeAge(23)
	fmt.Printf("%+v \n", up1)
}
