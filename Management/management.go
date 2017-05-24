package main

import (
	"github.com/DamienAy/epflDedisABTU/ABTU"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	. "github.com/DamienAy/epflDedisABTU/Management/Document"
	"runtime"
	"os/exec"
	"fmt"
	"log"
	"github.com/gorilla/websocket"
	"net/http"
)

const (
	defaultListenPort = 11000
	maxMessageSize = 1024
)

// A structure to communicate with frontend and store documents
type Management struct {
	// Channels to communicate between frontend and management
	frontendToMgmt chan []byte
	mgmtToFrontend chan []byte

	// A map of all documents
	documents map[uint32]*Document
}

func newManagement() *Management {
	return &Management{
		frontendToMgmt: make(chan []byte, 20),
		mgmtToFrontend: make(chan []byte, 20),
		documents: make(map[uint32]*Document),
	}
}

/* Message type to communicate with the front-end*/
type ManagementMessage struct {
	MessageType string
	Content []byte
}

func NewManagementMessage (mtype string, content []byte) (*ManagementMessage, error) {
	msg := &ManagementMessage{
		MessageType: mtype,
		Content: content,
	}
	return msg, nil
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	ws, err := websocket.Upgrade(w, r, w.Header(), maxMessageSize, maxMessageSize)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	fmt.Println("Frontend connected")

	//defer func() {
	//	ws.Close()
	//}()
	go handleMessages(ws)
}

func handleMessages(ws *websocket.Conn) {

	for {
		ty, mm, e := ws.ReadMessage()
	}
}


func (m *Management) InitFrontend(listenPort int) error {

	http.HandleFunc("/ws", wsHandler)
	err := http.ListenAndServe(":" + string(listenPort), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

	// Check OS and start CouchDB correspondingly
	switch os := runtime.GOOS; os {
	case "darwin":
		err = exec.Command("open", "../peer-to-peer-doc-editing/index.html").Start()
	case "linux":
		err = exec.Command("xdg-open", "../peer-to-peer-doc-editing/index.html").Start()
	default:
		// freebsd, openbsd, plan9, windows...
		fmt.Printf("Don't know how to start CouchDB or frontend in your OS", "%s.", os)
	}
	if err != nil {
		log.Fatal("Couldn't start frontend", err)
	}


	return nil
}


func main() {

	m := newManagement()

	if err := m.InitFrontend(defaultListenPort); err != nil {
		log.Fatal("Frontend is not started", err)
	}

	// All elements needed to start an ABTUInstance, those would be taken from database.
	var siteId SiteId = 1
	var numberOfSites uint32 = 4

	var initialSiteTimestamp Timestamp = NewTimestamp(numberOfSites)
	var initialHistoryBuffer []Operation = make([]Operation, 0)
	var initialRemoteBuffer []Operation = make([]Operation, 0)

	// TODO
	// Check if requested document exists in management
	// If yes, retrieve the data
	// If not, create newDocument

	// Create an ABTUInstance
	var abtu *ABTU.ABTUInstance
	abtu = ABTU.Init(siteId, initialSiteTimestamp, initialHistoryBuffer, initialRemoteBuffer)

	// Run the ABTUInstance
	m.frontendToABTU, m.ABTUToFrontend, m.peersToABTU, m.ABTUToPeers = abtu.Run()

	for {
		select {
		case <-m.frontendToMgmt:

		case <-m.frontendToABTU:

		case <-m.peersToMgmt:

		case <-m.peersToABTU:

		}
	}


}

