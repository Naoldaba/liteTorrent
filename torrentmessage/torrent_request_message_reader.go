package torrentmessage

import (
	"encoding/binary"
	"errors"
)

func ParseRequestMessage(payload []byte) (RequestMessage, error) {
	if len(payload) < 12 {
		return RequestMessage{}, errors.New("invalid payload length for request message")
	}

	pieceIndex := binary.BigEndian.Uint32(payload[0:4])
	begin := binary.BigEndian.Uint32(payload[4:8])
	length := binary.BigEndian.Uint32(payload[8:12])

	return RequestMessage{
		PieceIndex: int(pieceIndex),
		Begin:      int(begin),
		Length:     int(length),
	}, nil
}
