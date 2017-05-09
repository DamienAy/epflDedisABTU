package ABTU

import (
	com "github.com/DamienAy/epflDedisABTU/communication"

	"fmt"
	. "github.com/DamienAy/epflDedisABTU/operation"
	"log"
	. "github.com/DamienAy/epflDedisABTU/timestamp"
	"sync"
	. "github.com/DamienAy/epflDedisABTU/singleTypes"
	"time"
	. "github.com/DamienAy/epflDedisABTU/remoteBufferManager"
)

var (
	ID SiteId
	SV Timestamp
	H []Operation
	RB []Operation
	lastOp int
	lock sync.Mutex
	RBLock sync.Mutex
	communicationService *com.CommunicationService
)

func main() {
	/*communicationService, err := com.SetupCommunicationService(1, printOp)
	if err != nil{
		log.Fatal("fail.")
	}

	log.Println(communicationService)

	fmt.Println("Press enter when other peers ready")
	var ok string
	fmt.Scanln(&ok)

	//communicationService.Send(Operation{})*/
}

type ABTUInstance struct {
	id SiteId
	sv Timestamp // site timestamp
	h []Operation // history buffer, sorted in effect relation order
	rbm RemoteBufferManager // receiving buffer for remote operations
	lh []Operation // history of local operations for undo.

	manager chan string // channel to manage the instance (stop, etc...)

	lIn chan Operation // channel for receiving from local frontend
	lAckIn chan bool
	lOut chan Operation // channel for sending to local frontend

	rIn chan Operation
	rOut chan Operation
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
	for {
		select {
		case localOperation, done := <- abtu.lIn:


		}
	}
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

