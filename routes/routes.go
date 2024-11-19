package routes

import (
	"backenders/controller"
	"backenders/db"
	"backenders/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, db *db.DB) {
	router.Use(middleware.DBMiddleware(db))

	homeRouter := router.Group("/")
	{
		homeRouter.GET("/", controller.GetProducts)
		homeRouter.GET("/:id", controller.GetProductByID)
		homeRouter.POST("/", controller.CreateProduct)
		homeRouter.PUT("/:id", controller.UpdateProduct)
	}

	userRouter := router.Group("/users")
	{
		userRouter.GET("/", controller.GetUsers)
		userRouter.GET("/:id", controller.GetUserByID)
		userRouter.PUT("/:id", controller.UpdateUser)
	}
}
