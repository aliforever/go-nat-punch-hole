package server

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

type UDP struct {
	localAddress string
	peers        peerStorage
}

func NewUDP(localAddress string) UDP {
	return UDP{
		localAddress: localAddress,
		peers:        newPeerStorage(),
	}
}

func (u *UDP) Listen() (err error) {
	var addr *net.UDPAddr
	addr, err = net.ResolveUDPAddr("udp", u.localAddress)
	if err != nil {
		return
	}

	var conn *net.UDPConn
	conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return
	}

	err = u.processUpdates(conn)
	return
}

func (u *UDP) processUpdates(conn *net.UDPConn) (err error) {
	for {
		var buffer = make([]byte, 2048)

		var (
			readCount int
			udpAddr   *net.UDPAddr
		)
		readCount, udpAddr, err = conn.ReadFromUDP(buffer)
		if err != nil {
			break
		}

		buffer = buffer[0:readCount]
		var requestId = int64(buffer[0])
		var length = int64(buffer[1])
		var data = buffer[2:]

		if int(length) != len(data) {
			err = errors.New("invalid_data_length")
			break
		}

		byteData := bytes.Split(data, actionsSplitterBytes)

		if bytes.Equal(byteData[0], registerAction) {
			u.peers.storePeer(string(byteData[1]), udpAddr)
			conn.WriteToUDP(makeResponseBytes(requestId, okMessage), udpAddr)
			fmt.Println("stored peer", udpAddr.String())
		} else if bytes.Equal(byteData[0], findPeerByNameAction) {
			// Find Peer By Name
			peer, _ := u.peers.findPeer(string(byteData[1]), udpAddr)
			if peer == nil {
				conn.WriteToUDP(makeResponseBytes(requestId, peerNotFoundMessage), udpAddr)
			} else {
				conn.WriteToUDP(makeResponseBytes(requestId, []byte(peer.String())), udpAddr)
			}
		} else if bytes.Equal(byteData[0], deleteRoomAction) {
			u.peers.removeRoom(string(byteData[1]))
			conn.WriteToUDP(makeResponseBytes(requestId, okMessage), udpAddr)
		} else {
			err = errors.New("connection_closed_due_to_invalid_action")
			conn.WriteToUDP(makeResponseBytes(requestId, invalidActionMessage), udpAddr)
			break
		}
	}
	return
}
