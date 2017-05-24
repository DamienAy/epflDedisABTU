package document

type Document struct {
	// Channels pass frontend messages from management to ABTUInstance
	frontendToABTU chan []byte
	ABTUToFrontend chan []byte

	// Channels to communicate between management and peers
	peersToMgmt chan []byte
	mgmtToPeers chan []byte

	// Channels to pass remote operations from management to ABTUInstance
	peersToABTU chan []byte
	ABTUToPeers chan []byte
}


