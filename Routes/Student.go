package Routes

import (
	"example.com/example/Controllers"
	"github.com/gofiber/fiber/v2"
)

func StudentRoute(route fiber.Router) {
	route.Post("/new", Controllers.StudentCreateNew)
	route.Post("/get_all/:page/:rowsperpage", Controllers.StudentGetAll)
	route.Post("/delete_by_id/:id", Controllers.StudentDelete)
	route.Post("/modify/:id", Controllers.StudentModify)
	route.Post("/modify_status/:id/:new_status", Controllers.StudentSetStatus)
}
