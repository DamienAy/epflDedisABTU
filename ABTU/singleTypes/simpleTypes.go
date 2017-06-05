package singleTypes

type SiteId int

type OpType int
const(INS OpType = iota;
	DEL;
	UNIT)

type Position int
type Char []byte