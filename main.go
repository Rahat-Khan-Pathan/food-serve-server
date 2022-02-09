package main

import (
	"fmt"

	"example.com/example/DBManager"
	"example.com/example/Routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func SetupRoutes(app *fiber.App) {
	Routes.FoodItemRoute(app.Group("/food-item"))
	Routes.StudentRoute(app.Group("/student"))
	Routes.DistributionRoute(app.Group("/distribution"))
}

func main() {
	fmt.Println("Hello DakBox")

	fmt.Print("Initializing DataBase Connections ... ")
	initState := DBManager.InitCollections()
	if initState {
		fmt.Println("[OK]")
	} else {
		fmt.Println("[FAILED]")
		return
	}

	fmt.Print("Initializing the server ... ")
	app := fiber.New()
	app.Use(cors.New())
	// app.Use(Middlewares.Auth)
	app.Use(pprof.New())

	SetupRoutes(app)
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		c.Status(200).Send([]byte("Hello DB"))

		return nil
	})
	app.Post("/go", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		c.Status(200).Send([]byte("Hello DB"))

		return nil
	})
	fmt.Println("[OK]")

	app.Listen(":2022")
}
