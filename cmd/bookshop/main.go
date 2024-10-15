package main

import (
	"context"
	"github.com/liam-mackie/book-shop/internal/models"
	"github.com/liam-mackie/book-shop/pkg/data"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var tracer = otel.Tracer("gin-server")

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("failed to create logger")
	}

	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		logger.Panic("failed to connect database", zap.Error(err))
	}

	err = multierr.Append(err, db.AutoMigrate(&models.Author{}))
	err = multierr.Append(err, db.AutoMigrate(&models.Book{}))
	if err != nil {
		logger.Panic("failed to migrate schema", zap.Error(err))
	}

	// Generate some data for testing
	data.Generate(logger, db, 10, 100)

	r := gin.New()
	r.Use(otelgin.Middleware("bookshop"))

	r.GET("/books", getBooks(db))
	r.GET("/books/:id", getBook(db))
	r.POST("/books/:id/buy", sellBook(db))

	err = r.Run(":8080")
	if err != nil {
		logger.Panic("failed to start server", zap.Error(err))
	}

}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func getBooks(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		_, span := tracer.Start(c.Request.Context(), "getBooks")
		defer span.End()

		var books []models.Book
		err := db.Model(&models.Book{}).Preload(clause.Associations).Find(&books).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, books)
	}
	return fn
}

func getBook(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		_, span := tracer.Start(c.Request.Context(), "getBook")
		defer span.End()

		bookId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid book id",
			})
			return
		}

		book := models.Book{
			Model: gorm.Model{ID: uint(bookId)},
		}

		res := db.Model(&models.Book{}).Find(&book)
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if res.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "book not found",
			})
			return
		}

		c.JSON(http.StatusOK, book)
	}
	return fn
}

func sellBook(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		_, span := tracer.Start(c.Request.Context(), "sellBook")
		defer span.End()

		bookId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid book id",
			})
			return
		}

		book := models.Book{
			Model: gorm.Model{ID: uint(bookId)},
		}

		res := db.Model(&models.Book{}).Find(&book)
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if res.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "book not found",
			})
			return
		}

		if book.Sold {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "book already sold",
			})
			return
		}

		book.Sold = true
		db.Save(&book)
		c.JSON(http.StatusOK, book)
	}
	return fn
}
