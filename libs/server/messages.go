package server

var (
	okMessage            = []byte("OK")
	peerNotFoundMessage  = []byte("PEER_NOT_FOUND")
	invalidActionMessage = []byte("INVALID_ACTION")
)

func makeResponseBytes(requestId int64, data []byte) (responseBytes []byte) {
	responseBytes = make([]byte, len(data)+2)
	responseBytes[0] = byte(requestId)
	responseBytes[1] = byte(len(data))
	copy(responseBytes[2:], data)
	return
}
