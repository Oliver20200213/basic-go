package main

func Sum[T Number](vals ...T) T { //这里的Number是T的约束，来约束T的数据类型
	var res T
	for _, val := range vals {
		res = res + val
	}
	return res
}

//错误用法：
//func Sum1[T number](vals...number) number{
//	var t T
//  return t
//}

type Number interface {
	~int | float64 | int64 //该接口允许的类型，~int允许int的衍生类型
}

type Integer int
