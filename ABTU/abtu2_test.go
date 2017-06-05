package ABTU

import (
	"log"
	"testing"
	"time"
)



func TestABTU2(t *testing.T) {
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// --------------------------------------------------------
	// Sample of code for an example of use of and ABTUInstance
	// --------------------------------------------------------

	abtu := setupABTUInstance(2)
	// Run the ABTUInstance
	_, ABTUToFrontend, PeersToABTU, ABTUToPeers := abtu.Run()

	comService := setupCommunicationService(2)
	mgmtToPeers, peersToMgmt := comService.Run()

	done := make(chan bool)

	go func() {
		for {
			select {
			case msg := <-ABTUToFrontend:
				log.Println("Message to frontend:")
				log.Println(string(msg[:]))
				done <- true
			case msg := <-ABTUToPeers:
				log.Println("Message to peers: ")
				log.Println(string(msg[:]))
				mgmtToPeers <- msg
			case msg := <-peersToMgmt:
				log.Println("Message from peers to ABTU: ")
				log.Println(string(msg[:]))
				PeersToABTU <- msg
			}
		}
	}()

	time.Sleep(15 * time.Second)

	<-done


}


