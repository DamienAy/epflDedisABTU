package ABTU

import (
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	"log"
)

// Executes local thread algorithm, does not make any changes to localOperation
func (abtu *ABTUInstance) LocalThread(localOperation Operation) Operation {
	localOp := DeepCopyOperation(localOperation)

	abtu.sv.Increment(abtu.id)
	abtu.localTimestampHistory = append(abtu.localTimestampHistory, abtu.sv)

	localOp.AddV(abtu.sv)

	return abtu.IntegrateL(localOp)
}

// Decodes local undo operation from toUndo, and executes local thread algorithm.
func (abtu *ABTUInstance) LocalThreadUndo(toUndo int32) Operation {
	toUndoTimestamp := abtu.localTimestampHistory[len(abtu.localTimestampHistory) - int(toUndo)]

	// Need to find undo op in H
	toUndoOp := UnitOperation(abtu.id)

	for i:=len(abtu.h)-1; i>=0 ; i-- {
		isContainedIn, err := toUndoTimestamp.IsContainedIn(abtu.h[i].V())
		if err != nil {
			log.Fatal(err)
		}

		if isContainedIn && abtu.h[i].Id()==abtu.id{
			toUndoOp = abtu.h[i]
		}
	}

	if toUndoOp.OpType() == UNIT {
		log.Fatal("Did not find the operation to undo.")
	}

	if len(toUndoOp.Uv()) != 0 || len(toUndoOp.Dv())!=0 {
		// If operation has allready been undone or some other operation is dependent on this one.
		return UnitOperation(abtu.id)
	} else {
		undoOp, err := toUndoOp.GetInverse(abtu.id)
		if err != nil {
			log.Fatal(err)
		}

		abtu.sv.Increment(abtu.id)

		undoOp.AddV(abtu.sv)
		undoOp.AddAllOv(toUndoOp.V())
		//Need to find undo op in H
		toUndoOp.AddUv(abtu.sv)

		return abtu.IntegrateL(undoOp)
	}
}

// Executes IntegrateL algorithm: integrate toIntegrateOp into the history buffer and updates all timestamps and positions.
// Does not make any changes on toIntegrateOp.
func (abtu *ABTUInstance) IntegrateL(toIntegrateOp Operation) Operation {
	localOp := DeepCopyOperation(toIntegrateOp)

	k := len(abtu.h)

	var offset Position
	if localOp.OpType() == INS {
		offset = 1
	} else {
		offset = -1
	}

	if len(localOp.Ov()) == 0 { // o is a normal operation
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
		var i int = -1
		for j := range abtu.h {
			intersectionIsNotEmpty, err := IntersectionIsNotEmpty(localOp.Ov(), abtu.h[j].V())
			if err != nil {
				log.Fatal(err)
			}
			if intersectionIsNotEmpty {
				i = j
				break
			}
		}

		if i == -1 {
			log.Fatal("There is no operation to undo.")
		}

		localOp.AddAllTv(abtu.h[i].V())
		k = i + 1

		for j:=k; j<len(abtu.h); j++ {
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
