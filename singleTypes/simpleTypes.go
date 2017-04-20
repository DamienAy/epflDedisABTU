package singleTypes

type SiteId uint8

type OpType int
const(INS OpType = iota; DEL; UNIT)

type Position uint64
type Char byte
