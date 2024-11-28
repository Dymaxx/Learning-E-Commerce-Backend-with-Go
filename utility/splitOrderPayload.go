package utility

import (
	"backenders/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SplitOrderPayload(c *gin.Context) (models.Order, []models.OrderItems, error) {
	var payload models.OrderPayload
	if err := c.ShouldBind(&payload); err != nil {
		return models.Order{}, []models.OrderItems{}, err
	}
	fmt.Println("this is the payload", payload)
	order := models.Order{
		UserID:     payload.UserID,
		TotalPrice: payload.TotalPrice,
		Status:     payload.Status,
	}

	orderItems := payload.Items
	return order, orderItems, nil
}
