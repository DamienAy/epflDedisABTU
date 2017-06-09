package ABTU

import (
	"log"
	"testing"
	"time"
	"github.com/DamienAy/epflDedisABTU/ABTU/encoding"
	"encoding/json"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	"github.com/DamienAy/epflDedisABTU/ABTU/operation"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
)

func TestABTUWithCommunication1(t *testing.T) {
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	abtu := setupABTUInstance(1)
	// Run the ABTUInstance
	frontendToABTU, ABTUToFrontend, PeersToABTU , ABTUToPeers := abtu.Run()

	comService := setupCommunicationService(1)
	mgmtToPeers, peersToMgmt := comService.Run()

	time.Sleep(10 * time.Second)

	feedABTU(frontendToABTU)

	go func() {
		for {
			select {
			case msg := <-ABTUToFrontend:
				log.Println("Message to frontend:")
				log.Println(string(msg[:]))
			case msg := <-ABTUToPeers:
				log.Println("Message to peers: ")
				log.Println(string(msg[:]))
				mgmtToPeers <- msg
			case msg := <- peersToMgmt:
				log.Println("Message from peers to ABTU: ")
				log.Println(string(msg[:]))
				PeersToABTU <- msg
			}
		}
	}()

	time.Sleep(3 * time.Second)

}


func TestABTU2Instances(t *testing.T) {
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	abtu1 := setupABTUInstance(1)
	frontendToABTU1, ABTUToFrontend1, PeersToABTU1, ABTUToPeers1 := abtu1.Run()

	abtu2 := setupABTUInstance(2)
	frontendToABTU2, ABTUToFrontend2, PeersToABTU2, ABTUToPeers2 := abtu2.Run()


	feedABTU(frontendToABTU1)

	go func() {
		for {
			select {
			case msg := <-ABTUToFrontend1:
				log.Println("1 - Message to frontend:")
				log.Println(string(msg[:]))
			case msg := <-ABTUToPeers1:
				log.Println("1 - Message to peers: ")
				log.Println(string(msg[:]))
				PeersToABTU2 <- msg
			case msg := <-ABTUToFrontend2:
				log.Println("2 - Message to frontend:")
				log.Println(string(msg[:]))
				bytes, err := json.Marshal(encoding.FrontendMessage{"ackRemoteOperation", []byte{}})
				if err != nil {
					log.Fatal(err)
				}

				frontendToABTU2 <- bytes

			case msg := <- ABTUToPeers2:
				log.Println("2 - Message to peers: ")
				log.Println(string(msg[:]))
				PeersToABTU1 <- msg
			}
		}
	}()

	time.Sleep(5*time.Second)

}


func TestOneABTUInstance(t *testing.T) {
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	abtu1 := setupABTUInstance(1)
	FrontendToABTU1, ABTUToFrontend1, PeersToABTU1, ABTUToPeers1 := abtu1.Run()

	timestamp := NewTimestamp( 2)
	timestamp.Increment(0)
	bytes := []byte("a")
	op := operation.NewOperation(0, INS, 0, bytes, []Timestamp{timestamp}, []Timestamp{}, []Timestamp{}, []Timestamp{}, []Timestamp{})

	bytes, err := op.EncodeToPeers()
	if err != nil {
		log.Fatal(err)
	}

	PeersToABTU1 <- bytes

	var char Char = make(Char, 1)
	char[0] = 'a'

	localOperation1 := operation.FrontendOperation{INS, char, 0}

	encoded, err := json.Marshal(localOperation1)
	if err!=nil {
		log.Fatal(err)
	}

	frontendMsg := encoding.FrontendMessage{encoding.LocalOp, encoded}

	encodedFrontend, err := json.Marshal(frontendMsg)

	FrontendToABTU1 <- encodedFrontend

	go func() {
		for {
			select {
			case msg := <-ABTUToFrontend1:
				log.Println("1 - Message to frontend:")
				log.Println(string(msg[:]))
			case msg := <-ABTUToPeers1:
				log.Println("1 - Message to peers: ")
				log.Println(string(msg[:]))
			}
		}
	}()

	time.Sleep(3*time.Second)

}
