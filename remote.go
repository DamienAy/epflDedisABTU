package main

import (
	. "github.com/DamienAy/epflDedisABTU/operation"
	. "github.com/DamienAy/epflDedisABTU/singleTypes"
)

func RemoteThread(o Operation) {
	lock.Lock()
	if o.GetOv() == nil {
		var i int
		for index, h := range H {
			if o.GetOv().IsContainedIn(h.GetV()) {
				i = index
				break
			}
		}

		if H[i].GetUv() != nil { // H[i] has already been undone => merge two operations
			var j int
			for index, h := range H {
				if o.GetUv().IsContainedIn(h.GetV()) {
					j = index
					break
				}
			}
			H[j].AddV(o.GetV()[0]) // okey to add only first?????
			SV[o.GetId()]++
			return
		} else {
			H[i].SetUv(&(o.GetV()[0]))
			//o.set//----------------------------------------------------

		}
	}
	o2 := IntegrateR(o)
	SV[o.GetPos()]++
	if o2.GetOpType()!= UNIT {
		Execute(o2)
	}

	lock.Unlock()
}


func IntegrateR(o Operation) Operation {
	k := len(H)

	if len(o.GetTv()) == 0 {
		for i, h := range H {
			if h.IsConcurrentWith(o) {
				var offset Position
				if o.GetOpType() == INS {
					offset = 1
				} else {
					offset = -1
				}

				if len(h.GetTv())!=0 || h.IsSmallerC(o) {
					o.SetPos(o.GetPos()+offset)
				} else if h.IsGreaterC(o) {
					k = i
					break
				} else { //h = o
					o.SetToUnit()
					h.AddV(o.GetV()[0]) // okey to add only first ??????????
					break
				}
			} else if h.IsGreaterC(o) {
				k = i
				break
			}
		}
	} else { //o.tv is not empty, also covers undo case
		var i int
		for index := range H {
			if true { //////o.GetTv().IsContainedIn(h.GetV())
				i = index
				break
			}
		}
		o.SetPos(H[i].GetPos())
		k = i + 1

		//for ;  ; {//sort it out!!!
		{
			if H[k].IsSmallerC(o){
				var offset Position
				if H[k].GetOpType() == INS {
					offset = 1
				} else {
					offset = -1
				}
				o.SetPos(o.GetPos()+offset)
			} else if H[k].IsGreaterC(o) {
				//break
			} else {
				o.SetToUnit()
				H[k].AddV(o.GetV()[0]) //okey to add only first???
				//break
			}
			k++
		}
	}

	if o.GetOpType() != UNIT {
		var offset Position
		if o.GetOpType() == INS {
			offset = 1
		} else {
			offset = -1
		}

		for j:=k; j<=len(H); j++ {
			H[j].SetPos(H[j].GetPos()+offset)
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
}