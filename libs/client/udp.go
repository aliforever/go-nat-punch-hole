package client

import (
	"errors"
	"fmt"
	"net"
)

type UDP struct {
	signalingServerAddress string
	signalingServerAddr    *net.UDPAddr
	localServerAddress     string
	responseChan           responseChan
	conn                   *net.UDPConn
}

func NewUDP(signalingServerAddress, localServerAddress string) UDP {
	return UDP{signalingServerAddress: signalingServerAddress, localServerAddress: localServerAddress, responseChan: newResponseChannel()}
}

func (u *UDP) Start() (err error) {
	u.signalingServerAddr, err = net.ResolveUDPAddr("udp", u.signalingServerAddress)
	if err != nil {
		return
	}

	var addr *net.UDPAddr
	addr, err = net.ResolveUDPAddr("udp", u.localServerAddress)
	if err != nil {
		return
	}

	u.conn, err = net.ListenUDP("udp", addr)
	if err == nil {
		go u.processUpdates()
	}
	return
}

func (u *UDP) processUpdates() (err error) {
	for {
		var buffer = make([]byte, 2048)

		fmt.Println("incoming data")

		var readCount int
		readCount, err = u.conn.Read(buffer)
		if err != nil {
			fmt.Println("error reading incoming data", err)
			return
		}

		fmt.Println("reading incoming data till here")

		buffer = buffer[0:readCount]

		var requestId = int64(buffer[0])

		ch, _ := u.responseChan.findChannel(requestId)
		if ch == nil {
			fmt.Println("not found the response channel by request id")
			continue
		}

		var length = int64(buffer[1])
		var data = buffer[2:]

		if int(length) != len(data) {
			err = errors.New("data_length_doesnt_match_given_length")
			fmt.Println("error reading incoming data", err)
			break
		}

		ch <- data
		fmt.Println("written response to channel")
	}
	return
}

func (u *UDP) Register(roomName string) (success bool, err error) {
	var requestData = make([]byte, len(roomName)+1+len(registerAction))
	copy(requestData[0:], registerAction)
	copy(requestData[len(registerAction):], actionsSplitterBytes)
	copy(requestData[len(registerAction)+len(actionsSplitterBytes):], []byte(roomName))
	_, err = u.conn.WriteToUDP(makeRequestBytes(1, requestData), u.signalingServerAddr)
	if err != nil {
		return
	}

	var ch = make(chan []byte)
	u.responseChan.addChannel(1, ch)

	data := <-ch
	if string(data) == "OK" {
		success = true
	}
	return
}
