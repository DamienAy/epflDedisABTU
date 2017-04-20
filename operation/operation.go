package operation

import "errors"
const N int = 3

type SiteId int

type OpType int
const(INS OpType = iota; DEL; UNIT)

type Position int
type Char byte

type Timestamp [N]int




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
func (o Operation) getId() SiteId {
	return o.id
}

// Returns the OpType of the operation o.
func (o Operation) getOpType() OpType {
	return o.opType
}

// Returns the Position of the operation o.
func (o Operation) getPos() Position {
	return o.position
}

/*
Sets the positition of the operation o to the value p.
Returns an error if the value p is negative.
 */
func (o Operation) setPos(p Position) error {
	if p >= 0 {
		o.position = p
	} else {
		return errors.New("Position of an operation cannot be less than 0.")
	}

	return nil
}

// Returns the character of the operation o.
func (o Operation) getChar() Char {
	return o.character
}


// Returns a slice containing the timestamps of operation o.
func (o Operation) getV() []Timestamp {
	return o.v
}

// Appends the timestamp t to the timestamps slice of o.
func (o Operation) addV(t Timestamp) {
	o.v = append(o.v, t)
}

// Returns a slice containing the timestamps of operations that depend on operation o.
func (o Operation) getDv() []Timestamp {
	return o.dv
}

// Appends the Timestamp t to the timestamps slice of operations that depend on operation o.
func (o Operation) addDv(t Timestamp) {
	o.dv = append(o.dv, t)
}

// Returns a slice containing the timestamps of operations whose effect objects tie with o.c.
func (o Operation) getTv() []Timestamp {
	return o.tv
}

// Appends the timestamp t to the timestamps slice of operations whose effect objects tie with o.c.
func (o Operation) addTv(t Timestamp) {
	o.tv = append(o.tv, t)
}

// Returns the timestamp of the original operation o undoes (if operation o is an undo, otherwise nil).
func (o Operation) getOv() Timestamp {
	return o.ov
}

// Sets the timestamp ov of the operation o to t.
func (o Operation) setOv(t Timestamp) {
	o.ov = t
}

// Returns the timestamp of the operation that undoes o.
func (o Operation) getUv() Timestamp {
	return o.ov
}

// Sets the timestamp uv of the operation o to t.
func (o Operation) setUv(t Timestamp) {
	o.ov = t
}




