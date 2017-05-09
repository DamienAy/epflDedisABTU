package remoteBufferManager

import . "github.com/DamienAy/epflDedisABTU/operation"

type GetOp struct {
	Ret chan []Operation
}

type RemoveRearrangeOp struct {
	Index int
	Ret chan bool
}

type AddOp struct {
	Operation Operation
	Ret chan bool
}

type RemoteBufferManager struct {
	Add chan AddOp
	Get chan GetOp
	RemoveRearrange chan RemoveRearrangeOp

	rb []Operation
}

func (rbm *RemoteBufferManager) Start(rb []Operation){
	rbm.Add = make(chan AddOp)
	rbm.Get = make(chan GetOp)
	rbm.RemoveRearrange = make(chan RemoveRearrangeOp)
	//rbm.Remove = make(chan RemoveOp)

	rbm.rb = make([]Operation, len(rb))
	copy(rbm.rb, rb)

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
					rbm.rb = append(rbm.rb, addOp.Operation)
					addOp.Ret <- true
				}
			// Return a copy of rb
			case getOp := <- rbm.Get:
				rb = make([]Operation, len(rbm.rb))
				copy(rb, rbm.rb)

				getOp.Ret <- rb
			case removeRearrangeOp := <- rbm.RemoveRearrange: //Remove entry from rb only if the Delete flag is set to true
				if removeRearrangeOp.Index>=len(rbm.rb) { // If one item has already been removed, discard any future removes.
					removeRearrangeOp.Ret <- false
				} else {
					newBuffer := make([]Operation, len(rbm.rb)-1)

					for i, el := range rbm.rb {
						if i<removeRearrangeOp.Index{
							newBuffer[i] = el
						} else if i>removeRearrangeOp.Index{
							newBuffer[i-1] = el
						}
					}

					rbm.rb = newBuffer

					removeRearrangeOp.Ret <- true
				}
			}
		}
	}()


}


