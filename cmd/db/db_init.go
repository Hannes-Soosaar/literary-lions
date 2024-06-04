package main

import (
	"fmt"

	utils "gitea.kood.tech/hannessoosaar/literary-lions/pck/uitils"
)

/*
This will be a small program that is responsible for setting up and managing the DB environment.
It can be executed as a standalone
	create the sqlite3 db
	run scripts to populate the db
	run db related tests

*/

func main() {
	fmt.Println("Starting the DB initializer")
	utils.CreateDatabase()
	utils.InitiateDb()
}
