package router

import (
	"github.com/gofiber/fiber/v2"
	
	"github.com/radianhanggata/siesta-coding-test/lending-account-svc/lending"
)

func SetupRoute(
	app *fiber.App,
	lh lending.Handler,
) {
	api := app.Group("api")

	api.Post("/lending/insert", lh.InsertLending)
	api.Get("/simulate/:id", lh.Simulate)
	api.Get("/repayment/:id", lh.GetRepayment)
}
