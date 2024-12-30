package handshake

import "io"


func ReadHandShake(reader io.Reader) (HandShake, error) {
	var headerText [20]byte
	var reserved [8]byte
	var infoHash [20]byte
	var peerId [20]byte

	_, err := reader.Read(headerText[:])
	if err != nil {
		return HandShake{}, err
	}
	_, err = reader.Read(reserved[:])
	if err != nil {
		return HandShake{}, err
	}
	_, err = reader.Read(infoHash[:])
	if err != nil {
		return HandShake{}, err
	}
	_, err = reader.Read(peerId[:])
	if err != nil {
		return HandShake{}, err
	}

	return HandShake{
		HeaderText: string(headerText[:]),
		InfoHash:   infoHash,
		PeerId:     peerId,
	}, nil
}
