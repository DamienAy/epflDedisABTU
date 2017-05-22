package ABTU

import (
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
)

// Executes remotethread algorithm with remoteOperation and returns the operation to execute locally.
// Returns the resulting history buffer and local timestamp SV without affecting abtu.h and abtu.sv
func (abtu *ABTUInstance) RemoteThread(remoteOperation Operation) (Operation, []Operation, Timestamp){
	remoteOp := DeepCopyOperation(remoteOperation)
	H := DeepCopyOperations(abtu.h)
	SV := DeepCopyTimestamp(abtu.sv)

	if remoteOp.Ov() != nil { // remoteOp is undo.
		var i int
		for k:=0; k<len(H); k++ {
			if IntersectionIsNotEmpty(H[k].V(), remoteOp.Ov()) {
				i = k
				break
			}
		}

		if H[i].Uv()!=nil { // H[i] has already been undone => merge two operations
			var j int
			for k:=0; k<len(H); k++ {
				if IntersectionIsNotEmpty(H[k].V(), remoteOp.Uv()) {
					i = k
					break
				}
			}

			H[j].AddAllV(remoteOp.V())
			SV.Increment(remoteOp.Id())

			return PartialOperation(abtu.id, UNIT, 0, 0), H, SV

		} else {
			H[i].AddAllUv(remoteOp.Uv())
			remoteOp.AddAllOv(H[i].V())
		}
	}

	toExecuteLocallyOp, H := IntegrateR(remoteOp, H)

	SV.Increment(remoteOp.Id())

	return toExecuteLocallyOp, H, SV
}

// Executes integrateR algortithm, does not modify toIntegrateRemoteOp
// Returns the operation to integrate (only if type!=UNIT) and the updated history buffer.
func IntegrateR(toIntegrateRemoteOperation Operation, historyBuffer []Operation) (Operation, []Operation) {
	toIntegrateRemoteOp := DeepCopyOperation(toIntegrateRemoteOperation)
	H := historyBuffer


	k := len(H)

	if toIntegrateRemoteOp.Tv()==nil {
		for i:=0; i<len(H); i++ {
			if H[i].IsConcurrentWith(toIntegrateRemoteOp) {
				var offset Position
				if toIntegrateRemoteOp.OpType() == INS {
					offset = 1
				} else {
					offset = -1
				}

				if H[i].Tv()!=nil || H[i].IsSmallerC(toIntegrateRemoteOp) {
					toIntegrateRemoteOp.SetPos(toIntegrateRemoteOp.Pos()+offset)
				} else if H[i].IsGreaterC(toIntegrateRemoteOp) {
					k = i
					break
				} else { //H[i] == toIntegrateRemoteOp.
					toIntegrateRemoteOp.SetToUnit()
					H[i].AddAllV(toIntegrateRemoteOp.V())
					break
				}

			} else if H[i].IsGreaterC(toIntegrateRemoteOp) { // H[i].happenedBefore(toIntegrateRemoteOp) holds.
				k = i
				break
			}
		}
	} else { //toIntegrateRemoteOp.tv is not empty, also covers undo case
		var i int
		for j:=0; j<len(H); j++ {
			if IntersectionIsNotEmpty(H[k].V(), toIntegrateRemoteOp.Tv()) {
				i = j
				break
			}
		}

		toIntegrateRemoteOp.SetPos(H[i].Pos())
		k = i + 1

		//for ;  ; {//sort it out!!!
		for ; IntersectionIsNotEmpty(H[k].Tv(), H[i].V()); {
			if H[k].IsSmallerC(toIntegrateRemoteOp){
				var offset Position
				if H[k].OpType() == INS {
					offset = 1
				} else {
					offset = -1
				}

				toIntegrateRemoteOp.SetPos(toIntegrateRemoteOp.Pos()+offset)

			} else if H[k].IsGreaterC(toIntegrateRemoteOp) {
				break
			} else {
				toIntegrateRemoteOp.SetToUnit()
				H[k].AddAllV(toIntegrateRemoteOp.V())
				break
			}

			k++
		}
	}

	if toIntegrateRemoteOp.OpType() != UNIT {
		var offset Position
		if toIntegrateRemoteOp.OpType() == INS {
			offset = 1
		} else {
			offset = -1
		}

		for j:=k; j<=len(H); j++ {
			H[j].SetPos(H[j].Pos()+offset)
		}

		newH := make([]Operation, len(H)+1)

		for index := range newH {
			if index < k {
				newH[index] = H[index]
			} else if index == k {
				newH[index] = toIntegrateRemoteOp
			} else {
				newH[index] = H[index-1]
			}
		}

		return toIntegrateRemoteOp, newH
	}
}