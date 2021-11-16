package server

import (
	"errors"
	"net"
	"sync"
)

type peerStorage struct {
	sync.Mutex
	peers map[string][]*net.UDPAddr
}

func newPeerStorage() peerStorage {
	return peerStorage{peers: map[string][]*net.UDPAddr{}}
}

func (ps *peerStorage) storePeer(name string, peer *net.UDPAddr) {
	ps.Lock()
	defer ps.Unlock()
	ps.peers[name] = append(ps.peers[name], peer)
}

func (ps *peerStorage) findPeer(name string, currentPeer *net.UDPAddr) (peer *net.UDPAddr, err error) {
	ps.Lock()
	defer ps.Unlock()
	if room, ok := ps.peers[name]; !ok {
		err = errors.New("invalid_name")
		return
	} else {
		for _, conn := range room {
			if conn.String() != currentPeer.String() {
				peer = conn
				return
			}
		}
	}
	err = errors.New("peer_not_found")
	return
}
