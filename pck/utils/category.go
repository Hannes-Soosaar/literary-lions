package utils

import (
	"database/sql"
	"log"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
)

func GetActiveCategories() []models.Category{
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
return []models.Category{}
}


//TODO CreateCategory