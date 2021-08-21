package main

import (
	"flag"
	"learn_tun"
	"log"
	"net"
)

var tunnelAddress string
var tunnelPeerAddress string
var serverAddress string
var clientAddress string

func init() {
	flag.StringVar(&tunnelAddress, "tunnelAddress", "10.0.0.2", "tunnel interface address")
	flag.StringVar(&tunnelPeerAddress, "tunnelPeerAddress", "10.0.0.1", "tunnel interface address")
	flag.StringVar(&serverAddress, "serverAddress", "", "server address ip:port")
	flag.StringVar(&clientAddress, "clientAddress", "", "client address ip:port")

	flag.Parse()
}

func main() {
	checkFlags()

	iface := learn_tun.CreateTun("")
	defer iface.Close()

	learn_tun.TunCommand(tunnelAddress, tunnelPeerAddress, iface.Name())

	sock := learn_tun.DialUdp(serverAddress)
	defer sock.Close()

	//remoteSock := learn_tun.ListenUdp(serverAddress)
	//defer remoteSock.Close()

	go func() {
		for {
			b := make([]byte, 1<<16)
			n, err := iface.Read(b)
			if err != nil {
				log.Fatalln(err)
			}
			_, err = sock.Write(b[:n])
			if err != nil {
				log.Println("write msg", err)
			}
		}
	}()

	go func() {
		for {
			b := make([]byte, 1<<16)
			n, raddr, err := sock.ReadFrom(b)
			if err != nil {
				log.Println("receive from msg", raddr)
				continue
			}
			_, err = iface.Write(b[:n])
			if err != nil {
				log.Println("iface write", err)
				continue
			}
		}
	}()

	learn_tun.RegisterSignal()
}

func checkFlags() {
	if net.ParseIP(tunnelAddress) == nil {
		log.Fatalln("tunnelAddress is unvalid ip")
	}
	if net.ParseIP(tunnelPeerAddress) == nil {
		log.Fatalln("tunnelPeerAddress is unvalid ip")
	}
	if serverAddress == "" {
		log.Fatalln("serverAddress is empty")
	}
}
