package middleware

import (
	"backenders/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DBMiddleware(database *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", database)
		c.Next()
	}
}

func GetDB(c *gin.Context) (*db.DB, error) {

	dbInstance, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return nil, http.ErrBodyNotAllowed
	}

	database, ok := dbInstance.(*db.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid database connection"})
		return nil, http.ErrNotSupported
	}

	return database, nil
}
