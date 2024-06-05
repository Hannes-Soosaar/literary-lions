package main

import (
	"fmt"
	"log"

	utils "gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func main() {

	fmt.Println("Hello Lions!")

	for i := 0; i < 5; i++ {

		Id, err := utils.GenerateUUID()
		if err != nil {
			log.Printf("Error 1! %v", err)
		}
		fmt.Printf("Id %d, is %s", i, Id)

	}

}
