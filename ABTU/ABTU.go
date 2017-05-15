package ABTU

import (

	"fmt"
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	"log"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	"time"
	. "github.com/DamienAy/epflDedisABTU/ABTU/remoteBufferManager"
)

type ABTUInstance struct {
	id SiteId
	sv Timestamp // site timestamp
	h []Operation // history buffer, sorted in effect relation order
	rbm RemoteBufferManager // receiving buffer for remote operations
	lh []Operation // history of local operations for undo.

	manager chan string // channel to manage the instance (stop, etc...)

	lIn chan Operation // channel for receiving from local frontend
	lOut chan Operation // channel for sending to local frontend

	rIn chan Operation // channel for receiving remote operations
	rOut chan Operation // channel for dispatching local operations to remote sites.
}

// Initializes the ABTUInstance abtu.
func Init(id SiteId, sv Timestamp, h []Operation, rb []Operation) *ABTUInstance {
	var abtu *ABTUInstance
	abtu.id = id
	abtu.sv = sv

	abtu.h = DeepCopyOperations(h)

	abtu.rbm.Start(rb)

	abtu.manager = make(chan string, 2)

	abtu.lIn = make(chan Operation, 20)
	abtu.lOut = make(chan Operation, 20)

	abtu.rIn = make(chan Operation, 20)
	abtu.rOut = make(chan Operation, 20)

	return abtu
}

func (abtu *ABTUInstance) Run() (chan<- Operation, <-chan Operation, chan<- Operation, <-chan Operation){
	go abtu.run()
	return abtu.lIn, abtu.lOut, abtu.rIn, abtu.rOut
}

func (abtu *ABTUInstance) Stop(){
	abtu.manager <- "stop"
}

func (abtu *ABTUInstance) run() {
	for {
		select {
		case localOperation, notDone := <- abtu.lIn:
			if notDone {}
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

