package ABTU

import (
	. "github.com/DamienAy/epflDedisABTU/operation"
	. "github.com/DamienAy/epflDedisABTU/singleTypes"
)


func (abtu *ABTUInstance) LocalThread(localOp Operation) {

	if localOp.OpType() != UNIT { //normal operation case
		o := localOp
		abtu.sv.Increment(abtu.id)
		o.AddV(abtu.sv)

		Execute(o)
		o2 := IntegrateL(o)
		abtu.rOut <- o2.DeepCopy()
	} else { //undo case
		toUndo := &(H[lastOp])
		if len(toUndo.Uv())!= 0 || len(toUndo.GetDv())!=0 {
			return
		} else {
			o, err := toUndo.GetInverse(); check(err)
			SV[ID]++
			o.AddV(SV)
			o.SetOv(&(toUndo.V()[0])) // Right to take the first one???
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
	if o.OpType() == INS {
		offset = 1
	} else {
		offset = -1
	}

	if o.GetTv() == nil { // o is a normal operation
		for i:= len(H)-1 ; i>=0; i-- {
			if H[i].IsGreaterH(o) {
				k = i
				H[i].SetPos(H[i].Pos()+offset)
			} else if H[i].IsSmallerH(o) {
				break
			} else {
				o.AddTv(H[i].V()[0]) // ok to add only first????????
				if o.OpType() == DEL {
					H[i].AddDv(o.V()[0])
				}
			}
		}
	} else {
		var i int
		for index, h := range H {
			if o.Ov().IsContainedIn(h.V()) {
				i = index
				break
			}
		}

		k = i + 1
		o.AddTv(H[i].V()[0]) //ok to add only first ?????

		for j:=k; j<=len(H); j++ {
			H[j].SetPos(H[j].Pos()+offset)
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
