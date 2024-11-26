package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		DryRun: true, //只输出语句，不会执行
	})
	// root:root@tcp(localhost:3306)/your_db
	// 用户名：密码@协议（数据库地址：数据库端口）/数据库名字
	//db.err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/your_db"),&grom.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = db.Debug() //设置debug模式
	// 迁移 schema
	//建表
	db.AutoMigrate(&Product{})

	// Create
	//插入数据
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	//搜索
	var product Product
	db.First(&product, 1)                 // 根据整型主键查找  First查询到的数据只获取一条
	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录 建议使用这种写法

	// Update - 将 product 的 price 更新为 200
	// update不会解决数据类型师傅匹配
	db.Model(&product).Update("Price", 200) //Price列名 200值
	// Update - 更新多个字段
	// 仅更新非零值字段：这句话是说只更新Price和Code两个字段
	// 等价于SET `Price`=200,`Code`='F42'
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"})                    // 仅更新非零值字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"}) //map结构更新

	// Delete - 删除 product
	db.Delete(&product, 1)
}
