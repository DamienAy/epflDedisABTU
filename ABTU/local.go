package ABTU

import (
	. "github.com/DamienAy/epflDedisABTU/operation"
	"sync"
	. "github.com/DamienAy/epflDedisABTU/singleTypes"
	. "github.com/DamienAy/epflDedisABTU/timestamp"
	com "github.com/DamienAy/epflDedisABTU/communication"
	"log"
)

func check(err error){
	if err != nil {
		log.Fatal(err)
	}
}

func LocalThread(localOp SimpleOperation, lock *sync.Mutex, id SiteId, SV *Timestamp, H []Operation, comService com.CommunicationService, lastOp int) {
//----------------------------------------------
//GET OUT THE WAY
	lock.Lock()
//----------------------------------------------

	if localOp.OpType != UNIT { //normal operation case
		o := localOp.GetOperation()
		SV[id]++
		o.AddV(*SV)

		Execute(o)
		o2 := IntegrateL(o)
		comService.Send(o2)
	} else { //undo case
		toUndo := &(H[lastOp])
		if len(toUndo.GetUv())!= 0 || len(toUndo.GetDv())!=0 {
			return
		} else {
			o, err := toUndo.GetInverse(); check(err)
			SV[id]++
			o.AddV(*SV)
			o.SetOv(toUndo.GetV()[0]) // Right to take the first one???
			toUndo.SetUv(*SV)
			Execute(o)
			o2 := IntegrateL(o)
			comService.Send(o2)
		}

	}
//----------------------------------------------
//Free
	lock.Unlock()
//----------------------------------------------



}

func IntegrateL(o Operation) Operation {
	return o
}

func Execute(o Operation)  {
}