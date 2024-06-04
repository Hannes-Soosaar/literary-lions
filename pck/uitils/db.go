package utils

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
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
	fmt.Println("Database Open")

	sql, err := ioutil.ReadFile(config.INIT_SQL)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Println("Database Open", config.INIT_SQL)
		log.Fatal(err)
	}

	//This is an example
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int
		err = rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}

}
