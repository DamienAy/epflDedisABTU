package main

import (
	"github.com/DamienAy/epflDedisABTU/ABTU"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	"encoding/json"
	"log"
	"github.com/DamienAy/epflDedisABTU/ABTU/encoding"
)



func main() {
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// --------------------------------------------------------
	// Sample of code for an example of use of and ABTUInstance
	// --------------------------------------------------------

	// All elements needed to start an ABTUInstance, those would be taken from database.
	var siteId SiteId = 1

	var numberOfSites int32 = 4
	var initialSiteTimestamp Timestamp = NewTimestamp(numberOfSites)

	var initialHistoryBuffer []Operation = make([]Operation, 0)

	var initialRemoteBuffer []Operation = make([]Operation, 0)

	// Create an ABTUInstance
	var abtu *ABTU.ABTUInstance
	abtu = ABTU.Init(siteId, initialSiteTimestamp, initialHistoryBuffer, initialRemoteBuffer)

	// Channels needed for the communication with the ABTUInstance
	// !!!! Those will be transformed into channels of bytes for Json objects to be sent !!!!
	var frontendToABTU chan<- []byte
	var ABTUToFrontend <-chan []byte

	//var peersToABTU chan<- []byte
	var ABTUToPeers <-chan []byte

	// Run the ABTUInstance
	frontendToABTU, ABTUToFrontend, _, ABTUToPeers = abtu.Run()

	var char Char = make(Char, 1)
	char[0] = 'a'

	localOperation1 := FrontendOperation{DEL, char, 0}

	encoded, err := json.Marshal(localOperation1)
	check(err)

	frontendMsg := encoding.FrontendMessage{encoding.LocalOp, encoded}

	encodedFrontend, err := json.Marshal(frontendMsg)

	frontendToABTU <- encodedFrontend
	frontendToABTU <- encodedFrontend

	for i:= 0; i<4; i++ {

		select{
		case msg := <- ABTUToFrontend:
			var frontendMsg encoding.FrontendMessage
			err := json.Unmarshal(msg, &frontendMsg)
			check(err)
			log.Println(frontendMsg)
			log.Println("printed Frontend msg")
		case msg := <- ABTUToPeers:
			remoteOp, err := DecodeFromPeers(msg)
			check(err)
			log.Println(remoteOp)
			log.Println("printed peer msg")

		}
	}


}


func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}