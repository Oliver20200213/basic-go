package main

func IfNewVariable(start, end int) string {
	if distance := end - start; distance > 100 {
		return "太远了"
	} else if distance > 60 {
		return "有点远"
	} else {
		return "还挺好"
	}
	//println(distance) //distance只能在if中使用
}
