package ABTU

import (
	"log"
	"testing"
	"time"
	"encoding/json"
	"github.com/DamienAy/epflDedisABTU/ABTU/encoding"
)



func TestABTUWithCommunication2(t *testing.T) {
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	abtu := setupABTUInstance(2)
	// Run the ABTUInstance
	frontendToABTU, ABTUToFrontend, PeersToABTU, ABTUToPeers := abtu.Run()

	comService := setupCommunicationService(2)
	mgmtToPeers, peersToMgmt := comService.Run()


	go func() {
		for {
			select {
			case msg := <-ABTUToFrontend:
				log.Println("Message to frontend:")
				log.Println(string(msg[:]))
				bytes, err := json.Marshal(encoding.FrontendMessage{"ackRemoteOperation", []byte{}})
				if err != nil {
					log.Fatal(err)
				}

				frontendToABTU <- bytes
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



}


