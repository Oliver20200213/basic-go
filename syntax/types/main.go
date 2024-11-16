package main

import "fmt"

func main() {
	//u1 := &User{}
	//println(u1)
	//fmt.Printf("%p \n", u1)
	var u3 User = User{Name: "tom"}
	fmt.Println(u3)
	u4 := User{"tom", "feng", 18}
	fmt.Println(u4)

	//NewUser()
	//ChangeUser()
	//UserInt()
	//UseFish()

	var l List       //List是接口类型
	l = &ArrayList{} //接口变量存储的是值的类型和值的内存地址，而不是值本身，所以需要用&
	l = &LinkedList{}
	fmt.Println(l)

}
func DoSomething(l List) {

}
