package main

import (
	db "SchoolProject/Config"
	routes "SchoolProject/Routes"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func main() {

	fmt.Println("Go School API started...")

	// Connect to the database
	db.Connect()

	// Create a new Fiber application
	app := fiber.New()

	// Set up the routes for the application
	routes.Setup(app)

	// Start the Fiber application on port 30001
	err := app.Listen(":30001")
	if err != nil {
		// If there's an error starting the server, print it to the console
		fmt.Println("Failed to start server:", err)
	}
}
