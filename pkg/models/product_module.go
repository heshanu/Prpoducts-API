package models

import (
	"github.com/heshanu/go-service/pkg/config"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price" gorm:"not null"`
	Quantity    int     `json:"quantity" gorm:"default:0"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Product{})

}
