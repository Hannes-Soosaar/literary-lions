package main

import (
	"fmt"

	utils "gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func main() {

	user := utils.FindUserByUserName("bob")
	fmt.Printf("User in main %v 'n", user)

}
