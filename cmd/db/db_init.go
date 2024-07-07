package main

import (
	"fmt"

	utils "gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

/*
This will be a small program  responsible for setting up the initial DB .
It can be executed as a standalone
*/

func main() {
	fmt.Println("Starting the DB initializer")
	utils.WipeDb()
	utils.CreateDatabase()
	utils.InitiateDb()
	utils.PasswordHashing()
}
