package handshake

type HandShake struct {
	HeaderText string   
	InfoHash   [20]byte 
	PeerId     [20]byte
}

func New(infoHash, peerId [20]byte) HandShake {
	return HandShake{
		HeaderText: "BitTorrent protocol",
		InfoHash:   infoHash,
		PeerId:     peerId,
	}
}


func (handShake *HandShake) ToBytes() []byte {
	buf := make([]byte, len(handShake.HeaderText)+49)
	buf[0] = byte(len(handShake.HeaderText))
	curr := 1
	curr += copy(buf[curr:], []byte(handShake.HeaderText))
	curr += copy(buf[curr:], make([]byte, 8)) // 8 reserved bytes
	curr += copy(buf[curr:], handShake.InfoHash[:])
	curr += copy(buf[curr:], handShake.PeerId[:])
	return buf
}
