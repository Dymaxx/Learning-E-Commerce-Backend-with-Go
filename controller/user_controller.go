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

func GetUsers(c *gin.Context) {
	database, err := middleware.GetDB(c)
	if err != nil {
		return
	}
	users, err := models.GetUsers(database.Conn)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	database, err := middleware.GetDB(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch user"})
		return
	}
	id, err := utility.Convert_params(c)

	user, err := models.GetUserByID(database.Conn, id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch user"})
		return
	}
	fmt.Println(user, "This is the user with id ", id)
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	database, err := middleware.GetDB(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch user"})
		return
	}
	id, err := utility.Convert_params(c)

	// Parse JSON payload into the product struct
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update the user in the database
	user, err := models.UpdateUser(database.Conn, id, updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the user"})
	}

	c.JSON(http.StatusOK, user)

}
