package download

type PieceJobProgress struct {
	PieceIndex int
	Buffer []byte
	TotalDownloaded int
	PieceLength int
}
