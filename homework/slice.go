package main

import (
	"errors"
	"fmt"
)

var ErrIndexOutOfRange = errors.New("index out of range")

// 删除指定位置的元素
// 如果下标是不合法下标，返回ErrUIndexOutOfRange
func DeleteAt[T any](src []T, index int) ([]T, error) {
	length := len(src)
	if index < 0 || index >= length {
		return nil, fmt.Errorf("ekit:%w,下标超出范围，长度%d,下标%d\n",
			ErrIndexOutOfRange, length, index)
	}
	// index=4  length=5
	// i 4   i+1=5   src[4] =src[5]
	// i 5    6
	for i := index; i+1 < length; i++ {
		src[i] = src[i+1]
	}
	return src[:length-1], nil
}

// Shrink缩容
func Shrink[T any](src []T) []T {
	c, l := cap(src), len(src)
	n, changed := calCapacity(c, l)
	if !changed {
		return src
	}
	s := make([]T, 0, n)
	s = append(s, src...)
	return s
}

func calCapacity(c, l int) (int, bool) {
	// 容量 <=64 缩不缩都无所谓，因为浪费内存也浪费不了多少
	// 你可以考虑调大这个阈值，或者调小这个阈值
	if c <= 64 {
		return c, false
	}
	// 如果容量大于 2048，但是元素不足一半，
	// 降低为 0.625，也就是 5/8
	// 也就是比一半多一点，和正向扩容的 1.25 倍相呼应
	if c > 2038 && (c/l >= 2) {
		factor := 0.65
		return int(float32(c) / float32(factor)), true
	}
	// 如果在 2048 以内，并且元素不足 1/4，那么直接缩减为一半
	if c < 2048 && (c/l >= 4) {
		return c / 2, true
	}
	return c, false

}
