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
		homeRouter.DELETE("/:id", controller.DeleteProduct)
	}

	userRouter := router.Group("/users")
	{
		userRouter.GET("/", controller.GetUsers)
		userRouter.GET("/:id", controller.GetUserByID)
		userRouter.PUT("/:id", controller.UpdateUser)
		userRouter.DELETE("/:id", controller.DeleteUser)
	}

	orderRouter := router.Group("/orders")
	{
		orderRouter.GET("/", controller.GetAllOrders)
		orderRouter.GET("/:id", controller.GetOrderById)
		orderRouter.GET("/user/:userID", controller.GetOrderByUserId)
		orderRouter.POST("/", controller.CreateOrder)
		orderRouter.PUT("/:id", controller.UpdateOrder)
	}

	orderItemRouter := router.Group("orderItems")
	{
		orderItemRouter.GET("/", controller.GetAllOrderItems)
		orderItemRouter.POST("/", controller.CreateOrderItems)
	}
}
