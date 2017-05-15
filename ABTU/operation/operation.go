package operation

import (
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypesTypes";
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestampstamp";
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


func NewOperation(
	id SiteId,
	opType OpType,
	position Position,
	character Char,
	v []Timestamp,
	dv []Timestamp,
	tv []Timestamp,
	ov Timestamp,
	uv Timestamp) Operation {
	return Operation{
		id,
		opType,
		position,
		character,
		DeepCopyTimestamps(v),
		DeepCopyTimestamps(dv),
		DeepCopyTimestamps(tv),
		DeepCopyTimestamp(ov),
		DeepCopyTimestamp(uv)}
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

func DeepCopy(o Operation) Operation {
	return Operation{
		o.id,
		o.opType,
		o.position,
		o.character,
		DeepCopyTimestamps(o.v),
		DeepCopyTimestamps(o.dv),
		DeepCopyTimestamps(o.tv),
		DeepCopyTimestamp(o.ov),
		DeepCopyTimestamp(o.uv)}
}

func DeepCopyOperations(operations []Operation) []Operation {
	operationsCopy := make([]Operation, len(operations))

	for i, o := range operations {
		operationsCopy[i] = DeepCopy(o)
	}

	return operationsCopy
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
	return DeepCopyTimestamps(o.v)
}

// Appends the timestamp t to the timestamps slice of o.
func (o *Operation) AddV(t Timestamp) {
	o.v = append(o.v, DeepCopyTimestamp(t))
}

// Returns a copy of the slice containing the timestamps of operations that depend on operation o.
func (o *Operation) Dv() []Timestamp {
	return DeepCopyTimestamps(o.dv)
}

// Appends the Timestamp t to the timestamps slice of operations that depend on operation o.
func (o *Operation) AddDv(t Timestamp) {
	o.dv = append(o.dv, t.DeepCopy())
}

// Returns a slice containing the timestamps of operations whose effect objects tie with o.c.
func (o *Operation) Tv() []Timestamp {
	return DeepCopyTimestamps(o.tv)
}

// Appends the timestamp t to the timestamps slice of operations whose effect objects tie with o.c.
func (o *Operation) AddTv(t Timestamp) {
	o.tv = append(o.tv, t.DeepCopy())
}

// Returns the timestamp of the original operation o undoes (if operation o is an undo, otherwise nil).
func (o *Operation) Ov() Timestamp {
	ov := o.ov.DeepCopy()
	return ov
}

// Sets the timestamp ov of the operation o to t.
func (o *Operation) SetOv(t Timestamp) {
	o.ov = t.DeepCopy()
}

// Returns a copy of the timestamp of the operation that undoes o.
func (o *Operation) Uv() Timestamp {
	return o.uv.DeepCopy()
}

// Sets the timestamp uv of the operation o to t.
func (o *Operation) SetUv(t Timestamp) {
	o.uv = t.DeepCopy()
}

// Returns true if and only if operation o1 happened before operation o2.
func (o1 *Operation) HappenedBefore(o2 Operation) (bool, error) {
	for _, e1:= range o1.v{
		for _, e2:= range o2.v {
			happenedBefore, err := e1.HappenedBefore(e2)
			if err!=nil {
				return false, err
			} else {
				return happenedBefore, nil
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

type publicOp struct {
	Id SiteId
	OpType OpType
	Position Position
	Character Char
	V []Timestamp
	Dv []Timestamp
	Tv []Timestamp
	Ov Timestamp
	Uv Timestamp
}

//Transforms an Operation into a publicOp.
func OperationToPublicOp(o Operation) publicOp {
	copy := DeepCopy(o)
	return publicOp{
		copy.Id(),
		copy.OpType(),
		copy.Pos(),
		copy.Char(),
		copy.V(),
		copy.Dv(),
		copy.Tv(),
		copy.Ov(),
		copy.Uv()}
}

//Transforms a publicOp into an Operation.
func publicOpToOperation(o publicOp) Operation {
	return DeepCopy(NewOperation(
		o.Id,
		o.OpType,
		o.Position,
		o.Character,
		o.V,
		o.Dv,
		o.Tv,
		o.Ov,
		o.Uv))
}

type FrontendOperation struct {
	OpType OpType
	Character Char
	Position Position
}

func (frontendOperation *FrontendOperation) GetOperation(siteId SiteId) Operation {
	return PartialOperation(siteId, frontendOperation.OpType, frontendOperation.Position, frontendOperation.Character)
}
