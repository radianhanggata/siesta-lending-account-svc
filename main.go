package main

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/radianhanggata/siesta-coding-test/lending-account-svc/config"
	"github.com/radianhanggata/siesta-coding-test/lending-account-svc/lending"
	"github.com/radianhanggata/siesta-coding-test/lending-account-svc/model"
	"github.com/radianhanggata/siesta-coding-test/lending-account-svc/router"
)

func main() {
	app := fiber.New()

	dsn := "host=localhost user=postgres password=123456 dbname=demo port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.Lending{}, &model.Repayment{}, &model.Config{})
	if err != nil {
		panic(err)
	}

	if db.Migrator().HasTable(&model.Config{}) {
		if err := db.First(&model.Config{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			config := model.Config{
				ID:                   "lending",
				Fee:                  5,
				Interest:             1.99,
				OutstandingThreshold: 5000000,
				OutstandingFee:       10000,
			}
			err = db.Create(&config).Error
			if err != nil {
				panic(err)
			}
		}
	}

	configRepo := config.SetupRepository(db)
	lendingRepo := lending.SetupRepository(db)

	lendingHandler := lending.SetupHandler(lendingRepo, configRepo)

	router.SetupRoute(app, lendingHandler)

	log.Fatal(app.Listen(":3000"))
}
