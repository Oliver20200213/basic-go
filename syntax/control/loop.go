package main

import "fmt"

func ForLoop() {
	for i := 0; i < 10; i++ {
		println(i)
	}

	for i := 0; i < 10; {
		println(i)
		i++
	}

	i := 0
	for ; i < 10; i++ {
		println(i)
	}

}

func Loop2() {
	i := 0
	for i < 10 {
		i++
		println(i)
	}

	for {
		i++
		println(i)
	}
	//或者
	//for true {
	//	i++
	//	println(i)
	//}

}

func ForArr() {
	arr := [3]int{1, 2, 3}
	for index, val := range arr {
		fmt.Println("下标：", index, "值：", val)
	}

	for index := range arr {
		fmt.Println("下标：", index, "值：", arr[index])
	}
}

func ForSlice() {
	arr := []int{1, 2, 3}
	for index, val := range arr {
		fmt.Println("下标：", index, "值：", val)
	}

	for index := range arr {
		fmt.Println("下标：", index, "值：", arr[index])
	}
}

func ForMap() {
	m := map[string]int{
		"key1": 1001,
		"key2": 1002,
	}
	for k, v := range m { //遍历map是随机没有顺序的
		fmt.Println(k, v)
	}
	for k := range m {
		fmt.Println(k, m[k])
	}

}

// 千万不要对迭代参数取地址！！！
/*
在内存里面，迭代参数都是放在一个同一个地方的，
你循环开始就确定了，所以你一旦取地址，那么你
拿到的就是这个地址。
所以，右边的 map 中的键值对的值，最终都是同
一个，也就是最后一个。
*/
func LoopBug() {
	users := []User{
		{
			name: "oliver",
		},
		{
			name: "OLIVER",
		},
	}
	m := make(map[string]*User, 2)
	for _, u := range users {
		fmt.Printf("%p\n", &u)
		m[u.name] = &u
	}
	fmt.Printf("%v\n", m)
}

type User struct {
	name string
}

func LoopContinue() {
	i := 0
	for i < 10 {
		i++
		if i%2 == 0 {
			continue
		}
		println(i)
	}
}
