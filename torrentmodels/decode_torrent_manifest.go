package torrentmodels

import (
	"crypto/sha1"
	"fmt"
	"path"

	"github.com/IncSW/go-bencode"
)

func DecodeTorrentManifest(data interface{}) TorrentManifest {
	torrentMap := data.(map[string]interface{})

	announce := string(torrentMap["announce"].([]byte))
	announceList := []string{}
	if torrentMap["announce-list"] != nil {
		for _, item := range torrentMap["announce-list"].([]interface{}) {
			announceList = append(announceList, string(item.([]interface{})[0].([]byte)))
		}
	}
	comment := ""
	if torrentMap["comment"] != nil {
		comment = string(torrentMap["comment"].([]byte))
	}
	createdBy := ""
	if torrentMap["created by"] != nil {
		createdBy = string(torrentMap["created by"].([]byte))
	}

	info := torrentMap["info"].(map[string]interface{})

	pieceLength := info["piece length"].(int64)
	torrentName := string(info["name"].([]byte))

	infoBytes, _ := bencode.Marshal(info)
	infoHash := sha1.Sum(infoBytes)

	pieces := info["pieces"].([]byte)

	pieceHashes := [][20]byte{}
	for i := 0; i < len(pieces); i += 20 {
		var currentHash [20]byte
		copy(currentHash[:], pieces[i:i+20])
		pieceHashes = append(pieceHashes, currentHash)
	}

	filesMetadata := []FileMetadata{}
	var offset int64

	files := []interface{}{info}

	if info["files"] != nil {
		fmt.Println("Files exist")
		files = info["files"].([]interface{})
	}
	for _, file := range files {
		file := file.(map[string]interface{})

		filePathParts := []string{torrentName}

		if file["path"] != nil {
			for _, part := range file["path"].([]interface{}) {
				filePathParts = append(filePathParts, string(part.([]byte)))
			}
		} else {
			filePathParts = append(filePathParts, torrentName)
		}

		fileSize := file["length"].(int64)
		filesMetadata = append(filesMetadata, FileMetadata{
			FilePath:   path.Join(filePathParts...),
			FileName:   filePathParts[len(filePathParts)-1],
			FileSize:   fileSize,
			FileOffset: offset,
		})

		offset += fileSize
	}

	return TorrentManifest{
		Announce:      announce,
		AnnounceList:  announceList,
		InfoHash:      infoHash,
		PieceHashes:   pieceHashes,
		PieceLength:   pieceLength,
		TotalLength:   offset,
		TorrentName:   torrentName,
		Comment:       comment,
		CreatedBy:     createdBy,
		FilesMetadata: filesMetadata,
	}
}
