package client

func makeRequestBytes(requestId int64, action, data []byte) (requestBytes []byte) {
	var requestData = make([]byte, len(action)+1+len(data))
	copy(requestData[0:], action)
	copy(requestData[len(action):], actionsSplitterBytes)
	copy(requestData[len(action)+len(actionsSplitterBytes):], data)

	requestBytes = make([]byte, len(requestData)+2)
	requestBytes[0] = byte(requestId)
	requestBytes[1] = byte(len(requestData))
	copy(requestBytes[2:], requestData)
	return
}
