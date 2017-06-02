package singleTypes

type SiteId int32

type OpType int
const(INS OpType = iota;
	DEL;
	UNIT)

type Position int32
type Char []byte
