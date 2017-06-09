package receivingBufferManager

import (
	"testing"
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	"fmt"
)

func TestSimpleRequests ( *testing.T) {
	var rbm ReceivingBufferManager
	rbm.Start(make([]Operation,0))

	sv := NewTimestamp(3)
	sv.Increment(1)

	var char Char = make(Char, 1)
	char[0] = 'a'
	remoteOperation := PartialOperation(1, INS, 0, char);

	sv2 := DeepCopyTimestamp(sv)
	sv2.Increment(1)
	remoteOperation.AddV(sv2)

	// Adding.
	answer := make(chan bool)
	rbm.Add <- AddOp{remoteOperation, answer}
	<-answer

	// Get causally ready operation.
	ret := make(chan Operation)
	rbm.Get <- GetCausallyReadyOp{sv, ret}
	fmt.Println(<-ret)

	// Rearange.
	rbm.RemoveRearrange <- RemoveRearrangeOp{answer }
	<-answer
}


