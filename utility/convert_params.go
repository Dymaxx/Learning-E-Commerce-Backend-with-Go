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
func Convert_params2(c *gin.Context, param string) (int, error) {
	idStr := c.Param(param)
	if idStr == "" {
		fmt.Println("parameter is missing")
		return 0, fmt.Errorf("%d parameter is missing", param)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("invalid userId")
		return 0, fmt.Errorf("Invalid %d format", param)
	}
	fmt.Println("this is the userId %d", id)

	return id, nil
}
