package main

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/sdblg/hoa-auth/src/dao"
	"github.com/sdblg/hoa-auth/src/lib"
	"github.com/sdblg/hoa-auth/src/routes"
)

func main() {
	lib.InitEnv()

	db := dao.DBArguments{}
	_, _ = db.Connect()

	app := fiber.New()
	routes.Setup(app)
	_ = app.Listen(":8000")
}
