package main

import "errors"

/*泛型
 */

// 结构泛型
// T类型参数，名字叫做T，约束是any 等于没有约束 可以灵活的规定约束类型
type List[T any] interface {
	Add(idx int, t T)
	Append(t T)
}

func UserList() {
	var l List[int] //可以灵活的执行t的数据类型
	//l.Append("string") //规定了t的类型是int 就不能传字符串
	l.Append(18)

	var l2 List[string]
	l2.Append("string") //此时就可以传入字符串，可以灵活控制出入数据的类型

	var lany List[any]
	lany.Append("any")
	lany.Append(12)
	lany.Append(12.22)

	lk := LinkedList[int]{}
	intVal := lk.head.val
	println(intVal)
}

// 结构体泛型
type LinkedList[T any] struct {
	head *node[T]
	t    T
}

type node[T any] struct {
	val T
}

//方法泛型

func main() {
	println(Sum[int](1, 2, 3))
	println(Sum[Integer](1, 2, 3)) //衍生类型
	println(Sum[float64](1.1, 2.1, 3.1))

}

func Max[T Number](vals ...T) (T, error) {
	if len(vals) == 0 {
		var res T
		return res, errors.New("你的下标不对！")
	}
	res := vals[0]
	if i := 1; i < len(vals) {
		if vals[i] > res {
			res = vals[i]
		}
	}
	return res, nil

}
