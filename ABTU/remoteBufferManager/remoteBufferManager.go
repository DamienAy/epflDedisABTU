package remoteBufferManager

import . "github.com/DamienAy/epflDedisABTU/ABTU/operation"
import (
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
)

type GetCausallyReadyOp struct {
	currentTime Timestamp
	Ret chan Operation
}

type RemoveRearrangeOp struct {
	Ret chan bool
}

type AddOp struct {
	Operation Operation
	Ret chan bool
}

type RemoteBufferManager struct {
	Add chan AddOp
	Get chan GetCausallyReadyOp
	RemoveRearrange chan RemoveRearrangeOp

	rb []Operation
	siteId SiteId

	ABTUIsWaitingCausallyReadyOp bool
	ABTUSV Timestamp
	CausallyReadyOpRetChan chan Operation
	currentCausallyReadyOperationIndex int
}

func (rbm *RemoteBufferManager) Start(rb []Operation, siteId SiteId){
	rbm.Add = make(chan AddOp)
	rbm.Get = make(chan GetCausallyReadyOp)
	rbm.RemoveRearrange = make(chan RemoveRearrangeOp)

	rbm.rb = DeepCopyOperations(rb)
	rbm.siteId = siteId

	rbm.ABTUIsWaitingCausallyReadyOp = false
	rbm.currentCausallyReadyOperationIndex = -1;

	go func () {
		cont := true // True as long as the Add channel is not closed.

		for ; cont ; {
			select {
			// As long as the Add channel is not closed, put the operation in the receiving buffer.
			case addOp, notDone := <- rbm.Add:
				if !notDone {
					cont = false
					addOp.Ret <-false
				} else {
					rbm.rb = append(rbm.rb, DeepCopyOperation(addOp.Operation))
					addOp.Ret <- true

					// If ABTU is waiting for a causally ready operation, check againg.
					if rbm.ABTUIsWaitingCausallyReadyOp {
						causallyReadyOp, index := rbm.getFirstCausallyReadyOperation()
						if index >= 0 {
							rbm.currentCausallyReadyOperationIndex = index
							rbm.CausallyReadyOpRetChan <- causallyReadyOp
							rbm.ABTUIsWaitingCausallyReadyOp = false
						}
					}
				}
			// Return a copy of rb
			case getCausallyReadyOp := <- rbm.Get:
				rbm.ABTUSV = getCausallyReadyOp.currentTime
				causallyReadyOp, index := rbm.getFirstCausallyReadyOperation()

				if index >= 0 {
					rbm.currentCausallyReadyOperationIndex = index
					rbm.CausallyReadyOpRetChan <- causallyReadyOp
					rbm.ABTUIsWaitingCausallyReadyOp = false
				} else {
					rbm.ABTUIsWaitingCausallyReadyOp = true
				}

			case removeRearrangeOp := <- rbm.RemoveRearrange:
				if rbm.currentCausallyReadyOperationIndex>=len(rbm.rb) { // If one item has already been removed, discard any future removes.
					removeRearrangeOp.Ret <- false
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

					rbm.ABTUIsWaitingCausallyReadyOp = false

					removeRearrangeOp.Ret <- true
				}
			}
		}
	}()


}

// Returns the first causally ready operation from the receiving buffer rb.
// If there is no causally ready operation yet, it returns a unit operation.
// Does not make any changes to currentTime and rb
func (rbm *RemoteBufferManager) getFirstCausallyReadyOperation() (Operation, int){
	for i:=0; i<len(rbm.rb) ; i++ {
		for _, operationTimestamp := range rbm.rb[i].V() {
			if operationTimestamp.IsCausallyReady(rbm.ABTUSV, rbm.siteId) {
				return DeepCopyOperation(rbm.rb[i]), i
			}
		}
	}
	return PartialOperation(rbm.siteId, UNIT, 0, 0), -1
}


