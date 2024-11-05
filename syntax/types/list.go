package main

/*
基本语法： type 名字 interface {}
  - 里面只能有方法，方法也不需要 func 关键字。

啥是接口（interface）：接口是一组行为的抽象
  - 尽量用接口，以实现面向接口编程。

个人原则：即便是业务开发，也应该面向接口编程。
当你怀疑要不要定义接口的时候，加上去！
*/

type List interface {
	Add(idx int, val any) error //1.18之后用any 之前用interface{}
	Append(val any)
	Delete(idx int) (any, error)
}
