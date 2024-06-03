package main

import (
	"fmt"
	"log"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/uitils"
)

func main() {

	fmt.Println("Hello Lions!")
	
	for i := 0; i < 25; i++ {
		
		Id,err := utils.GenerateUUID()
		if err != nil {
			log.Printf("Error 1! %v",err)
		}
		fmt.Printf("Id %d, is %s",i,Id)
	} 

}