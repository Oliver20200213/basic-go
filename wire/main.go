package wire

import (
	"basic-go/wire/repository"
	"basic-go/wire/repository/dao"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	db, err := gorm.Open(mysql.Open("dsn"))
	if err != nil {
		panic("failed to connect database")
	}
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	fmt.Println(repo)

	InitRepository()
}
