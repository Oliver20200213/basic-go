package main

import "fmt"

func Map() {
	m1 := map[string]int{
		"key1": 123,
	}
	m1["hello"] = 456

	// 12表示容量，如果不知道容量可以不传，使用默认，默认是16
	m2 := make(map[string]int, 12)
	m2["key2"] = 12

	//map常用接收键值的写法
	val, ok := m1["dakunkun"]
	if ok {
		//有这个键值对
		println(val)
	}

	//对应的键值不存在，有用一个变量来接取map的值，则为对应键值的零值
	val = m1["oliver"]
	println("m1中oliver对应的值：", val) //0

	//删除,无返回值
	delete(m1, "hello")

	//map的遍历是随机的，也就是说你遍历两边，输出的结果都不一样
}

func UseKey() {
	m := map[string]int{
		"key1": 123,
	}
	keys := Keys(m) //编译不通过， map必须是值和键的类型都对的上才行
	fmt.Println(keys)
}

func Keys(m map[string]any) []string {
	return []string{"hello"}
}
