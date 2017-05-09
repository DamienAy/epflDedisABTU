package operation

import (
	. "github.com/DamienAy/epflDedisABTU/singleTypes";
	. "github.com/DamienAy/epflDedisABTU/timestamp";
	"errors"
)

type SimpleOperation struct {
	OpType OpType
	Character Char
	Position Position
}

func (simpleOp *SimpleOperation) GetOperation() Operation {
	return Operation{
		opType: simpleOp.OpType,
		character: simpleOp.Character,
		position: simpleOp.Position}
}

type Operation struct {
	id SiteId
	opType OpType
	position Position
	character Char
	v []Timestamp
	dv []Timestamp
	tv []Timestamp
	ov *Timestamp
	uv *Timestamp
}


func NewOperation(
	id SiteId,
	opType OpType,
	position Position,
	character Char,
	v []Timestamp,
	dv []Timestamp,
	tv []Timestamp,
	ov *Timestamp,
	uv *Timestamp) Operation {
	myOv := *ov
	myUv := *uv
	return Operation{id, opType, position, character, v, dv, tv, &myOv, &myUv}
}

// Returns a new operation.
// Only sets id, opType, position, and character.
func PartialOperation(
	id SiteId,
	opType OpType,
	position Position,
	character Char) Operation {
	return Operation{id: id, opType:DEL, position:position, character:character}
}


// Returns the siteId where the operation o has been generated.
func (o *Operation) Id() SiteId {
	return o.id
}

// Returns the OpType of the operation o.
func (o *Operation) OpType() OpType {
	return o.opType
}

// Sets the type of the operation o to UNIT
func (o *Operation) SetToUnit(){
	o.opType = UNIT
}

// Returns the Position of the operation o.
func (o *Operation) Pos() Position {
	return o.position
}

//Sets the positition of the operation o to the value p.
func (o *Operation) SetPos(p Position) {
	o.position = p
}

// Returns the character of the operation o.
func (o *Operation) Char() Char {
	return o.character
}


// Returns a copy of the slice containing the timestamps of operation o.
func (o *Operation) V() []Timestamp {
	v := make([]Timestamp, len(o.v))
	copy(v, o.v)
	return v
}

// Appends the timestamp t to the timestamps slice of o.
func (o *Operation) AddV(t Timestamp) {
	o.v = append(o.v, t)
}

// Returns a copy of the slice containing the timestamps of operations that depend on operation o.
func (o *Operation) GetDv() []Timestamp {
	dv := make([]Timestamp, len(o.dv))
	copy(dv, o.dv)
	return dv
}

// Appends the Timestamp t to the timestamps slice of operations that depend on operation o.
func (o *Operation) AddDv(t Timestamp) {
	o.dv = append(o.dv, t)
}

// Returns a slice containing the timestamps of operations whose effect objects tie with o.c.
func (o *Operation) GetTv() []Timestamp {
	tv := make([]Timestamp, len(o.tv))
	copy(tv, o.tv)
	return o.tv
}

// Appends the timestamp t to the timestamps slice of operations whose effect objects tie with o.c.
func (o *Operation) AddTv(t Timestamp) {
	o.tv = append(o.tv, t)
}

// Returns the timestamp of the original operation o undoes (if operation o is an undo, otherwise nil).
func (o *Operation) Ov() *Timestamp {
	ov := *o.ov
	return &ov
}

// Sets the timestamp ov of the operation o to t.
func (o *Operation) SetOv(t *Timestamp) {
	myOv := *t
	o.ov = &myOv
}

// Returns a copy of the timestamp of the operation that undoes o.
func (o *Operation) Uv() *Timestamp {
	uv := *o.uv
	return &uv
}

// Sets the timestamp uv of the operation o to t.
func (o *Operation) SetUv(t *Timestamp) {
	myOv := *t
	o.ov = &myOv
}

// Returns true if and only if operation o1 happened before operation o2.
func (o1 *Operation) HappenedBefore(o2 Operation) (bool, error) {
	for _, e1:= range o1.v{
		for _, e2:= range o2.v {
			happenendBefore, err := e1.HappenedBefore(e2)
			if err!=nil {
				return false, err
			} else {
				return happenendBefore, nil
			}
		}
	}

	return false, nil
}

// Returns true if and only if operation o1 is concurrent with operation o2.
func (o1 *Operation) IsConcurrentWith(o2 Operation) (bool, error) {
	o1Ho2, err := o1.HappenedBefore(o2)
	if err!=nil {
		return false, err
	}
	o2Ho1, err2 := o2.HappenedBefore(*o1)
	if err2!=nil {
		return false, err
	}

	return !(o1Ho2 || o2Ho1), nil
}

// Returns true if and only if operation o1 is smaller in effect relation order than o2
// when both operations have the same definition state.
func (o1 *Operation) IsSmallerC(o2 Operation) bool {
	p1 := o1.position < o2.position
	p2 := o1.position==o2.position && o1.opType==INS && o2.opType==DEL
	p3 := o1.position==o2.position && o1.opType==INS && o2.opType==INS && o1.id < o2.id

	return p1 || p2 || p3
}

// Returns true if and only if operation o1 is greater in effect relation order than o2
// when both operations have the same definition state.
func (o1 *Operation) IsGreaterC(o2 Operation) bool {
	p1 := o1.position > o2.position
	p2 := o1.position==o2.position && o1.opType==DEL && o2.opType==INS
	p3 := o1.position==o2.position && o1.opType==INS && o2.opType==INS && o1.id > o2.id

	return p1 || p2 || p3
}

// Returns true if and only if operation o1 is smaller in effect relation order than o2
// when dst(o2) = dst(o1) * o1
func (o1 *Operation) IsSmallerH(o2 Operation) bool {
	p1 := o1.position < o2.position
	p2 := o1.position==o2.position && o1.opType==DEL && o2.opType==DEL

	return p1 || p2
}

// Returns true if and only if operation o1 is greater in effect relation order than o2
// when dst(o2) = dst(o1) * o1
func (o1 *Operation) IsGreaterH(o2 Operation) bool {
	p1 := o1.position > o2.position
	p2 := o1.position==o2.position && o1.opType==INS && o2.opType==INS

	return p1 || p2
}

// Returns the inverse of the operation o. An error is returned if the operation o is unit.
// Only sets id, opType, position, and character.
func (o *Operation) GetInverse(siteId SiteId) (Operation, error) {
	if o.opType == INS {
		inverse := Operation{id: siteId, opType:DEL, position:o.position, character:o.character}
		return inverse, nil

	} else if o.opType == DEL{
		inverse := Operation{id: siteId, opType:INS, position:o.position, character:o.character}
		return inverse, nil
	} else {
		return *o, errors.New("Computing the inverse of a unit operation.")
	}

}

