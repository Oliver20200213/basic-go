package main

func YourName(name string, aliases ...string) {
	//aliases是一个切片

}

func CallYourName(name string) {
	YourName("大坤坤")
	YourName("大坤坤", "大坤", "老冯")
	aliases := []string{"大坤", "老冯"}
	YourName(name, aliases...)

}
