package Management

import (
	"github.com/DamienAy/epflDedisABTU/ABTU"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
)

func main() {
	// --------------------------------------------------------
	// Sample of code for an example of use of and ABTUInstance
	// --------------------------------------------------------

	// All elements needed to start an ABTUInstance, those would be taken from database.
	var siteId SiteId = 1

	var numberOfSites uint64 = 4
	var initialSiteTimestamp Timestamp = NewTimestamp(numberOfSites)

	var initialHistoryBuffer []Operation = make([]Operation, 0)

	var initialRemoteBuffer []Operation = make([]Operation, 0)

	// Create an ABTUInstance
	var abtu *ABTU.ABTUInstance
	abtu = ABTU.Init(siteId, initialSiteTimestamp, initialHistoryBuffer, initialRemoteBuffer)

	// Channels needed for the communication with the ABTUInstance
	// !!!! Those will be transformed into channels of bytes for Json objects to be sent !!!!
	var frontendToABTU chan<- Operation
	var ABTUToFrontend <-chan Operation

	var peersToABTU chan<- Operation
	var ABTUToPeers <-chan Operation

	// Run the ABTUInstance
	frontendToABTU, ABTUToFrontend, peersToABTU, ABTUToPeers = abtu.Run()


}

