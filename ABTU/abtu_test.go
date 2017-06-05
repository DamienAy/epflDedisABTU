package ABTU

import (
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	"encoding/json"
	"github.com/DamienAy/epflDedisABTU/ABTU/encoding"
	"github.com/DamienAy/epflDedisABTU/Management/peerCommunication"
	"log"
)

func feedABTU(frontendToABTU chan<- []byte) {
	{
		var char Char = make(Char, 1)
		char[0] = 'a'

		localOperation1 := FrontendOperation{INS, char, 0}

		encoded, err := json.Marshal(localOperation1)
		if err!=nil {
			log.Fatal(err)
		}

		frontendMsg := encoding.FrontendMessage{encoding.LocalOp, encoded}

		encodedFrontend, err := json.Marshal(frontendMsg)

		frontendToABTU <- encodedFrontend
		/*frontendToABTU <- encodedFrontend

		var toUndoIndex int32 = 1
		toUndoIndexBytes, err := json.Marshal(toUndoIndex)
		if err!=nil {
			log.Fatal(err)
		}

		undoFrontendMessage := encoding.FrontendMessage{encoding.Undo, toUndoIndexBytes}
		undoFrontendMessageBytes, err := json.Marshal(undoFrontendMessage)
		if err!=nil {
			log.Fatal(err)
		}

		frontendToABTU <- undoFrontendMessageBytes*/
	}
}

func setupABTUInstance(siteId SiteId) *ABTUInstance {
	// All elements needed to start an ABTUInstance, those would be taken from database.

	var numberOfSites int = 3
	var initialSiteTimestamp Timestamp = NewTimestamp(numberOfSites)
	var initialHistoryBuffer []Operation = make([]Operation, 0)
	var initialRemoteBuffer []Operation = make([]Operation, 0)

	// Create an ABTUInstance
	return Init(siteId, initialSiteTimestamp, initialHistoryBuffer, initialRemoteBuffer)
}

func setupCommunicationService(siteId SiteId) *peerCommunication.CommunicationService{
	peer1 := peerCommunication.ABTUPeer{1,"QmVvtzcZgCkMnSFf2dnrBPXrWuNFWNM9J3MpZQCvWPuVZf", "127.0.0.1", "1234" }
	peer2 := peerCommunication.ABTUPeer{2,"QmT1VesmGjDy4LnGzqSAbkr7ntqh67cgedU2dhsMk7dVGL", "127.0.0.1", "1235" }
	peer3 := peerCommunication.ABTUPeer{3,"QmeHfv4QTtnWs12MT3yCAqoc5Vd45WHVuvjAuQkquNY7YM", "127.0.0.1", "1236" }

	ABTUPeers := map[SiteId]peerCommunication.ABTUPeer{1:peer1, 2:peer2, 3:peer3}


	return peerCommunication.Init(siteId, ABTUPeers)
}