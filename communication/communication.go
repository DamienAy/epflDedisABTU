package communication

import (
	. "github.com/DamienAy/epflDedisABTU/constants"
	. "github.com/DamienAy/epflDedisABTU/operation"
	"os"
	"log"
	"bufio"
	"errors"
	"strconv"
	host "github.com/libp2p/go-libp2p-host"


	"context"
	//"fmt"
	//"log"
	//"strings"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	inet "github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"
	//"os"
	"strings"
)

type CommunicationService struct {
	host host.Host
}

func check(error error) {
	if error!=nil {
		log.Fatal(error)
	}
}

func SetupCommunicationSerice(myId int) (CommunicationService, error) {
	f, err := os.Open("communication.txt")
	if err!= nil {
		log.Fatal(err)
	}
	defer f.Close()

	var addresses [N]string
	scanner := bufio.NewScanner(f)

	for i:=0; i<N; i++ {
		if scanner.Scan(){
			addresses[i] = scanner.Text()
		} else {
			return nil, errors.New("Communication file does contain less than " + strconv.FormatInt(int64(N), 10) + " addresses.")
		}
	}

	address := strings.Split(addresses[i], "/ipfs/")
	myIpTcpAddress := address[0]
	peerId, err := peer.IDB58Decode(address[1]); check(err)

	host, err := makeBasicHost(myIpTcpAddress, peerId); check(err)

	for i:=0; i<N; i++ {
		if i==myId {
			//Skip
		} else {
			address := strings.Split(addresses[i], "/ipfs/")
			multiAddress, err := ma.NewMultiaddr(address[0]); check(err)
			peerId, err := peer.IDB58Decode(address[1]); check(err)

			host.Peerstore().AddAddr(peerId, multiAddress, pstore.PermanentAddrTTL)
		}
	}

	return CommunicationService{host}, nil
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

func (c *CommunicationService) Send(o Operation) {
	//I sent it
}


