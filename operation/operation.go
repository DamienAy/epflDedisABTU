package operation

import (
	. "github.com/DamienAy/epflDedisABTU/singleTypes";
	. "github.com/DamienAy/epflDedisABTU/timestamp";
	"errors"
)

type Operation struct {
	id SiteId
	opType OpType
	position Position
	character Char
	v []Timestamp
	dv []Timestamp
	tv []Timestamp
	ov Timestamp
	uv Timestamp
}

// Returns the siteId where the operation o has been generated.
func (o *Operation) GetId() SiteId {
	return o.id
}

// Returns the OpType of the operation o.
func (o *Operation) GetOpType() OpType {
	return o.opType
}

// Returns the Position of the operation o.
func (o *Operation) GetPos() Position {
	return o.position
}

//Sets the positition of the operation o to the value p.
func (o *Operation) SetPos(p Position) {
	o.position = p
}

// Returns the character of the operation o.
func (o *Operation) GetChar() Char {
	return o.character
}


// Returns a slice containing the timestamps of operation o.
func (o *Operation) GetV() []Timestamp {
	return o.v
}

// Appends the timestamp t to the timestamps slice of o.
func (o *Operation) AddV(t Timestamp) {
	o.v = append(o.v, t)
}

// Returns a slice containing the timestamps of operations that depend on operation o.
func (o *Operation) GetDv() []Timestamp {
	return o.dv
}

// Appends the Timestamp t to the timestamps slice of operations that depend on operation o.
func (o *Operation) AddDv(t Timestamp) {
	o.dv = append(o.dv, t)
}

// Returns a slice containing the timestamps of operations whose effect objects tie with o.c.
func (o *Operation) GetTv() []Timestamp {
	return o.tv
}

// Appends the timestamp t to the timestamps slice of operations whose effect objects tie with o.c.
func (o *Operation) AddTv(t Timestamp) {
	o.tv = append(o.tv, t)
}

// Returns the timestamp of the original operation o undoes (if operation o is an undo, otherwise nil).
func (o *Operation) GetOv() Timestamp {
	return o.ov
}

// Sets the timestamp ov of the operation o to t.
func (o *Operation) SetOv(t Timestamp) {
	o.ov = t
}

// Returns the timestamp of the operation that undoes o.
func (o *Operation) GetUv() Timestamp {
	return o.ov
}

// Sets the timestamp uv of the operation o to t.
func (o *Operation) SetUv(t Timestamp) {
	o.ov = t
}

// Returns true if and only if operation o1 happened before operation o2.
func (o1 *Operation) HapenedBefore(o2 Operation) bool {
	for _, e1:= range o1.GetV(){
		for _, e2:= range o2.GetV() {
			if e1.HappenedBefore(e2) {
				return true
			}
		}
	}

	return false
}

// Returns true if and only if operation o1 is concurrent with operation o2.
func (o1 *Operation) isConcurrentWith(o2 Operation) bool {
	return !(o1.HapenedBefore(o2) || o2.HapenedBefore(o1))
}

// Returns true if and only if operation o1 is smaller in effect relation order than o2.
func (o1 *Operation) IsSmallerInEffectRelationOrder(o2 Operation) bool {
	p1 := o1.position < o2.position
	p2 := o1.position==o2.position && o1.opType==INS && o2.opType==DEL
	p3 := o1.position==o2.position && o1.opType==INS && o2.opType==INS && o1.id < o2.id

	return p1 || p2 || p3
}

// Returns true if and only if operation o1 is greater in effect relation order than o2.
func (o1 *Operation) IsGreaterInEffectRelationOrder(o2 Operation) bool {
	p1 := o1.position > o2.position
	p2 := o1.position==o2.position && o1.opType==DEL && o2.opType==INS
	p3 := o1.position==o2.position && o1.opType==INS && o2.opType==INS && o1.id > o2.id

	return p1 || p2 || p3
}

// Returns the inverse of the operation o. An error is returned if the operation o is unit.
func (o *Operation) inverse() (Operation, error) {
	if o.opType == INS {
		inverse := Operation{opType:DEL, position:o.position, character:o.character}
		return inverse, nil

	} else if o.opType == DEL{
		inverse := Operation{opType:INS, position:o.position, character:o.character}
		return inverse, nil
	} else {
		return nil, errors.New("Computing the inverse of a unit operation.")
	}

}

