package protocol

type Protocol interface {
	MsgLen(data []byte) (length uint64, err error)
	Unpack(data []byte) (msg []byte, err error)
	Pack(data []byte) (smg []byte, err error)
}
