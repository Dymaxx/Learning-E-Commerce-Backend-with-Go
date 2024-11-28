package models

import (
	"backenders/middleware"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type OrderItem struct {
	ProductID   int     `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type Order struct {
	ID         int     `db:"id"`
	UserID     int     `db:"user_id"`
	Status     string  `db:"status"`
	TotalPrice float64 `db:"total_price"`
	CreatedAt  string  `db:"created_at"` // Matches database column 'created_at'
	UpdatedAt  string  `db:"updated_at"`
}

type NewOrder struct {
	UserID     int     `json:"user_id" db:"user_id"`
	TotalPrice float64 `json:"total_price" db:"total_price"`
	Status     string  `json:"status" db:"status"`
}

type RetrunOrder struct {
	ID         int         `json:"id"`
	UserID     int         `json:"userId"`
	Status     string      `json:"status"`
	TotalPrice float64     `json:"totalPrice"`
	Items      []OrderItem `json:"items"`
}

type OrderPayload struct {
	UserID     int          `json:"userID"`
	Status     string       `json:"status"`
	TotalPrice float64      `json:"totalPrice"`
	Items      []OrderItems `json:"items"`
}

func GetAllOrders(c *gin.Context) ([]Order, error) {
	database, err := middleware.GetDB(c)
	if err != nil {
		log.Println("Error", err)
	}
	var orders []Order
	query := `SELECT * FROM orders`
	if err := database.Conn.Select(&orders, query); err != nil {
		fmt.Println(err)
		return []Order{}, err
	}
	fmt.Println(orders)
	return orders, err
}

func GetOrdersWithItems(c *gin.Context, id int) ([]RetrunOrder, error) {

	query := `SELECT 
        o.id AS order_id,
        o.user_id,
        o.total_price,
        oi.quantity,
		p.id AS id,
        p.name AS product_name,
        oi.price AS product_price
    FROM 
        orders o
    JOIN 
        order_items oi ON o.id = oi.order_id
    JOIN 
        products p ON oi.product_id = p.id
    WHERE 
        o.id = $1  
    ORDER BY 
        o.id, oi.id;`

	db, err := middleware.GetDB(c)
	if err != nil {
		log.Println("Error", err)
	}
	rows, err := db.Conn.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a map to group the order items by order
	orderMap := make(map[int]*RetrunOrder)

	for rows.Next() {
		var orderID, userID, productID int
		var totalPrice, price float64
		var quantity int
		var productName string

		if err := rows.Scan(&orderID, &userID, &totalPrice, &quantity, &productID, &productName, &price); err != nil {
			return nil, err
		}

		// Check if the order already exists in the map
		if _, exists := orderMap[orderID]; !exists {
			// If the order does not exist, create a new one and add the first item
			orderMap[orderID] = &RetrunOrder{
				ID:         orderID,
				UserID:     userID,
				TotalPrice: totalPrice,
				Items:      []OrderItem{},
			}
		}

		// Always append the item to the order's Items array, regardless of whether it already exists
		orderMap[orderID].Items = append(orderMap[orderID].Items, OrderItem{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    quantity,
			Price:       price,
		})
	}

	// Convert the map to a slice for easier return
	var result []RetrunOrder
	for _, order := range orderMap {
		result = append(result, *order)
	}

	fmt.Println(result)

	return result, nil
}

func GetOrderByUserId(c *gin.Context, userId int) ([]Order, error) {

	db, err := middleware.GetDB(c)
	if err != nil {
		log.Println("Faild to get database connection:", err)
	}

	query := `SELECT * FROM orders WHERE user_id = $1`
	var orders []Order
	// Retrieve orders
	err = db.Conn.Select(&orders, query, userId)
	if err != nil {
		log.Println("Error executing query:", err)
		return []Order{}, err
	}
	if len(orders) <= 0 {
		log.Println("There are no orders found with this user_id")
		return []Order{}, fmt.Errorf("Error", fmt.Sprintf("There are no orders with this user_id: %d", userId))

	}

	return orders, err

}

func UpdateOrder(c *gin.Context, id int, order Order) (Order, error) {
	setClause := []string{}
	args := map[string]interface{}{
		"id": id,
	}

	if order.UserID >= 0 {
		setClause = append(setClause, "user_id = :user_id")
		args["user_id"] = order.UserID
	}
	if order.TotalPrice >= 0 {
		setClause = append(setClause, "total_price = :total_price")
		args["total_price"] = order.TotalPrice
	}
	if order.Status != "" {
		setClause = append(setClause, "status = :status")
		args["status"] = order.Status

	}
	if order.UpdatedAt != "" {
		setClause = append(setClause, "updated_at = :updated_at")
		args["updated_at"] = order.UpdatedAt

	}

	if len(setClause) == 0 {
		return Order{}, fmt.Errorf("no fields to update")
	}

	// Build the query with the dynamic SET clause
	query := fmt.Sprintf(`
		UPDATE Orders
		SET %s
		WHERE id = :id
		RETURNING id, user_id, total_price, status, updated_at
	`, strings.Join(setClause, ", "))

	// Execute the query
	var updatedOrder Order
	db, err := middleware.GetDB(c)
	if err != nil {
		log.Println("Failed to get database connection:", err)
		return Order{}, err
	}
	rows, err := db.Conn.NamedQuery(query, args)
	if err != nil {
		log.Println("Error executing query", err)
		return Order{}, err
	}

	// Use the StructScan to scan the seult into tour updatedUser struct
	if rows.Next() {
		if err := rows.StructScan(&updatedOrder); err != nil {
			log.Println("Error scanning result:", err)
			return Order{}, err
		}
	} else {
		// If no rows were returnd, handle the case
		return Order{}, fmt.Errorf("User not found")
	}
	return updatedOrder, nil

}

func CreateOrder(c *gin.Context, order Order) (Order, error) {
	db, err := middleware.GetDB(c)
	if err != nil {
		fmt.Println("Failed to get the database connection")
		return Order{}, err
	}
	query := `
		INSERT INTO orders (user_id, total_price, status) 
		VALUES (:user_id, :total_price, :status) 
		RETURNING id, user_id,total_price, status
	`
	var createdOrder Order
	fmt.Printf("Order data being inserted: %+v\n", order)
	rows, err := db.Conn.NamedQuery(query, order)
	if err != nil {
		log.Println("Error inserting product:", err)
		return Order{}, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.StructScan(&createdOrder)
		if err != nil {
			log.Println("Error scanning product:", err)
			return Order{}, err
		}
	} else {
		return Order{}, fmt.Errorf("no rows returned after insert")
	}

	fmt.Println("Order created:", createdOrder)
	return createdOrder, nil
}

func DeleteOrder(c *gin.Context, id int) error {
	query := `DELETE FROM orders WHERE id = :id`
	db, err := middleware.GetDB(c)
	if err != nil {
		log.Println("Failed to get database connection:", err)
		return err
	}
	result, err := db.Conn.NamedExec(query, map[string]interface{}{"id": id})
	if err != nil {

		return fmt.Errorf("Error", err)
	}
	// Check if there is a order deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no orders found with id %d ", id)
	}

	return nil
}
