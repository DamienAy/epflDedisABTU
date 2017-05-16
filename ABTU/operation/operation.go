package operation

import (
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes";
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp";
	"errors"
	"encoding/json"
)

// Represents an operation as defined in the ABTU paper.
// All fields are private.
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

// Returns a new Operation, copies all arguments.
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

// Returns a deep copy of the operation o
// Timestamps and slices of timestamps are also deep copied.
func DeepCopyOperation(o Operation) Operation {
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

// Returns a deep copy of the slice of Operation operations.
func DeepCopyOperations(operations []Operation) []Operation {
	operationsCopy := make([]Operation, len(operations))

	for i, o := range operations {
		operationsCopy[i] = DeepCopyOperation(o)
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

// Appends a copy of the timestamp t to the timestamps slice of o.
func (o *Operation) AddV(t Timestamp) {
	o.v = append(o.v, DeepCopyTimestamp(t))
}

// Returns a copy of the slice containing the timestamps of operations that depend on operation o.
func (o *Operation) Dv() []Timestamp {
	return DeepCopyTimestamps(o.dv)
}

// Appends a copy of the Timestamp t to the timestamps slice of operations that depend on operation o.
func (o *Operation) AddDv(t Timestamp) {
	o.dv = append(o.dv, DeepCopyTimestamp(t))
}

// Appends a copy of the []Timestamp timestamps to o.Dv.
func (o *Operation) AddAllDv(timestamps []Timestamp) {
	for _, t := range timestamps {
		o.AddDv(t)
	}
}

// Returns a copy of the slice containing the timestamps of operations whose effect objects tie with o.c.
func (o *Operation) Tv() []Timestamp {
	return DeepCopyTimestamps(o.tv)
}

// Appends a copy of the timestamp t to the timestamps slice of operations whose effect objects tie with o.c.
func (o *Operation) AddTv(t Timestamp) {
	o.tv = append(o.tv, DeepCopyTimestamp(t))
}

// Appends a copy of the []Timestamp timestamps to o.tv.
func (o *Operation) AddAllTv(timestamps []Timestamp) {
	for _, t := range timestamps {
		o.AddTv(t)
	}
}

// Returns a copy of the timestamp of the original operation o undoes (if operation o is an undo, otherwise nil).
func (o *Operation) Ov() Timestamp {
	return DeepCopyTimestamp(o.ov)
}

// Sets the timestamp ov of the operation o to a copy of t.
func (o *Operation) SetOv(t Timestamp) {
	o.ov = DeepCopyTimestamp(t)
}

// Returns a copy of the timestamp of the operation that undoes o.
func (o *Operation) Uv() Timestamp {
	return DeepCopyTimestamp(o.uv)
}

// Sets the timestamp uv of the operation o to a copy of t.
func (o *Operation) SetUv(t Timestamp) {
	o.uv = DeepCopyTimestamp(t)
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

// Same as Operation, all fields are public
// Used for encoding.
type PublicOperation struct {
	Id SiteId
	OpType OpType
	Position Position
	Character Char
	V []PublicTimestamp
	Dv []PublicTimestamp
	Tv []PublicTimestamp
	Ov PublicTimestamp
	Uv PublicTimestamp
}

// Returns the PublicOperation corresponding to the Operation o
// The Timestamps contained in the Operation o are also transformed to PublicTimestamps
func OperationToPublicOperation(o Operation) PublicOperation {
	copy := DeepCopyOperation(o)
	return PublicOperation{
		copy.Id(),
		copy.OpType(),
		copy.Pos(),
		copy.Char(),
		TimestampsToPublicTimestamps(copy.V()),
		TimestampsToPublicTimestamps(copy.Dv()),
		TimestampsToPublicTimestamps(copy.Tv()),
		TimestampToPublicTimestamp(copy.Ov()),
		TimestampToPublicTimestamp(copy.Uv())}
}

// Returns the Operation corresoponding to the PublicOperation publicOP.
// The PublicTimestamps contained in the PublicOperation publicOp are also transformed to Timestamps.
func publicOperationToOperation(publicOp PublicOperation) Operation {
	return DeepCopyOperation(NewOperation(
		publicOp.Id,
		publicOp.OpType,
		publicOp.Position,
		publicOp.Character,
		PublicTimestampsToTimestamps(publicOp.V),
		PublicTimestampsToTimestamps(publicOp.Dv),
		PublicTimestampsToTimestamps(publicOp.Tv),
		PublicTimestampToTimestamp(publicOp.Ov),
		PublicTimestampToTimestamp(publicOp.Uv)))
}

// Represents an operation as sent to frontend.
// Only contains useful information for frontend.
type FrontendOperation struct {
	OpType OpType
	Character Char
	Position Position
}

// Returns the Operation corresponding to the frontendOperation.
func FrontendOperationToOperation(frontendOperation FrontendOperation, siteId SiteId) Operation {
	return PartialOperation(siteId, frontendOperation.OpType, frontendOperation.Position, frontendOperation.Character)
}

// Returns the FrontendOperation corresponding to the Operation.
func OperationToFrontendOperation(operation Operation) FrontendOperation {
	return FrontendOperation{operation.opType, operation.character, operation.position}
}

// Returns the encoding in json format (frontend) corresponding to the operation o.
// Returns an error if the encoding failed.
func (o *Operation) EncodeFrontend() ([]byte, error) {
	bytes, err := json.Marshal(&OperationToFrontendOperation(*o))
	if err != nil {
		return nil, errors.New("Json encoding failed :" + err.Error())
	}

	return bytes, nil
}

// Returns the operation corresponding to the json encoding of a frontendOperation.
// Returns an error if the decoding failed.
func DecodeFrontend(bytes []byte, siteId SiteId) (Operation, error) {
	var frontendOperation FrontendOperation

	err := json.Unmarshal(bytes, &frontendOperation)

	if err != nil {
	return nil, errors.New("Json decoding failed :" + err.Error())
	}

	return FrontendOperationToOperation(frontendOperation, siteId), nil
}

// Returns the encoding in json format (for peers) corresponding to the operation o.
// Returns an error if the encoding failed.
func (o *Operation) EncodePeers() ([]byte, error) {
	bytes, err := json.Marshal(OperationToPublicOperation(*o))
	if err != nil {
		return nil, errors.New("Json encoding failed :" + err.Error())
	}

	return bytes, nil
}

// Returns the operation corresponding to the json encoding of a publicOperation.
// Returns an error if the decoding failed.
func DecodePeers(bytes []byte) (Operation, error) {
	var publicOperation PublicOperation

	err := json.Unmarshal(bytes, &publicOperation)

	if err != nil {
		return nil, errors.New("Json decoding failed :" + err.Error())
	}

	return publicOperationToOperation(publicOperation), nil
}