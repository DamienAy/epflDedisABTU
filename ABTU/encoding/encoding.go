package encoding

import (
	"encoding/json"
	"errors"
)

const (
	localOp = "localOperation"
	Undo = "undo"
	AckRemoteOp = "ackRemoteOperation"
	AckLocalOperation = "ackLocalOperation"
	NackLocalUndo = "nackLocalUndo"
	AckLocalUndo = "ackLocalUndo"
	RemoteOp = "remoteOperation"
)

type FrontendMessage struct {
	Type string
	Content []byte
}




