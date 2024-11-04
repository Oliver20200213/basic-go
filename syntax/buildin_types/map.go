package main

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
}
