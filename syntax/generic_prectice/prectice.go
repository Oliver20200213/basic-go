package generic_prectice

import (
	"errors"
)

func Sum[T Number](ts []T) T {
	var res T
	for _, n := range ts {
		res += n
	}
	return res
}

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func Max[T Number](vals ...T) (T, error) {
	var res T
	if len(vals) == 0 {
		return res, errors.New("输入不正确")
	}
	res = vals[0]
	for i := 1; i < len(vals); i++ {
		if vals[i] > res {
			res = vals[i]
		}
	}
	return res, nil
}

func Min[T Number](vals ...T) (T, error) {
	var res T
	if len(vals) == 0 {
		return res, errors.New("输入不正确")
	}
	res = vals[0]
	for i := 1; i < len(vals); i++ {
		if vals[i] < res {
			res = vals[i]
		}
	}
	return res, nil
}

type matchFunc[T any] func(src T) bool

// type声明新的类型
// matchFunc函数类型的名称   [T any]表示泛型定义 T为类型参数 any为约束
// func(src T) bool  func表示是一个函数类型 (src T)表示函数接收一个类型为T的参数src bool表示返回值的是布尔类型

func find[T any](src []T, match matchFunc[T]) (T, bool) {
	for _, val := range src {
		if match(val) {
			return val, true
		}
	}
	var t T
	return t, false
}

func AddSlice[T any](slice []T, idx int, val T) ([]T, error) {
	//如果idx是负数，或者超过了slice的界限
	if idx < 0 || idx >= len(slice) {
		return nil, errors.New("下标出错")
	}
	res := make([]T, len(slice)+1)
	for i := 0; i < idx; i++ {
		res = append(res, slice[i])
	}
	res[idx] = val
	for i := idx; i < len(slice); i++ {
		res = append(res, slice[i])
	}
	return res, nil

}

func AddSlice2[T any](slice []T, idx int, val T) ([]T, error) {
	if idx < 0 || idx >= len(slice) {
		return nil, errors.New("下标出错了")
	}
	res := append(append(slice[:idx], val), slice[idx+1:]...)
	return res, nil

}
