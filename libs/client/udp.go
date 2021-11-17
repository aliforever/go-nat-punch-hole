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

		var (
			readCount int
			udpAddr   *net.UDPAddr
		)
		readCount, udpAddr, err = u.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("error reading incoming data", err)
			return
		}

		buffer = buffer[0:readCount]

		var requestId = int64(buffer[0])
		var length = int64(buffer[1])
		var data = buffer[2:]

		if int(length) != len(data) {
			err = errors.New("data_length_doesnt_match_given_length")
			fmt.Println(data, string(data))
			fmt.Println("error reading incoming data", err)
			break
		}

		ch, _ := u.responseChan.findChannel(requestId)
		if ch != nil {
			ch <- data
		} else {
			fmt.Printf("received %s from %s\n", data, udpAddr)
		}

	}
	return
}

func (u *UDP) GetPeerAddress(roomName string) (peerAddress *net.UDPAddr, err error) {
	_, err = u.conn.WriteToUDP(makeRequestBytes(1, findPeerByNameAction, []byte(roomName)), u.signalingServerAddr)
	if err != nil {
		return
	}

	var ch = make(chan []byte)
	u.responseChan.addChannel(1, ch)

	data := <-ch

	if string(data) == "PEER_NOT_FOUND" {
		err = errors.New("peer_not_found")
		return
	}

	peerAddress, err = net.ResolveUDPAddr("udp", string(data))
	return
}

func (u *UDP) Register(roomName string) (success bool, err error) {
	_, err = u.conn.WriteToUDP(makeRequestBytes(1, registerAction, []byte(roomName)), u.signalingServerAddr)
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

func (u *UDP) RemoveRoom(roomName string) (success bool, err error) {
	_, err = u.conn.WriteToUDP(makeRequestBytes(1, deleteRoomAction, []byte(roomName)), u.signalingServerAddr)
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

func (u *UDP) ConnectToPeer(peerAddress *net.UDPAddr) (err error) {
	// fmt.Println("Attempting to connect to", peerAddress.String())
	_, err = u.conn.WriteToUDP(makeRequestBytes(3, []byte("MESSAGE"), []byte("Hello")), peerAddress)
	return
}
