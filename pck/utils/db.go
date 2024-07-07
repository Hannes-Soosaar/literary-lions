package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
	_ "github.com/mattn/go-sqlite3"
)

func CreateDatabase() {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Database created")
}

func InitiateDb() {
	fmt.Println("Opening Database")
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting SQL script Open")
	sql, err := os.ReadFile(config.INIT_SQL)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Println("Database Open", config.INIT_SQL)
		log.Fatal(err)
	}
	fmt.Println("Database populated")
}

//! THIS FUNCTION IS TO BE RAN ONLY ONCE IT WILL HASH ALL PWs IF ALREADY RAN IT WILL DO A HASH OF A HASH  
func PasswordHashing(){
	fmt.Println("Opening Database for initial PW encryption")
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var users []models.User
	queryRead := "SELECT id,password FROM users"
	rows,_ := db.Query(queryRead)
	for rows.Next(){
		var user models.User
		err := rows.Scan(&user.ID,&user.Password)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}
	stmt, err := db.Prepare("UPDATE users SET password=? WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmt.Close()
	for _, user := range users {
		_, err = stmt.Exec(HashString(user.Password),user.ID)
		if err != nil {
			panic(err.Error())
		}
	}
	fmt.Println("PWs encrypted!")
}

func WipeDb() {
	fmt.Println("Resetting Database")
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sql, err := os.ReadFile(config.RESET_SQL)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Println("Failed to execute SQL script")
		log.Fatal(err)
	}
	fmt.Println("Database wiped!")
}
