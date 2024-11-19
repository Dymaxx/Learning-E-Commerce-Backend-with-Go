package main

import (
	"backenders/db"
	"backenders/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize database connection
	database := db.InitializeDB()
	defer database.Conn.Close()

	// Create Gin router and pass the database
	router := gin.Default()
	routes.InitializeRoutes(router, database)

	router.Run(":8080")
}
