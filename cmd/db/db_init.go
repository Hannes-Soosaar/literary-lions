package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	utils "gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
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
	reader := bufio.NewReader(os.Stdin)
	prompt := "Do you want to delete the database? (yes/no)"
	var operation string
	var err error
	for i := 0; i < 2; i++ {
		if i == 0 {
			fmt.Println(prompt)
			operation, err = reader.ReadString('\n')
			operation = strings.TrimSpace(operation)
			if err != nil {
				log.Fatal(err)
			}
		}

		if i == 1 && operation == "yes" {
			prompt = "Are you sure? (yes/no)"
			fmt.Println(prompt)
			operation, err = reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			operation = strings.TrimSpace(operation)

			if operation == "yes" {
				utils.WipeDb()
				fmt.Println("The database has been wiped!")
			}
		}

		if i == 0 && operation != "yes" {
			break
		}
	}
}
