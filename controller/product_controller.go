package controller

import (
	"backenders/middleware"
	"backenders/models"
	"backenders/utility"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	// Use the database connection instance (db.Conn)
	database, error := middleware.GetDB(c)
	if error != nil {
		return
	}
	// Use the database connection
	products, err := models.GetAllProducts(database.Conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch products"})
		return
	}
	// Debugging log
	fmt.Println(products, "These are the products")

	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {

	database, error := middleware.GetDB(c)
	if error != nil {
		return
	}
	id, err := utility.Convert_params(c)
	product, err := models.GetProductByID(database.Conn, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch product"})
		return
	}
	fmt.Println(product, "This is the product")

	c.JSON(http.StatusOK, product)

}

func CreateProduct(c *gin.Context) {
	database, error := middleware.GetDB(c)
	if error != nil {
		return
	}

	// Parse JSON payload into the product struct
	var product models.NewProduct
	if err := c.ShouldBindJSON(&product); err != nil {
		fmt.Println(err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the product in the database
	newProduct, err := models.CreateProduct(database.Conn, product)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// Respond with the created product
	c.JSON(http.StatusOK, newProduct)
}

func UpdateProduct(c *gin.Context) {
	database, error := middleware.GetDB(c)
	if error != nil {
		return
	}
	id, err := utility.Convert_params(c)

	// Parse JSON payload into the product struct
	var updatedProduct models.NewProduct
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the product in the database
	product, err := models.UpdateProduct(database.Conn, id, updatedProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Respond with the updated product
	c.JSON(http.StatusOK, product)
}
