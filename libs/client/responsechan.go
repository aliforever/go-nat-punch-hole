package client

import (
	"errors"
	"sync"
)

type responseChan struct {
	sync.Mutex
	responses map[int64]chan []byte
}

func newResponseChannel() responseChan {
	return responseChan{
		responses: map[int64]chan []byte{},
	}
}

func (rc *responseChan) addChannel(requestId int64, ch chan []byte) {
	rc.Lock()
	defer rc.Unlock()
	rc.responses[requestId] = ch
}

func (rc *responseChan) findChannel(requestId int64) (ch chan []byte, err error) {
	rc.Lock()
	defer rc.Unlock()

	var found bool
	ch, found = rc.responses[requestId]
	if !found {
		err = errors.New("not_found")
	}
	return
}
