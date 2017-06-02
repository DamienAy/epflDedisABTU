package ABTU

import (
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	"log"
)

// Executes remotethread algorithm with remoteOperation and returns the operation to execute locally.
// Returns the resulting history buffer and local timestamp SV without affecting abtu.h and abtu.sv
func (abtu *ABTUInstance) RemoteThread(remoteOperation Operation) (Operation, []Operation, Timestamp){
	remoteOp := DeepCopyOperation(remoteOperation)
	H := DeepCopyOperations(abtu.h)
	SV := DeepCopyTimestamp(abtu.sv)

	if len(remoteOp.Ov()) != 0 { // remoteOp is undo.
		var i int
		for k:=0; k<len(H); k++ {
			intersectionIsNotEmpty, err := IntersectionIsNotEmpty(H[k].V(), remoteOp.Ov())
			if err != nil {
				log.Fatal(err)
			}

			if intersectionIsNotEmpty {
				i = k
				break
			}
		}

		if len(H[i].Uv()) != 0 { // H[i] has already been undone => merge two operations
			var j int
			for k:=0; k<len(H); k++ {
				intersectionIsNotEmpty, err := IntersectionIsNotEmpty(H[k].V(), remoteOp.Uv())
				if err != nil {
					log.Fatal(err)
				}
				if intersectionIsNotEmpty {
					i = k
					break
				}
			}

			H[j].AddAllV(remoteOp.V())
			SV.Increment(remoteOp.Id())

			return UnitOperation(abtu.id), H, SV

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

	if len(toIntegrateRemoteOp.Tv()) == 0 {
		for i:=0; i<len(H); i++ {
			isConcurrentWith, err := H[i].IsConcurrentWith(toIntegrateRemoteOp)
			if err != nil {
				log.Fatal(err)
			}

			if isConcurrentWith {
				var offset Position
				if toIntegrateRemoteOp.OpType() == INS {
					offset = 1
				} else {
					offset = -1
				}

				if len(H[i].Tv())!=0 || H[i].IsSmallerC(toIntegrateRemoteOp) {
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
			intersectionIsNotEmpty, err := IntersectionIsNotEmpty(H[k].V(), toIntegrateRemoteOp.Tv())
			if err != nil {
				log.Fatal(err)
			}
			if intersectionIsNotEmpty {
				i = j
				break
			}
		}

		toIntegrateRemoteOp.SetPos(H[i].Pos())
		k = i + 1

		intersectionIsNotEmpty, err := IntersectionIsNotEmpty(H[k].Tv(), H[i].V())
		if err != nil {
			log.Fatal(err)
		}

		for ; intersectionIsNotEmpty; {
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

			intersectionIsNotEmpty, err = IntersectionIsNotEmpty(H[k].Tv(), H[i].V())
			if err != nil {
				log.Fatal(err)
			}
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

	return toIntegrateRemoteOp, H
}