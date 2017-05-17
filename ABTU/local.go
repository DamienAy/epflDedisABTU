package ABTU

import (
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	"errors"
	enc "github.com/DamienAy/epflDedisABTU/ABTU/encoding"
	"encoding/json"
)


func (abtu *ABTUInstance) LocalThread(bytes []byte) error {

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

	// Execute locally
	abtu.lOut <- ackToSendFrontend

	operationToDispatch := abtu.IntegrateL(localOp)

	bytesToDispatch, err := operationToDispatch.EncodePeers()
	if err != nil {
		return errors.New("Cannot send operation to rOut:" + err.Error())
	}

	abtu.rOut <- bytesToDispatch

	return nil
}

func (abtu *ABTUInstance) LocalThreadUndo(toUndo uint64) error {
	// Need to find undo op in H
	toUndoOp := &abtu.h[toUndo]

	if toUndoOp.Uv().Size()!= 0 || len(toUndoOp.Dv())!=0 {
		// If operation has allready been undone or some other operation is dependent on this one.
		return errors.New("Operation cannot be undone.")
	} else {
		undoOp, err := toUndoOp.GetInverse(abtu.id)
		if err != nil {
			return err
		}

		abtu.sv.Increment(abtu.id)

		undoOp.AddV(abtu.sv)
		undoOp.AddAllOv(toUndoOp.V())
		//Need to find undo op in H
		toUndoOp.SetUv(abtu.sv)

		UndoFrontendOperation, err := json.Marshal(OperationToFrontendOperation(undoOp))
		if err != nil {
			return errors.New("Could not encode simple operation:" + err.Error())
		}

		ackLocalUndo, err := json.Marshal(enc.FrontendMessage{enc.AckLocalUndo, UndoFrontendOperation})
		if err != nil {
			return errors.New("Could not send ackLocalUndo, Json encoding failed :" + err.Error())
		}

		// Execute locally
		abtu.lOut <- ackLocalUndo

		operationToDispatch := abtu.IntegrateL(undoOp)

		bytesToDispatch, err := operationToDispatch.EncodePeers()
		if err != nil {
			return errors.New("Cannot send operation to rOut:" + err.Error())
		}

		abtu.rOut <- bytesToDispatch
	}

	return nil
}

func (abtu *ABTUInstance) IntegrateL(toIntegrateOp Operation) Operation{
	localOp := DeepCopyOperation(toIntegrateOp)

	k := len(abtu.h)

	var offset Position
	if localOp.OpType() == INS {
		offset = 1
	} else {
		offset = -1
	}

	if localOp.Ov() == nil { // o is a normal operation
		for i:= len(abtu.h)-1 ; i>=0; i-- {
			if abtu.h[i].IsGreaterH(localOp) {
				k = i
				abtu.h[i].SetPos(abtu.h[i].Pos()+offset)
			} else if abtu.h[i].IsSmallerH(localOp) {
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
		for j := range abtu.h {
			if IntersectionIsNotEmpty(localOp.Ov(), abtu.h[j].V()) {
				i = j
				break
			}
		}

		localOp.AddAllTv(abtu.h[i].V())
		k = i + 1

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
			newH[index] = localOp
		} else {
			newH[index] = abtu.h[index-1]
		}
	}

	abtu.h = newH

	return DeepCopyOperation(localOp)
}
