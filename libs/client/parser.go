package client

import (
	"errors"
	"net"
)

func parseRequestId(conn *net.UDPConn) (requestId int64, err error) {
	var data = make([]byte, 1)
	_, err = conn.Read(data)
	if err != nil {
		return
	}
	requestId = int64(data[0])
	return
}

func parseLength(conn *net.UDPConn) (length int64, err error) {
	var data = make([]byte, 1)
	_, err = conn.Read(data)
	if err != nil {
		return
	}
	length = int64(data[0])
	return
}

func parseData(length int64, conn *net.UDPConn) (data []byte, err error) {
	data = make([]byte, length)

	var readCount int
	readCount, err = conn.Read(data)
	if err != nil {
		return
	}

	if readCount != int(length) {
		err = errors.New("data_length_not_equals_to_available_data")
		data = nil
	}
	return
}
