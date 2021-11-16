package main

import (
	"fmt"
	"net"
	"time"

	"github.com/aliforever/go-nat-punch-hole/libs/client"

	"github.com/aliforever/go-nat-punch-hole/args"
	"github.com/aliforever/go-nat-punch-hole/libs/server"
)

func main() {
	if args.Type == "" {
		panic("please specify type by --type=server/client")
	}

	if args.Type == "client" && args.SignalingServerAddress == "" {
		panic("please specify signaling server address using --server=ip/domain")
	}

	if args.Type == "server" {
		udpServer := server.NewUDP(args.LocalAddress)
		fmt.Println("Starting listening on", args.LocalAddress)
		err := udpServer.Listen()
		if err != nil {
			panic(err)
		}
	}

	if args.Type == "client" {
		udpClient := client.NewUDP(args.SignalingServerAddress, args.LocalAddress)
		err := udpClient.Start()
		if err != nil {
			panic(fmt.Sprintf("Error in udp client: %s", err))
		}
		ok, err := udpClient.Register("test_room")
		if err != nil {
			fmt.Println(err)
			return
		}
		if !ok {
			return
		}

		var peerAddr *net.UDPAddr
		for {
			peerAddr, err = udpClient.GetPeerAddress("test_room")
			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Second * 1)
				continue
			}
			fmt.Printf("peer address is: %s\n", peerAddr)
			break
		}

		if peerAddr == nil {
			return
		}

		udpClient.RemoveRoom("test_room")

		/*if strings.Contains(args.LocalAddress, "8182") {
			peerAddr, _ = net.ResolveUDPAddr("udp", "5.122.40.48:8183")
		} else if strings.Contains(args.LocalAddress, "8183") {
			peerAddr, _ = net.ResolveUDPAddr("udp", "5.122.40.48:8182")
		}*/

		for true {
			err = udpClient.ConnectToPeer(peerAddr)
			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Second * 1)
				continue
			}
			time.Sleep(time.Second * 1)
		}

	}

	for true {
		time.Sleep(time.Second * 1)
	}
}
