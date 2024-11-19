package models

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Product struct {
	ID          int     `db:"id"`          // Matches database column 'id'
	Name        string  `db:"name"`        // Matches database column 'name'
	Description string  `db:"description"` // Matches database column 'description'
	Price       float64 `db:"price"`       // Matches database column 'price'
	Stock       int     `db:"stock"`       // Matches database column 'stock'
	CreatedAt   string  `db:"created_at"`  // Matches database column 'created_at'
}

type NewProduct struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func GetAllProducts(db *sqlx.DB) ([]Product, error) {
	var products []Product
	err := db.Select(&products, "SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	return products, nil
}

func GetProductByID(db *sqlx.DB, id int) (Product, error) {
	var product Product
	err := db.Get(&product, "SELECT * FROM products WHERE id = 1")
	if err != nil {
		fmt.Println(err)
		return Product{}, err
	}
	return product, nil
}

func CreateProduct(db *sqlx.DB, product NewProduct) (Product, error) {
	query := `
		INSERT INTO products (name, description, price, stock) 
		VALUES (:name, :description, :price, :stock) 
		RETURNING id, name, description, price, stock, created_at
	`

	var createdProduct Product
	rows, err := db.NamedQuery(query, product)
	if err != nil {
		log.Println("Error inserting product:", err)
		return Product{}, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.StructScan(&createdProduct)
		if err != nil {
			log.Println("Error scanning product:", err)
			return Product{}, err
		}
	} else {
		return Product{}, fmt.Errorf("no rows returned after insert")
	}

	fmt.Println("Product created:", createdProduct)
	return createdProduct, nil
}

func UpdateProduct(db *sqlx.DB, id int, product NewProduct) (Product, error) {
	setClauses := []string{}
	args := map[string]interface{}{
		"id": id,
	}

	if product.Name != "" {
		setClauses = append(setClauses, "name = :name")
		args["name"] = product.Name
	}
	if product.Description != "" {
		setClauses = append(setClauses, "description = :description")
		args["description"] = product.Description
	}
	if product.Price != 0 {
		setClauses = append(setClauses, "price = :price")
		args["price"] = product.Price
	}
	if product.Stock != 0 {
		setClauses = append(setClauses, "stock = :stock")
		args["stock"] = product.Stock
	}

	if len(setClauses) == 0 {
		return Product{}, fmt.Errorf("no fields to update")
	}

	fmt.Println(setClauses)

	// Build the query with the dynamic SET clause
	query := fmt.Sprintf(`
		UPDATE products
		SET %s
		WHERE id = :id
		RETURNING id, name, description, price, stock, created_at
	`, strings.Join(setClauses, ", "))

	// Execute the query
	var updatedProduct Product

	// Run the query
	rows, err := db.NamedQuery(query, args)
	if err != nil {
		log.Println("Error executing query:", err)
		return Product{}, err
	}

	// Use StructScan to scan the result into your updatedProduct struct
	if rows.Next() {
		err := rows.StructScan(&updatedProduct)
		if err != nil {
			log.Println("Error scanning result:", err)
			return Product{}, err
		}
	} else {
		// If no rows were returned, handle the case (e.g., product not found)
		return Product{}, fmt.Errorf("product not found")
	}

	return updatedProduct, nil
}
