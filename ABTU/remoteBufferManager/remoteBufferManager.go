package remoteBufferManager

import . "github.com/DamienAy/epflDedisABTU/ABTU/operation"
import (
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	"log"
)

type GetCausallyReadyOp struct {
	CurrentTime Timestamp
	Return      chan Operation
}

type RemoveRearrangeOp struct {
	Ack chan bool
}

type AddOp struct {
	Operation Operation
	Ack chan bool
}

type RemoteBufferManager struct {
	Add chan AddOp
	Get chan GetCausallyReadyOp
	RemoveRearrange chan RemoveRearrangeOp

	rb []Operation
	siteId SiteId

	aBTUIsWaitingCausallyReadyOp bool
	aBTUSV Timestamp
	causallyReadyOpRetChan chan Operation
	currentCausallyReadyOperationIndex int
}

func (rbm *RemoteBufferManager) Start(rb []Operation, siteId SiteId){
	rbm.Add = make(chan AddOp)
	rbm.Get = make(chan GetCausallyReadyOp)
	rbm.RemoveRearrange = make(chan RemoveRearrangeOp)

	rbm.rb = DeepCopyOperations(rb)
	rbm.siteId = siteId

	rbm.aBTUIsWaitingCausallyReadyOp = false
	rbm.currentCausallyReadyOperationIndex = -1;

	go func () {
		cont := true // True as long as the Add channel is not closed.

		for ; cont ; {
			select {
			// As long as the Add channel is not closed, put the operation in the receiving buffer.
			case addOp, notDone := <- rbm.Add:
				if !notDone {
					cont = false
					//addOp.Ack <-false
				} else {
					rbm.rb = append(rbm.rb, DeepCopyOperation(addOp.Operation))
					addOp.Ack <- true

					// If ABTU is waiting for a causally ready operation, check againg.
					if rbm.aBTUIsWaitingCausallyReadyOp {
						causallyReadyOp, index := rbm.getFirstCausallyReadyOperation()
						if index >= 0 {
							rbm.currentCausallyReadyOperationIndex = index
							select {
							case rbm.causallyReadyOpRetChan <- causallyReadyOp:
							default:
							}
							rbm.aBTUIsWaitingCausallyReadyOp = false
						}
					}
				}
			// Return the first causally ready operation if awailable. DeepCopy the timestamp
			case getCausallyReadyOp := <- rbm.Get:
				rbm.aBTUSV = DeepCopyTimestamp(getCausallyReadyOp.CurrentTime)
				causallyReadyOp, index := rbm.getFirstCausallyReadyOperation()

				if index >= 0 {
					rbm.currentCausallyReadyOperationIndex = index
					select {
					case rbm.causallyReadyOpRetChan <- causallyReadyOp:
					default:
					}
					rbm.aBTUIsWaitingCausallyReadyOp = false
				} else {
					rbm.aBTUIsWaitingCausallyReadyOp = true
				}


			case removeRearrangeOp := <- rbm.RemoveRearrange:
				if rbm.currentCausallyReadyOperationIndex>=len(rbm.rb) { // If one item has already been removed, discard any future removes.
					removeRearrangeOp.Ack <- false
				} else {
					newBuffer := make([]Operation, len(rbm.rb)-1)

					for i, el := range rbm.rb {
						if i<rbm.currentCausallyReadyOperationIndex {
							newBuffer[i] = el
						} else if i>rbm.currentCausallyReadyOperationIndex {
							newBuffer[i-1] = el
						}
					}

					rbm.rb = newBuffer

					rbm.aBTUIsWaitingCausallyReadyOp = false

					removeRearrangeOp.Ack <- true
				}
			}
		}
	}()


}

// Returns the first causally ready operation from the receiving buffer rb.
// If there is no causally ready operation yet, it returns a unit operation.
// Does not make any changes to CurrentTime and rb
func (rbm *RemoteBufferManager) getFirstCausallyReadyOperation() (Operation, int){
	for i:=0; i<len(rbm.rb) ; i++ {
		for _, operationTimestamp := range rbm.rb[i].V() {
			isCausallyReady, err := operationTimestamp.IsCausallyReady(rbm.aBTUSV, rbm.siteId)

			if err != nil {
				log.Fatal(err)
			}

			if  isCausallyReady{
				return DeepCopyOperation(rbm.rb[i]), i
			}
		}
	}
	return UnitOperation(0), -1
}


