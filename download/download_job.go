package download

type PieceJob struct {
	PieceIndex int
	PieceHash [20]byte
	PieceLength int
}
