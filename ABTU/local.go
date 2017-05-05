package ABTU

import (
	. "github.com/DamienAy/epflDedisABTU/operation"
	. "github.com/DamienAy/epflDedisABTU/singleTypes"
)


func LocalThread(localOp SimpleOperation) {
	//----------------------------------------------
	//GET OUT THE WAY
	lock.Lock()
	//----------------------------------------------

	if localOp.OpType != UNIT { //normal operation case
		o := localOp.GetOperation()
		SV[ID]++
		o.AddV(SV)

		Execute(o)
		o2 := IntegrateL(o)
		communicationService.Send(o2)
	} else { //undo case
		toUndo := &(H[lastOp])
		if len(*(toUndo.GetUv()))!= 0 || len(toUndo.GetDv())!=0 {
			return
		} else {
			o, err := toUndo.GetInverse(); check(err)
			SV[ID]++
			o.AddV(SV)
			o.SetOv(&(toUndo.GetV()[0])) // Right to take the first one???
			toUndo.SetUv(&SV)
			Execute(o)
			o2 := IntegrateL(o)
			communicationService.Send(o2)
		}

	}
	//----------------------------------------------
	//Free
	lock.Unlock()
	//----------------------------------------------

}

func IntegrateL(o Operation) Operation {
	k := len(H)

	var offset Position
	if o.GetOpType() == INS {
		offset = 1
	} else {
		offset = -1
	}

	if o.GetTv() == nil { // o is a normal operation
		for i:= len(H)-1 ; i>=0; i-- {
			if H[i].IsGreaterH(o) {
				k = i
				H[i].SetPos(H[i].GetPos()+offset)
			} else if H[i].IsSmallerH(o) {
				break
			} else {
				o.AddTv(H[i].GetV()[0]) // ok to add only first????????
				if o.GetOpType() == DEL {
					H[i].AddDv(o.GetV()[0])
				}
			}
		}
	} else {
		var i int
		for index, h := range H {
			if o.GetOv().IsContainedIn(h.GetV()) {
				i = index
				break
			}
		}

		k = i + 1
		o.AddTv(H[i].GetV()[0]) //ok to add only first ?????

		for j:=k; j<=len(H); j++ {
			H[j].SetPos(H[j].GetPos()+offset)
		}
	}

	newH := make([]Operation, len(H)+1)

	for index := range newH {
		if index < k {
			newH[index] = H[index]
		} else if index == k {
			newH[index] = o
		} else {
			newH[index] = H[index-1]
		}
	}

	H = newH

	return o
}

func Execute(o Operation)  {
}
