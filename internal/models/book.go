package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	AuthorId uint
	Author   Author
	ISBN     string
	Title    string
	Price    decimal.Decimal
	Sold     bool
}
