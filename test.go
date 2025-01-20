package main

import (
	"fmt"
	"time"
)

type User struct {
	Id       int
	Email    string
	Password string
	Phone    string
	Ctime    time.Time
}

func main() {
	now := time.Now()
	now = time.UnixMilli(now.UnixMilli())
	u := User{
		Id:       123,
		Email:    "123@qq.com",
		Password: "this is a password",
		Phone:    "15512345678",
		Ctime:    now, //注意这里的时间需要保留到毫秒，time.Now()是到毫秒的
	}
	fmt.Println(u)
}
