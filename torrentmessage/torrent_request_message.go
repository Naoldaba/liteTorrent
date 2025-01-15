package torrentmessage

import (
	"encoding/binary"
)

type RequestMessage struct {
	PieceIndex int
	Begin      int
	Length     int
}

func (req RequestMessage) ToBytes() []byte {
	bytes := make([]byte, 0)
	bytes = binary.BigEndian.AppendUint32(bytes, uint32(req.PieceIndex))
	bytes = binary.BigEndian.AppendUint32(bytes, uint32(req.Begin))
	bytes = binary.BigEndian.AppendUint32(bytes, uint32(req.Length))
	return bytes
}
