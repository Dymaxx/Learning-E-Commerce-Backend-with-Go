package models

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Email     string `db:"email"`
	Username  string `db:"username"`
	CreatedAt string `db:"created_at"`
}

func GetUsers(db *sqlx.DB) ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return users, nil
}

func GetUserByID(db *sqlx.DB, id int) (User, error) {
	var user User
	query := `SELECT * FROM users WHERE id = $1`
	if err := db.Get(&user, query, id); err != nil {
		fmt.Println(err)
		return User{}, err
	}

	return user, nil
}

func UpdateUser(db *sqlx.DB, id int, user User) (User, error) {
	setClause := []string{}
	args := map[string]interface{}{
		"id": id,
	}

	if user.Name != "" {
		setClause = append(setClause, "name = :name")
		args["name"] = user.Name
	}
	if user.Email != "" {
		setClause = append(setClause, "email = :email")
		args["email"] = user.Email
	}
	if user.Username != "" {
		setClause = append(setClause, "username = :username")
		args["username"] = user.Username

	}

	if len(setClause) == 0 {
		return User{}, fmt.Errorf("no fields to update")
	}

	// Build the query with the dynamic SET clause
	query := fmt.Sprintf(`
		UPDATE USERS
		SET %s
		WHERE id = :id
		RETURNING id, name, email, username
	`, strings.Join(setClause, ", "))

	// Execute the query
	var updatedUser User

	rows, err := db.NamedQuery(query, args)
	if err != nil {
		log.Println("Error executing query", err)
		return User{}, err
	}

	// Use the StructScan to scan the seult into tour updatedUser struct
	if rows.Next() {
		if err := rows.StructScan(&updatedUser); err != nil {
			log.Println("Error scanning result:", err)
			return User{}, err
		}
	} else {
		// If no rows were returnd, handle the case
		return User{}, fmt.Errorf("User not found")
	}
	return updatedUser, nil

}
