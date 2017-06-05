package peerCommunication

import (
	"log"
	host "github.com/libp2p/go-libp2p-host"
	"context"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	ma "github.com/multiformats/go-multiaddr"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	"fmt"
	"github.com/libp2p/go-libp2p-protocol"
)

const(
	COMMUNICATION_PROTOCOL protocol.ID = "epflDedisABTU/Broadcast/0.0.1"
)

type CommunicationService struct {
	host host.Host
	PeersToMgmt <-chan []byte
	MgmtToPeers chan<- []byte
}

type ABTUPeer struct {
	Id SiteId
	PeerId string
	IpAddr string
	TCPPort string
}

func check(error error) {
	if error!=nil {
		log.Fatal(error)
	}
}

/*
Sets up a CommunicationService for site myId.
All received operations will be transmitted to the receivingFunction function.
Addresses in the peerCommunication.txt file should in the following format: /ip4/<ipv4Address>/tcp/<tcpPort>/ipfs/<ipfsId>
 */
func Init(myId SiteId, ABTUPeers map[SiteId]ABTUPeer) (*CommunicationService, error) {
	comService := &CommunicationService{}
	comService.PeersToMgmt = make(<-chan []byte, 20)
	comService.MgmtToPeers = make(chan<- []byte, 20)

	// Setup the host
	myIpTcpAddress := fmt.Sprintf("/ip4/%s/tcp/%s", ABTUPeers[myId].IpAddr, ABTUPeers[myId].TCPPort)

	peerId, err := peer.IDB58Decode(ABTUPeers[myId].PeerId)
	if err != nil {
		log.Fatal(err)
	}

	host, err := makeBasicHost(myIpTcpAddress, peerId)
	if err != nil {
		log.Fatal(err)
	}

	comService.host = host

	// Add all other peers into the peerstore.
	for sId, ABTUPeer := range ABTUPeers {
		if sId == myId {
			// skip
		} else {
			multiAddress, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", ABTUPeer.IpAddr, ABTUPeer.TCPPort)); check(err)
			peerId, err := peer.IDB58Decode(ABTUPeer.PeerId); check(err)

			host.Peerstore().AddAddr(peerId, multiAddress, pstore.PermanentAddrTTL)
		}
	}

	return comService, nil
}

func (comService *CommunicationService) Run() {
	// Setup stream handler for COMMUNICATION_PROTOCOL
	comService.host.SetStreamHandler(COMMUNICATION_PROTOCOL , func(s net.Stream) {
		defer s.Close()

		var incommingMsg []byte
		s.Read(incommingMsg)

		comService.PeersToMgmt <- incommingMsg
	})

	// Transfer incoming messages to the PeersToMgmt channel.
	for {
		select {
		case outGoingMsg := <- comService.PeersToMgmt:
			for _, peer := range comService.host.Peerstore().Peers() {
				outGoingStream, err := comService.host.NewStream(context.Background(), peer, string(COMMUNICATION_PROTOCOL))
				if err != nil {
					log.Fatal(err)
				}

				outGoingStream.Write(outGoingMsg)
			}
		}
	}
}

func makeBasicHost(listen string, pid peer.ID) (host.Host, error) {
	multiAddr, err := ma.NewMultiaddr(listen); check(err)

	ps := pstore.NewPeerstore()
	ctx := context.Background()

	// create a new swarm to be used by the service host
	netw, err := swarm.NewNetwork(ctx, []ma.Multiaddr{multiAddr}, pid, ps, nil)
	if err != nil {
		return nil, err
	}

	return bhost.New(netw), nil
}