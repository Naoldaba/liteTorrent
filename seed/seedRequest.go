package seed

import (
	"github.com/Naoldaba/Bit_Torrent/peer"
	"github.com/Naoldaba/Bit_Torrent/peercommunication"
)

type SeedRequest struct {
	Peer    *peer.Peer
	Message *peercommunication.Message
}
