package Routes

import (
	"example.com/example/Controllers"
	"github.com/gofiber/fiber/v2"
)

func DistributionRoute(route fiber.Router) {
	route.Post("/new", Controllers.DistributionCreateNew)
	route.Post("/get_all/:page/:rowsperpage", Controllers.DistributionGetAll)
	route.Post("/get_all_populated/:page/:rowsperpage", Controllers.DistributionGetAllPopulated)
}
