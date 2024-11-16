package main

/*
 type Name struct {
 •  fieldName FieldType
 •  // ...
 • }
结构体和结构体的字段都遵循大小写控制访问性的原则。
通过 . 这个符号来访问字段或者方法。
*/

type LinkedList struct {
	head *node
	tail *node

	//包外可访问
	Len int
}

func (l *LinkedList) Add(idx int, val any) error {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Append(val any) {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Delete(idx int) (any, error) {
	//TODO implement me
	panic("implement me")
}

//func (l LinkedList) Add(idx int, val any) {
//
//}
//
//// 方法接收器  receiver
//func (l *LinkedList) AddV1(idx int, val any) {
//
//}

type node struct {
	prev *node //结构体的自引用一定用指针
	next *node
}
