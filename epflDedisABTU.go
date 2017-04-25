package main

import (
	com "github.com/DamienAy/epflDedisABTU/communication"

	"fmt"
	. "github.com/DamienAy/epflDedisABTU/operation"
	"log"
)


func main() {
	communicationService, err := com.SetupCommunicationService(1, printOp)
	if err != nil{
		log.Fatal("fail.")
	}

	log.Println(communicationService)

	fmt.Println("Press enter when other peers ready")
	var ok string
	fmt.Scanln(&ok)

	//communicationService.Send(Operation{})
}


func printOp(o Operation){
	log.Println(o)
	fmt.Println("Press enter when other peers ready")
	var ok string
	fmt.Scanln(&ok)
}