package utility

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Convert_params(c *gin.Context) (int, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return 0, fmt.Errorf("id parameter is missing")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("Invalid ID format")
	}

	return id, nil
}
