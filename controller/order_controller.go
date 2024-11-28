package controller

import (
	"backenders/middleware"
	"backenders/models"
	"backenders/utility"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllOrders(c *gin.Context) {
	orders, err := models.GetAllOrders(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch orders"})
	}
	c.JSON(http.StatusOK, orders)
}

func GetOrderById(c *gin.Context) {
	id, err := utility.Convert_params(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err})
		return
	}
	order, err := models.GetOrdersWithItems(c, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch order"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func GetOrderByUserId(c *gin.Context) {
	user_id, err := utility.Convert_params2(c, "userID")
	fmt.Println("This is the userId", user_id)
	orders, err := models.GetOrderByUserId(c, user_id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch orders"})
	}
	c.JSON(http.StatusOK, orders)
}

func CreateOrder(c *gin.Context) {
	db, err := middleware.GetDB(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect with the database"})
	}

	// Parse JSON payload into the product struct
	order, orderItems, err := utility.SplitOrderPayload(c)
	fmt.Println("This is the order", order)

	// Create the product in the database
	newOrder, err := models.CreateOrder(c, order)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	orderID := newOrder.ID
	fmt.Println(orderItems, orderID)

	for i := 0; i < len(orderItems); i++ {
		orderItems[i].OrderID = orderID
		newOrderItem, err := models.CreateOrderItems(db.Conn, orderItems[i])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create orderItems"})
			return
		}
		fmt.Println(newOrderItem)
	}

	// Respond with the created product
	c.JSON(http.StatusOK, newOrder)
}

func UpdateOrder(c *gin.Context) {
	id, err := utility.Convert_params(c)
	// Parse JSON payload into the order struct
	var updatedOrder models.Order
	if err := c.ShouldBindJSON(&updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the order in the database
	order, err := models.UpdateOrder(c, id, updatedOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	// Respond with the updated product
	c.JSON(http.StatusOK, order)
}
