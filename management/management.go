package management

import (
	"github.com/DamienAy/epflDedisABTU/ABTU"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	"github.com/DamienAy/epflDedisABTU/management/document"
	"log"
	"github.com/gorilla/websocket"
	"net/http"
	"encoding/json"
	"strconv"
	"path/filepath"
	"github.com/DamienAy/epflDedisABTU/management/peerCommunication"
	"fmt"
)

const (
	defaultListenPort = 5050
	maxMessageSize = 1024
	frontendPath = "/Users/knikitin/projects/peer-to-peer-doc-editing"
)


// A structure to communicate with frontend and store documents
type Management struct {
	// A document being opened at the moment
	doc *document.Document

	// Channel to check whether an ABTU instance must be running
	isDocumentOpen chan bool

	// Communication service to send and receive operations fro network
	network *peerCommunication.CommunicationService
}


// Returns a pointer to a new Management structure
func NewManagement() *Management {
	return &Management{
		doc: nil,
		network: nil,
		isDocumentOpen: make(chan bool),
	}
}


/* Message type to communicate with the front-end and other peers*/
type collaborationMessage struct {
	Event string `json:"Event"`
	Content []byte `json:"Content"`
}

/* A function returning a new instance of collaborationMessage */
func newCollaborationMessage(event string, content []byte) *collaborationMessage {
	msg := &collaborationMessage{
		Event: event,
		Content: content,
	}
	return msg
}


/*Handles messages received from the frontend, either to be sent to an ABTU instance,
an access control message or other ones*/
func (mgmt *Management) handleFrontendMessage(received []byte) {
	var cm collaborationMessage
	err := json.Unmarshal(received, &cm)
	if err != nil {
		log.Println("Error while unmarshalling collaborationMessage:", err)
	}

	switch cm.Event {
	case "ABTU":
		mgmt.doc.FrontendToABTU <- cm.Content
	case "AccessControl":
	//	TODO Handle access control messages
	case "Cursor":
	//	TODO Handle cursor messages
	}
}


/*Handles messages received from the peers, either to be sent to an ABTU instance,
an access control message or other ones*/
func (mgmt *Management) handlePeersMessage(received []byte) {
	var cm collaborationMessage
	err := json.Unmarshal(received, &cm)
	if err != nil {
		log.Println("Error while unmarshalling collaborationMessage:", err)
	}

	switch cm.Event {
	case "ABTU":
		mgmt.doc.PeersToABTU <- cm.Content
	case "AccessControl":
	//	TODO Handle access control messages
	case "Cursor":
	//	TODO Handle cursor messages
	}
}


func serveWS(mgmt *Management, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connected through a websocket")
	ws, err := websocket.Upgrade(w, r, w.Header(), maxMessageSize, maxMessageSize)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	log.Println("Frontend connected to websocket")

	// Writing messages to the connection
	go func() {
		m2write := <- mgmt.doc.MgmtToFrontend
		// BinaryMessage =2 denotes a binary data message
		ws.WriteMessage(2, m2write)
	}()

	for {
		_, m2read, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		mgmt.doc.FrontendToMgmt <- m2read
	}
}


func initDocument() *document.Document {
	//// TODO
	//// Check if requested document exists in management
	//// If yes, retrieve the data
	//// If not, create newDocument
	doc := document.NewDocument()

	//// All elements needed to start an ABTUInstance, those would be taken from database.
	var siteId SiteId = 1
	var numberOfSites int = 2
	var initialSiteTimestamp Timestamp = NewTimestamp(numberOfSites)
	var initialHistoryBuffer []Operation = make([]Operation, 0)
	var initialRemoteBuffer []Operation = make([]Operation, 0)

	// Create an ABTUInstance
	var abtu *ABTU.ABTUInstance
	abtu = ABTU.Init(siteId, initialSiteTimestamp, initialHistoryBuffer, initialRemoteBuffer)
	// Run the ABTUInstance
	doc.FrontendToABTU, doc.ABTUToFrontend, doc.PeersToABTU, doc.ABTUToPeers = abtu.Run()

	/* Setup network communication */
	// Give details of peers
	peer1 := peerCommunication.ABTUPeer{1,"QmVvtzcZgCkMnSFf2dnrBPXrWuNFWNM9J3MpZQCvWPuVZf", "127.0.0.1", "1234" }
	peer2 := peerCommunication.ABTUPeer{2,"QmT1VesmGjDy4LnGzqSAbkr7ntqh67cgedU2dhsMk7dVGL", "127.0.0.1", "1235" }
	ABTUPeers := map[SiteId]peerCommunication.ABTUPeer{1:peer1, 2:peer2}
	// Initialize and run communication service
	comService := peerCommunication.Init(siteId, ABTUPeers)
	doc.MgmtToPeers, doc.PeersToMgmt = comService.Run()

	return doc
}

func serveHome(mgmt *Management, w http.ResponseWriter, r *http.Request) {
	// Serve requested files and dependencies
	//log.Println(r.URL.EscapedPath())
	http.ServeFile(w, r, filepath.Join(frontendPath, r.URL.EscapedPath()))

	/* Currently, creates a document when a user first time goes to Home
	and keeps it open all the time.
	Future: a new document is created and ABTU instance is run
	when a user opens a document (requested from a db),
	then nil the document and stop ABTU when the users quits the doc.
	TODO timely creating and erasing a document instance*/
	if mgmt.doc == nil {
		mgmt.doc = initDocument()
		mgmt.isDocumentOpen <- true
	}
}


func (mgmt *Management) Run() {
	/* Create an instance of Management and establish control communication with Frontend*/
	//mgmt := NewManagement()

	// Give handlers for http and websocket connection and start serving
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHome(mgmt, w, r)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(mgmt, w, r)
	})
	go func() {
		if err := http.ListenAndServe(":" + strconv.Itoa(defaultListenPort), nil); err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()
	fmt.Printf("Go to http://localhost:%v in your browser to access frontend\n", defaultListenPort)

	for {
		select {
		case isOpen := <-mgmt.isDocumentOpen:
			log.Println("A document is opened")
			if !isOpen {
				// Received "close" request
				log.Panicln("Received request to close a doc when none is opened")
			}
		}

		select {

		case received := <- mgmt.doc.FrontendToMgmt:
			log.Println(received)
			mgmt.handleFrontendMessage(received)

		case received := <- mgmt.doc.ABTUToFrontend:
			log.Println(received)
			cm := newCollaborationMessage("ABTU", received)
			message, err := json.Marshal(cm)
			if err != nil {
				log.Println("Error during json marshalling:", err)
			}
			mgmt.doc.MgmtToFrontend <- message

		case received := <- mgmt.doc.PeersToMgmt:
			log.Println(received)
			mgmt.handlePeersMessage(received)

		case received := <- mgmt.doc.ABTUToPeers:
			log.Println(received)
			cm := newCollaborationMessage("ABTU", received)
			message, err := json.Marshal(cm)
			if err != nil {
				log.Println("Error during json marshalling:", err)
			}
			mgmt.doc.MgmtToPeers <- message

		case isOpen := <-mgmt.isDocumentOpen:
			if !isOpen {
				// Document is closed, go back to waiting for opening
				continue
			} else {
				// Received "open" request
				log.Panicln("Received request to open a doc when an opened doc exists")
			}
		}
	}
}

