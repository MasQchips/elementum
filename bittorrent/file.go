package bittorrent

// File ...
type File struct {
	Index      int
	Name       string
	Size       int64
	Path       string
	Offset     int64
	PieceStart int
	PieceEnd   int
}
