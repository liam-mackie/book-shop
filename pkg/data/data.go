package data

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/liam-mackie/book-shop/internal/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"iter"
)

func Generate(logger *zap.Logger, db *gorm.DB, numberAuthors, numberBooks int) {
	for a := range generateAuthors(numberAuthors) {
		if err := db.Create(&a).Error; err != nil {
			logger.Error("failed to create author", zap.Error(err))
		}
		logger.Info("Created author",
			zap.Int("ID", int(a.ID)),
			zap.String("name", a.FirstName+" "+a.LastName),
		)
	}

	for b := range generateBooks(numberBooks, numberAuthors, 0) {
		if err := db.Create(&b).Error; err != nil {
			logger.Error("failed to create book", zap.Error(err))
		}
		logger.Info("Created book",
			zap.Int("ID", int(b.ID)),
			zap.String("title", b.Title),
			zap.String("ISBN", b.ISBN),
			zap.String("price", b.Price.String()),
		)
	}

}

func generateAuthors(number int) iter.Seq[models.Author] {
	return func(yield func(models.Author) bool) {
		for range number {
			generated := models.Author{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
			}
			if !yield(generated) {
				return
			}
		}
	}
}

func generateBooks(number, authorMax, authorIndexStart int) iter.Seq[models.Book] {
	return func(yield func(models.Book) bool) {
		for range number {
			generated := models.Book{
				AuthorId: uint(gofakeit.Number(authorIndexStart, authorMax)),
				ISBN:     gofakeit.Numerify("##########"),
				Title:    gofakeit.BookTitle(),
				Price:    decimal.NewFromFloat(gofakeit.Price(10, 100)),
			}
			if !yield(generated) {
				return
			}
		}
	}
}
