package models

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type OrderItems struct {
	ID        int     `db:"id" json:"id"`
	OrderID   int     `db:"order_id" json:"orderID"`
	ProductID int     `db:"product_id" json:"productID"`
	Quantity  int     `db:"quantity" json:"quantity"`
	Price     float64 `db:"price" json:"price"`
}

func GetAllOrderItems(db *sqlx.DB) ([]OrderItems, error) {

	var orderItems []OrderItems
	query := "SELECT * FROM order_items"
	err := db.Select(&orderItems, query)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	return orderItems, nil

}

func CreateOrderItems(db *sqlx.DB, newOrderItem OrderItems) (OrderItems, error) {
	query := `
	INSERT INTO order_items (order_id, product_id, quantity, price) 
	VALUES (:order_id, :product_id, :quantity, :price) 
	RETURNING id, order_id, product_id, quantity, price
`

	var createdOrderItem OrderItems

	rows, err := db.NamedQuery(query, newOrderItem)
	if err != nil {
		log.Println("Error inserting order_item:", err)
		return OrderItems{}, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.StructScan(&createdOrderItem)
		if err != nil {
			log.Panicln("Error scanning orderItems:", err)
		}
	} else {
		return OrderItems{}, err
	}
	fmt.Println("OrderItem created:", createdOrderItem)
	return createdOrderItem, nil
}
