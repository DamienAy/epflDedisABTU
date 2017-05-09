package ABTU

import (

	"fmt"
	. "github.com/DamienAy/epflDedisABTU/operation"
	"log"
	. "github.com/DamienAy/epflDedisABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/singleTypes"
	"time"
	. "github.com/DamienAy/epflDedisABTU/remoteBufferManager"
)

type ABTUInstance struct {
	id SiteId
	sv Timestamp // site timestamp
	h []Operation // history buffer, sorted in effect relation order
	rbm RemoteBufferManager // receiving buffer for remote operations
	lh []Operation // history of local operations for undo.

	manager chan string // channel to manage the instance (stop, etc...)

	lIn chan Operation // channel for receiving from local frontend
	lAckIn chan bool // channel for receiving acknowledgements for the execution of remote operations
	lOut chan Operation // channel for sending to local frontend

	rIn chan Operation // channel for receiving remote operations
	rOut chan Operation // channel for dispatching local operations to remote sites.
}

// Initializes the ABTUInstance abtu.
func (abtu *ABTUInstance) Init(id SiteId, sv Timestamp, h []Operation, rb []Operation) (chan<- Operation, chan<- bool, <-chan Operation, chan<- Operation, <-chan Operation){

	abtu.id = id
	abtu.sv = sv

	abtu.h = DeepCopyOperations(h)

	abtu.rbm.Start(rb)

	abtu.manager = make(chan string, 2)

	abtu.lIn = make(chan Operation, 20)
	abtu.lAckIn = make(chan bool, 20)
	abtu.lOut = make(chan Operation, 20)

	abtu.rIn = make(chan Operation, 20)
	abtu.rOut = make(chan Operation, 20)

	return abtu.lIn, abtu.lAckIn, abtu.lOut, abtu.rIn, abtu.rOut
}

func (abtu *ABTUInstance) Run(){
	go abtu.run()
	return
}

func (abtu *ABTUInstance) Stop(){
	abtu.manager <- "stop"
}

func (abtu *ABTUInstance) run() {
	for {
		select {
		case localOperation, done := <- abtu.lIn:
			if done {}
			abtu.LocalThread(localOperation)
		//
		}
	}
}

func (abtu *ABTUInstance) listenToRemote() {
	for {
		remoteOperation, done := <- abtu.rIn
		if done {}
		ack := make(chan bool)
		abtu.rbm.Add <- AddOp{remoteOperation, ack}
		<-ack
	}
}

func (abtu *ABTUInstance) launchController() {

}





func printOp(o Operation){
	log.Println(time.Now())

	var ok string
	fmt.Scanln(&ok)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

