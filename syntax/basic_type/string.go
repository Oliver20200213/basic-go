package main

import (
	"fmt"
	"unicode/utf8"
)

func Str() {
	// He said: "Hello Go!"
	fmt.Println("He said: \"Hello Go!\"")
	println("Hello, Go!")
	println(`反引号中不能有反引号没法转义
可以换行
“

下一行`)

	println(len("你好")) //先转成utf8在计算
	println(utf8.RuneCountInString("你好abc"))

}
