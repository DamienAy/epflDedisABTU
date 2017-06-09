package ABTU

import (

	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	"log"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	. "github.com/DamienAy/epflDedisABTU/ABTU/receivingBufferManager"
	"errors"
	"encoding/json"
	. "github.com/DamienAy/epflDedisABTU/ABTU/encoding"
)

type ABTUInstance struct {
	id                    SiteId
	sv                    Timestamp              // site timestamp
	h                     []Operation            // history buffer, sorted in effect relation order
	localTimestampHistory []Timestamp            // buffer storing timestamps of operations generated locally.
	rbm                   ReceivingBufferManager // receiving buffer for remote operations

	localToABTU chan []byte // channel for receiving from local frontend
	aBTUToLocal chan []byte // channel for sending to local frontend

	remoteToABTU chan []byte // channel for receiving remote operations
	aBTUToRemote chan []byte // channel for dispatching local operations to remote sites.
}

// Initializes the ABTUInstance abtu.
func Init(id SiteId, sv Timestamp, h []Operation, rb []Operation) *ABTUInstance {
	var abtu *ABTUInstance = &ABTUInstance{}
	abtu.id = id
	abtu.sv = sv

	abtu.h = DeepCopyOperations(h)
	abtu.localTimestampHistory = make([]Timestamp, 0)

	abtu.rbm.Start(rb)

	abtu.localToABTU = make(chan []byte, 20)
	abtu.aBTUToLocal = make(chan []byte, 20)

	abtu.remoteToABTU = make(chan []byte, 20)
	abtu.aBTUToRemote = make(chan []byte, 20)

	return abtu
}

func (abtu *ABTUInstance) Run() (chan<- []byte, <-chan []byte, chan<- []byte, <-chan []byte){
	go abtu.run()
	return abtu.localToABTU, abtu.aBTUToLocal, abtu.remoteToABTU, abtu.aBTUToRemote
}

func (abtu *ABTUInstance) Stop(){
	//abtu.manager <- "stop"
}

func (abtu *ABTUInstance) run() {
	go abtu.listenToRemote()

	go abtu.launchController()
}

func (abtu *ABTUInstance) listenToRemote() {
	more := true
	var bytes []byte

	for ;more ; {
		select {
		case bytes, more = <- abtu.remoteToABTU:
			if more {
				remoteOp, err := DecodeFromPeers(bytes)
				if err!=nil {
					log.Fatal(err)
				}

				ack := make(chan bool)
				abtu.rbm.Add <- AddOp{remoteOp, ack}
				<- ack
			} else {
				close(abtu.rbm.Add)
			}
		}
	}
}

func (abtu *ABTUInstance) launchController() {
	for {
		select {
		case bytes :=  <- abtu.localToABTU:
			abtu.handleFrontendMessage(bytes)
		default:
			// Channel with buffer of size 1, in order not to lock rbm.
			causallyReadyOperationChannel := make(chan Operation, 1)
			abtu.rbm.Get <- GetCausallyReadyOp{abtu.sv, causallyReadyOperationChannel}

			select {
			case causallyReadyOp := <- causallyReadyOperationChannel:
				abtu.handleCausallyReadyOperation(causallyReadyOp)
			case bytes := <- abtu.localToABTU:
				abtu.handleFrontendMessage(bytes)
			}
		}
	}
}

func (abtu *ABTUInstance) handleFrontendMessage(bytes []byte) {
	var frontendMsg FrontendMessage
	err := json.Unmarshal(bytes, &frontendMsg)
	if err != nil {
		log.Fatal(errors.New("Could not decode frontendMessage:" + err.Error()))
	}

	switch frontendMsg.Type {
	case LocalOp:
		abtu.handleLocalOperation(frontendMsg.Content)
	case Undo:
		abtu.handleLocalUndo(frontendMsg.Content)
	default:
		log.Fatal(errors.New("Unexpected message from frontend: " + frontendMsg.Type))
	}
}

func (abtu *ABTUInstance) handleLocalOperation(bytes []byte) {
	localOp, err := DecodeFromFrontend(bytes, abtu.id)
	if err != nil {
		log.Fatal(errors.New("Cannot decode local operation: " + err.Error()))
	}

	toDispatchOp := abtu.localThread(localOp)

	ackToSendFrontend, err := json.Marshal(FrontendMessage{AckLocalOp, []byte{}})
	if err != nil {
		log.Fatal(errors.New("Could not send ackLocalOperation, Json encoding failed :" + err.Error()))
	}

	// Execute locally
	abtu.aBTUToLocal <- ackToSendFrontend

	bytesToDispatch, err := toDispatchOp.EncodeToPeers()
	if err != nil {
		log.Fatal(errors.New("Cannot send operation to aBTUToRemote:" + err.Error()))
	}

	abtu.aBTUToRemote <- bytesToDispatch
}

func (abtu *ABTUInstance) handleLocalUndo(bytes []byte) {
	var toUndo int32
	err := json.Unmarshal(bytes, &toUndo)
	if err != nil {
		log.Fatal(errors.New("Cannot decode local undo: " + err.Error()))
	}

	toExecuteOp := abtu.localThreadUndo(toUndo)


	if toExecuteOp.OpType() == UNIT {
		nackLocalUndo, err := json.Marshal(FrontendMessage{NackLocalUndo, []byte{}})
		if err != nil {
			log.Fatal(errors.New("Cannot encode nacklocalundo: " + err.Error()))
		}

		abtu.aBTUToLocal <- nackLocalUndo
	} else {
		undoFrontendOp, err := json.Marshal(OperationToFrontendOperation(toExecuteOp))
		if err != nil {
			log.Fatal(errors.New("Could not encode simple operation:" + err.Error()))
		}

		ackLocalUndo, err := json.Marshal(FrontendMessage{AckLocalUndo, undoFrontendOp})
		if err != nil {
			log.Fatal(errors.New("Could not send ackLocalUndo, Json encoding failed :" + err.Error()))
		}

		// Execute locally
		abtu.aBTUToLocal <- ackLocalUndo

		toDispatchOp, err := toExecuteOp.EncodeToPeers()
		if err != nil {
			log.Fatal(errors.New("Cannot send operation to aBTUToRemote:" + err.Error()))
		}

		// Dispatch to other peers
		abtu.aBTUToRemote <- toDispatchOp
	}
}

func (abtu *ABTUInstance) handleCausallyReadyOperation(causallyReadyOp Operation){
	toExecuteOp, h, sv := abtu.remoteThread(causallyReadyOp)

	if toExecuteOp.OpType() != UNIT {
		frontendOp, err := json.Marshal(OperationToFrontendOperation(toExecuteOp))
		if err != nil {
			log.Fatal(errors.New("Could not encode simple operation:" + err.Error()))
		}

		remoteOperation, err := json.Marshal(FrontendMessage{RemoteOp, frontendOp})
		if err != nil {
			log.Fatal(errors.New("Could not send remoteOperation, Json encoding failed :" + err.Error()))
		}

		// Execute locally
		abtu.aBTUToLocal <- remoteOperation

		frontendBytes := <- abtu.localToABTU

		var frontendMsg FrontendMessage
		err = json.Unmarshal(frontendBytes, &frontendMsg)
		if err != nil {
			log.Fatal(errors.New("Could not decode frontendMessage:" + err.Error()))
		}

		switch frontendMsg.Type {
		case LocalOp:
			abtu.handleLocalOperation(frontendMsg.Content)
		case Undo:
			abtu.handleLocalUndo(frontendMsg.Content)
		case AckRemoteOp:
			abtu.sv = sv
			abtu.localTimestampHistory = append(abtu.localTimestampHistory, sv)
			abtu.h = h

			ack := make(chan bool)
			abtu.rbm.RemoveRearrange <- RemoveRearrangeOp{ack}
			<- ack
		}
	} else { // Operation is unit, we do not need to apply it, but we need to integrate effects in H and SV
		abtu.sv = sv
		abtu.h = h

		ack := make(chan bool)
		abtu.rbm.RemoveRearrange <- RemoveRearrangeOp{ack}
		<- ack
	}
}