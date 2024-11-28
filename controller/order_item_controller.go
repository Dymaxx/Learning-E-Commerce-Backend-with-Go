package controller

import (
	"backenders/middleware"
	"backenders/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllOrderItems(c *gin.Context) {
	db, err := middleware.GetDB(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Unable to connect to the db"})
	}
	orderItems, err := models.GetAllOrderItems(db.Conn)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Unable to fetch order_items"})
	}

	c.JSON(http.StatusOK, orderItems)
}

func CreateOrderItems(c *gin.Context) {
	db, err := middleware.GetDB(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Unable to connect to the db"})
		return
	}
	var newOrderItem models.OrderItems
	if err := c.ShouldBindJSON(&newOrderItem); err != nil {
		fmt.Println(err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderItems, err := models.CreateOrderItems(db.Conn, newOrderItem)
	// Respond with the created product
	c.JSON(http.StatusOK, orderItems)
}
