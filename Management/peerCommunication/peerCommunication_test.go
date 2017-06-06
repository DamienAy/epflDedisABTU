package peerCommunication

import (
	"testing"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
)



func setupCommunicationService(siteId SiteId) *CommunicationService{
	peer1 := ABTUPeer{1,"QmVvtzcZgCkMnSFf2dnrBPXrWuNFWNM9J3MpZQCvWPuVZf", "127.0.0.1", "1234" }
	peer2 := ABTUPeer{2,"QmT1VesmGjDy4LnGzqSAbkr7ntqh67cgedU2dhsMk7dVGL", "127.0.0.1", "1235" }
	peer3 := ABTUPeer{3,"QmeHfv4QTtnWs12MT3yCAqoc5Vd45WHVuvjAuQkquNY7YM", "127.0.0.1", "1236" }

	ABTUPeers := map[SiteId]ABTUPeer{1:peer1, 2:peer2, 3:peer3}


	return Init(siteId, ABTUPeers)
}
