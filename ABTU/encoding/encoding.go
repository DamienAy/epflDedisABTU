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

func EncodeFrontendMessage(frontendMsg FrontendMessage) ([]byte, error) {
	bytes, err := json.Marshal(frontendMsg)

	if err != nil {
		return nil, errors.New("Json encoding failed :" + err.Error())
	}

	return bytes, nil
}


