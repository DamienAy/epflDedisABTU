package document

type Document struct {
	// Channels to communicate between frontend and management
	FrontendToMgmt chan []byte
	MgmtToFrontend chan []byte

	// Channels pass frontend messages from management to ABTUInstance
	FrontendToABTU chan []byte
	ABTUToFrontend chan []byte

	// Channels to communicate between management and peers
	PeersToMgmt chan []byte
	MgmtToPeers chan []byte

	// Channels to pass remote operations from management to ABTUInstance
	PeersToABTU chan []byte
	ABTUToPeers chan []byte
}

func NewDocument() *Document {
	return &Document{
		FrontendToMgmt: make(chan []byte, 20),
		MgmtToFrontend: make(chan []byte, 20),
		FrontendToABTU: make(chan []byte, 20),
		ABTUToFrontend: make(chan []byte, 20),
		PeersToMgmt: make(chan []byte, 20),
		MgmtToPeers: make(chan []byte, 20),
		PeersToABTU: make(chan []byte, 20),
		ABTUToPeers: make(chan []byte, 20),
	}
}
