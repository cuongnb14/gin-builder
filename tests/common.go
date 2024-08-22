package tests

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type UserVO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func GetDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&User{})

	return db
}

func CreateUser(db *gorm.DB, total int) {
	for i := 0; i < total; i++ {
		db.Create(&User{Name: "User 1", Email: fmt.Sprintf("user%v@gmail.com", i), Age: 40 + i})
	}
}
