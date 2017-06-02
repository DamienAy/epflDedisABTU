package encoding

const (
	LocalOp       = "localOperation"
	Undo          = "undo"
	AckRemoteOp   = "ackRemoteOperation"
	AckLocalOp    = "ackLocalOperation"
	NackLocalUndo = "nackLocalUndo"
	AckLocalUndo  = "ackLocalUndo"
	RemoteOp      = "remoteOperation"
)

type FrontendMessage struct {
	Type string
	Content []byte
}




