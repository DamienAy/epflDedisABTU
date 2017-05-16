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
	AckLocalUndo = "ackLocalUndo"
	RemoteOp = "remoteOperation"
)

type FrontendMessage struct {
	Type string
	Content []byte
}




