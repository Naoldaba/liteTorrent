package main

import (
	"crypto/rand"
	"fmt"

	"github.com/Naoldaba/Bit_Torrent/bitfield"
	"github.com/Naoldaba/Bit_Torrent/common"
	"github.com/Naoldaba/Bit_Torrent/download"
	"github.com/Naoldaba/Bit_Torrent/fileUtils"
	messageutils "github.com/Naoldaba/Bit_Torrent/messageUtils"
	"github.com/Naoldaba/Bit_Torrent/peer"
	"github.com/Naoldaba/Bit_Torrent/peerinteraction"
	"github.com/Naoldaba/Bit_Torrent/piecehandler"
	"github.com/Naoldaba/Bit_Torrent/seed"
	"github.com/Naoldaba/Bit_Torrent/torrentmodels"
	"github.com/Naoldaba/Bit_Torrent/tracker"
	"log"
	mrand "math/rand"
	"net"
	"time"
)

func main() {
	manifest := fileUtils.ReadManifestFromFile("debian-11.6.0-amd64-netinst.iso.torrent")

	blobFile := fileUtils.LoadOrCreateDownloadBlob(&manifest)
	memCopy := make([]byte, manifest.TotalLength)
	blobFile.ReadAt(memCopy, 0)

	currentBitField, bitfieldFile := bitfield.LoadOrCreateBitFieldFromFile(&manifest)

	totalDownloaded := countDownloadedPieces(currentBitField)

	fmt.Println("Total downloaded", totalDownloaded)

	id := [20]byte{}
	rand.Read(id[:])

	peerAddresses, err := tracker.GetPeersList(manifest, id, common.Port)
	if err != nil {
		fmt.Println("Can't get peers", err)
		panic(err)
	}
	fmt.Println(peerAddresses)

	workChannel := make(chan download.PieceJob, len(manifest.PieceHashes))
	pieceJobResultChannel := make(chan *download.PieceJobResult)
	seedRequestChannel := make(chan *seed.SeedRequest)

	createWorkForPieces(&manifest, currentBitField, &workChannel)

	peers := make([]*peer.Peer, len(peerAddresses))

	for i, peerAddress := range peerAddresses {
		go peerinteraction.StartPeerWorker(peers, i, peerAddress, id, manifest, common.Port, &workChannel, currentBitField, &pieceJobResultChannel, &seedRequestChannel, nil)
	}

	go func() {
		for {
			seedRequest := <-seedRequestChannel
			go seed.HandleSeedingRequest(seedRequest, memCopy, currentBitField, &manifest)
		}
	}()

	go startSeedingServer(&peers, id, manifest, common.Port, &workChannel, currentBitField, &pieceJobResultChannel, &seedRequestChannel)

	go optimisticUnchoking(peers)

	for {
		pieceJobResult := <-pieceJobResultChannel
		if pieceJobResult == nil {
			continue
		}

		copy(memCopy[pieceJobResult.PieceIndex*int(manifest.PieceLength):], pieceJobResult.PieceData)
		piecehandler.WritePieceToFile(&manifest, pieceJobResult, blobFile)
		currentBitField.MarkPiece(pieceJobResult.PieceIndex)
		currentBitField.WriteToFile(&manifest, bitfieldFile)

		totalDownloaded++
		fmt.Printf("Downloaded %v/%v pieces\n", totalDownloaded, len(manifest.PieceHashes))

		for _, peer := range peers {
			if peer != nil {
				messageutils.SendHaveMessage(peer, pieceJobResult.PieceIndex)
			}
		}

		if totalDownloaded == len(manifest.PieceHashes) {
			fmt.Println("Download finished")
			fileUtils.WriteBlobToFiles(&manifest)
		}
	}
}

func optimisticUnchoking(peers []*peer.Peer) {
	for {
		if len(peers) != 0 {
			peerIndex := mrand.Intn(len(peers))
			if peers[peerIndex] != nil {
				if peers[peerIndex].IsChoked {
					peers[peerIndex].IsChoked = false
					go messageutils.SendUnchokeMessage(peers[peerIndex])
				}
			}
		}
		time.Sleep(31 * time.Second)
	}
}

func startSeedingServer(peers *[]*peer.Peer, id [20]byte, manifest torrentmodels.TorrentManifest, port int, workChannel *chan download.PieceJob, currentBitField *bitfield.Bitfield, pieceJobResultChannel *chan *download.PieceJobResult, seedRequestChannel *chan *seed.SeedRequest) {
	ListenAddr := ":" + fmt.Sprint(port)
	listener, err := net.Listen("tcp", ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("Listening on %s...\n", ListenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		peerInstance := peer.Peer{
			Address: peer.PeerAddress{
				IP:   conn.RemoteAddr().(*net.TCPAddr).IP,
				Port: uint16(port),
			},
			Conn:       conn,
			Interested: false,
			IsChoked:   true,
			IsChoking:  false,
			BitField:   make([]byte, len(manifest.PieceHashes)),
		}

		*peers = append(*peers, &peerInstance)

		go peerinteraction.StartPeerWorker(*peers, len(*peers)-1, peerInstance.Address, id, manifest, port, workChannel, currentBitField, pieceJobResultChannel, seedRequestChannel, &conn)
	}
}

func createWorkForPieces(manifest *torrentmodels.TorrentManifest, currentBitField *bitfield.Bitfield, workChannel *chan download.PieceJob) {
	for index, hash := range manifest.PieceHashes {
		if !currentBitField.HasPiece(index) {
			*workChannel <- download.PieceJob{
				PieceIndex:  index,
				PieceHash:   hash,
				PieceLength: piecehandler.GetPieceLength(index, int(manifest.PieceLength), int(manifest.TotalLength)),
			}
		}
	}
}

func countDownloadedPieces(bitField *bitfield.Bitfield) int {
	totalDownloaded := 0

	for _, piece := range *bitField {
		for i := 0; i < 8; i++ {
			if piece&(1<<uint(i)) != 0 {
				totalDownloaded++
			}
		}
	}

	return totalDownloaded
}
