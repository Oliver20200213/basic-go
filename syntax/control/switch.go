package main

func Switch(stauts int) {
	switch stauts { //switch后面跟有val，case后面跟的是对应的值
	case 0:
		println("初始化")
	case 1:
		println("运行中")
	default:
		println("未知状态")
	}
}

func SwitchBool(age int) {
	switch { //switch后面没val，case后面跟的是bool表达式
	case age >= 18:
		println("成年人")
	case age >= 12:
		println("少年")
	default:
		println("儿童")
	}
}
