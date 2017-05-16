package ABTU

import (
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	"errors"
	enc "github.com/DamienAy/epflDedisABTU/ABTU/encoding"
	"encoding/json"
)


func (abtu *ABTUInstance) LocalThread(bytes []byte, undo bool, toUndo uint64) error {

	if !undo {
		localOp, err := DecodeFrontend(bytes, abtu.id)

		if err != nil {
			return errors.New("Cannot decode local operation: " + err.Error())
		}

		abtu.sv.Increment(abtu.id)
		localOp.AddV(abtu.sv)

		ackToSendFrontend, err := json.Marshal(enc.FrontendMessage{enc.AckLocalOperation, []byte{}})
		if err != nil {
			return errors.New("Could not send ackLocalOperation, Json encoding failed :" + err.Error())
		}

		abtu.lOut <- ackToSendFrontend

		abtu.IntegrateL(&localOp)

		operationToSend, err := localOp.EncodePeers()
		if err != nil {
			return errors.New("Cannot send operation to rOut:" + err.Error())
		}

		abtu.rOut <- operationToSend
	} else {
		toUndoOp := &abtu.h[toUndo]

		if toUndoOp.Uv().Size()!= 0 || len(toUndoOp.Dv())!=0 {
			return nil
		} else {
			undoOp, err := toUndoOp.GetInverse(abtu.id); check(err)
			abtu.sv.Increment(abtu.id)
			undoOp.AddV(abtu.sv)
			undoOp.SetOv(&(toUndo.V()[0])) // Right to take the first one???
			toUndoOp.SetUv(&SV)
			Execute(o)
			o2 := IntegrateL(o)
			communicationService.Send(o2)
		}
	}

	return nil

}

func (abtu *ABTUInstance) IntegrateL(localOp *Operation) {

	k := len(abtu.h)

	var offset Position
	if localOp.OpType() == INS {
		offset = 1
	} else {
		offset = -1
	}

	if localOp.Tv() == nil { // o is a normal operation
		for i:= len(abtu.h)-1 ; i>=0; i-- {
			if abtu.h[i].IsGreaterH(*localOp) {
				k = i
				abtu.h[i].SetPos(abtu.h[i].Pos()+offset)
			} else if abtu.h[i].IsSmallerH(*localOp) {
				break
			} else {
				localOp.AddAllTv(abtu.h[i].V())
				if localOp.OpType() == DEL {
					abtu.h[i].AddAllDv(localOp.V())
				}
			}
		}
	} else {
		var i int
		for index, h := range abtu.h {
			if localOp.Ov().IsContainedIn(h.V()) {
				i = index
				break
			}
		}

		k = i + 1
		localOp.AddAllTv(abtu.h[i].V())

		for j:=k; j<=len(abtu.h); j++ {
			abtu.h[j].SetPos(abtu.h[j].Pos()+offset)
		}
	}

	// Insert localOp in H at position i
	newH := make([]Operation, len(abtu.h)+1)

	for index := range newH {
		if index < k {
			newH[index] = abtu.h[index]
		} else if index == k {
			newH[index] = DeepCopyOperation(*localOp)
		} else {
			newH[index] = abtu.h[index-1]
		}
	}

	abtu.h = newH
}

func Execute(o Operation)  {
}


func Dispatch(o Operation) {

}
