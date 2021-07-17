package main

import (
	"fmt"
	"log"

	"github.com/pyoko/gorest/pkg/models"
	"gorm.io/gorm"
)

func main() {
	database, err := models.DbConnect()
	if err != nil {
		panic(fmt.Sprintf("unable to cannot to database: %+v", err))
	}	

	// Tables
	database.Transaction(func (tx *gorm.DB) (err error) {
		if err = tx.Exec("DROP TABLE IF EXISTS posts;").Error; err != nil {
			log.Printf("%+v", err)
			return
		}

		if err = tx.AutoMigrate(&models.Post{}); err != nil {
			log.Printf("%+v", err)
			return
		}

		return nil
	})

	// Data
	posts := []models.Post{
		{
			Title: "Post A",
		},
		{
			Title: "Post B",
		},
		{
			Title: "Post C",
		},
		{
			Title: "Post D",
		},
	}

	database.Transaction(func (tx *gorm.DB) (err error) {
		if err = tx.Create(&posts).Error; err != nil {
			log.Printf("%+v", err)
			return
		}

		return nil
	})
}