package remoteBufferManagerTest

import (
	"testing"
	. "github.com/DamienAy/epflDedisABTU/operation"
	. "github.com/DamienAy/epflDedisABTU/remoteBufferManager"
	"fmt"
)

func TestJustTryIt ( *testing.T) {
	operation := PartialOperation(0, 0, 2, 0)

	var rbm RemoteBufferManager

	rbm.Start(make([]Operation,0 ))

	answer := make(chan bool)
	answer2 := make(chan bool)
	rbm.Add <- AddOp{operation, answer}
	fmt.Println(<-answer)
	rbm.Add <- AddOp{operation, answer2}
	fmt.Println(<-answer2)

	ret := make(chan []Operation)
	rbm.Get <- GetOp{ret}
	fmt.Println(<-ret)
	rbm.RemoveRearrange <- RemoveRearrangeOp{1, answer }
	fmt.Println(<-answer)


	rbm.Get <- GetOp{ret}
	rbm.RemoveRearrange <- RemoveRearrangeOp{0, answer }
	//fmt.Println(<-answer)


	/*rbm.RemoveRearrange <- RemoveRearrangeOp{0, answer }
	fmt.Println(<-answer)*/
	//close(rbm.Add)
	//close(answer)
}
