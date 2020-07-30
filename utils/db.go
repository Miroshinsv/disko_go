package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func GetConnect() *gorm.DB {
	db, err := gorm.Open("postgres", "host=ec2-107-20-15-85.compute-1.amazonaws.com port=5432 user=wkfbpsznhwjzdl dbname=der0jl2atd0f3j password=1bd72584b189d9431a39b54742db5205aa148e44ec836edaa4f0310323198034 ")
	if err != nil {
		panic(err)
	}
	return db
}
