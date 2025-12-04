package shared

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Logger(app *fiber.App) {
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format:        "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
		DisableColors: true,
	}))
}
