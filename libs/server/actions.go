package server

var (
	actionsSplitterBytes = []byte(":")
	registerAction       = []byte("REGISTER")
	findPeerByNameAction = []byte("PEER")
	deleteRoomAction     = []byte("DELETE_ROOM")
)
