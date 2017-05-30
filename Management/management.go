package main

import (
	"github.com/DamienAy/epflDedisABTU/ABTU"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	"github.com/DamienAy/epflDedisABTU/Management/document"
	"fmt"
	"log"
	"github.com/gorilla/websocket"
	"net/http"
	//"encoding/json"
	"strconv"
	"path/filepath"
)

const (
	defaultListenPort = 5050
	maxMessageSize = 1024
	frontendPath = "/Users/knikitin/projects/peer-to-peer-doc-editing"
)


// A structure to communicate with frontend and store documents
type Management struct {
	// Channels to communicate between frontend and management
	controlFromFrontend chan []byte
	controlToFrontend chan []byte

	// A document being opened at the moment
	doc *document.Document

	// Channel to check whether an ABTU instance must be running
	isDocumentOpen chan bool
}


// Returns a pointer to a new Management structure
func newManagement() *Management {
	return &Management{
		controlFromFrontend: make(chan []byte, 20),
		controlToFrontend: make(chan []byte, 20),
		doc: nil,
	}
}


///* Message type to communicate with the front-end*/
//type ControlMessage struct {
//	Event string
//	Content []byte
//}

//func NewControlMessage (event string, content []byte) (*ControlMessage, error) {
//	msg := &ControlMessage{
//		Event: event,
//		Content: content,
//	}
//	return msg, nil
//}


//func (mgmt *Management) handleControlMessages(ws *websocket.Conn) {
//	var cm ControlMessage
//	for {
//		_, message, err := ws.ReadMessage()
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
//				log.Printf("error: %v", err)
//			}
//			break
//		}
//
//		/*Unmarshal the control message { Event: "...", Content: []byte} */
//		err = json.Unmarshal(message, &cm)
//		if err != nil {
//			log.Println("Couldn't decode a management message from frontend:", err)
//		}
//
//		switch cm.Event {
//		case cm.Event != "Document":
//			var dm DocumentMessage
//			var dst interface{}
//
//			if err := json.Unmarshal(cm.Content, &dm); err != nil {
//				log.Println("Couldn't decode a document message from frontend:", err)
//			}
//
//			switch dm.Type {
//			case "OpenDocument":
//				dst = new(DocumentToOpen)
//				if err := json.Unmarshal(dm.Content, dst); err != nil {
//					log.Println("Couldn't decode a content of OpenDocument message:", err)
//				}
//			case "CloseDocument":
//				dst = new(DocumentToClose)
//				if err := json.Unmarshal(dm.Content, dst); err != nil {
//					log.Println("Couldn't decode a content of CloseDocument message:", err)
//				}
//			default:
//				log.Println("Unknown type of document message:", dm.Type)
//
//			}
//		default:
//			log.Println("Wrong control message type:", cm.Event)
//		}
//	}
//}


func serveWS(mgmt *Management, w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, w.Header(), maxMessageSize, maxMessageSize)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	log.Println("Frontend connected to websocket")

	// Writing messages to the connection
	go func() {
		message := <- mgmt.doc.MgmtToFrontend
		// BinaryMessage =2 denotes a binary data message
		ws.WriteMessage(2, message)
	}()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		mgmt.doc.FrontendToMgmt <- message
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
	var numberOfSites uint32 = 4
	var initialSiteTimestamp Timestamp = NewTimestamp(numberOfSites)
	var initialHistoryBuffer []Operation = make([]Operation, 0)
	var initialRemoteBuffer []Operation = make([]Operation, 0)

	// Create an ABTUInstance
	var abtu *ABTU.ABTUInstance
	abtu = ABTU.Init(siteId, initialSiteTimestamp, initialHistoryBuffer, initialRemoteBuffer)

	// Run the ABTUInstance
	doc.FrontendToABTU, doc.ABTUToFrontend, doc.PeersToABTU, doc.ABTUToPeers = abtu.Run()

	return doc
}

func serveHome(mgmt *Management, w http.ResponseWriter, r *http.Request) {
	// Serve requested files and dependencies
	log.Println(r.URL.EscapedPath())
	http.ServeFile(w, r, filepath.Join(frontendPath, r.URL.EscapedPath()))


	/* Currently, creates a document now when a user first time goes to Home
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


func main() {
	/* Create an instance of Management and establish control communication with Frontend*/
	mgmt := newManagement()

	// Give handlers for http and websocket connection and start serving
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHome(mgmt, w, r)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(mgmt, w, r)
	})
	if err := http.ListenAndServe(":" + strconv.Itoa(defaultListenPort), nil); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
	fmt.Printf("Go to http://localhost:%v in your browser to access frontend\n", defaultListenPort)

	for {
		select {
		case isOpen := <-mgmt.isDocumentOpen:
			if !isOpen {
				// Received "close" request
				log.Panicln("Received request to close a doc when none is opened")
			}
		}
		select {
		case message := <- mgmt.doc.FrontendToMgmt:
		//	TODO
		case message := <- mgmt.doc.MgmtToFrontend:
		//	TODO
		case message := <- mgmt.doc.FrontendToABTU:
		//	TODO
		case message := <- mgmt.doc.ABTUToFrontend:
		//	TODO
		case message := <- mgmt.doc.PeersToMgmt:
		//	TODO
		case message := <- mgmt.doc.MgmtToPeers:
		//	TODO
		case message := <- mgmt.doc.PeersToABTU:
		//	TODO
		case message := <- mgmt.doc.ABTUToPeers:
		//	TODO
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

