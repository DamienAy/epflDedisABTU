package main

import (
	com "github.com/DamienAy/epflDedisABTU/communication"

	"fmt"
	. "github.com/DamienAy/epflDedisABTU/operation"
	"log"
	. "github.com/DamienAy/epflDedisABTU/timestamp"
	"sync"
	. "github.com/DamienAy/epflDedisABTU/singleTypes"
)

var (
	ID SiteId
	SV Timestamp
	H []Operation
	lastOp int
	lock sync.Mutex
	communicationService com.CommunicationService
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

func Init(){
	ID = 1;
	H = make([]Operation, 0)
}


func printOp(o Operation){
	log.Println(o)
	fmt.Println("Press enter when other peers ready")
	var ok string
	fmt.Scanln(&ok)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

