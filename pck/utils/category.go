package utils

import (
	"database/sql"
	"fmt"
	"log"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
)

func GetActiveCategories() []models.Category {
	var activeCategories []models.Category
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := "SELECT * FROM categories WHERE active = ? ORDER BY category ASC" // Used to the the result in ascending order.
	rows, err := db.Query(query, config.ACTIVE)
	if err != nil {
		fmt.Printf("there is an error getting rows %v \n", err)
		return []models.Category{}
	}
	defer rows.Close()
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Category, &category.Active, &category.CreatedAt)
		if err != nil {
			fmt.Printf("error reading from a row %v  \n", err)
			return activeCategories
		}
		activeCategories = append(activeCategories, category)
	}
	err = rows.Err()
	if err != nil {
		fmt.Printf("error occurred during rows iteration %v \n", err)
		return activeCategories
	}
	return activeCategories
}