package client

func makeRequestBytes(requestId int64, data []byte) (requestBytes []byte) {
	requestBytes = make([]byte, len(data)+2)
	requestBytes[0] = byte(requestId)
	requestBytes[1] = byte(len(data))
	copy(requestBytes[2:], data)
	return
}
