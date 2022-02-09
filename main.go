package main

import (
	"fmt"
	"os"

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
	fmt.Println("[OK]")
	fmt.Println(os.Getenv("PORT"))
	app.Listen(":" + os.Getenv("PORT"))
}
