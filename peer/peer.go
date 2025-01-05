package peer

import (
	"github.com/Naoldaba/Bit_Torrent/bitfield"
	"net"
)

type Peer struct {
	Conn       net.Conn
	Address    PeerAddress
	Interested bool
	IsChoked   bool
	IsChoking  bool
	BitField   bitfield.Bitfield
}
