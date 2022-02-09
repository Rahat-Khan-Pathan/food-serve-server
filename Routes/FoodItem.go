package Routes

import (
	"example.com/example/Controllers"
	"github.com/gofiber/fiber/v2"
)

func FoodItemRoute(route fiber.Router) {
	route.Post("/new", Controllers.FoodItemCreateNew)
	route.Post("/get_all/:page/:rowsperpage", Controllers.FoodItemGetAll)
	route.Post("/delete_by_id/:id", Controllers.FoodItemDelete)
	route.Post("/modify/:id", Controllers.FoodItemModify)
}
