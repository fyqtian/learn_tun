package learn_tun

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func RegisterSignal() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("received %v - initiating shutdown\n", <-sigc)
	os.Exit(1)
}

func ListenUdp(addr string) *net.UDPConn {
	udpAddr := ParseUdpAddr(addr)
	sock, err := net.ListenUDP("udp4", udpAddr)

	if err != nil {
		log.Fatalln("listen udp=", udpAddr, err)
	}
	return sock
}

func DialUdp(addr string) *net.UDPConn {
	udpAddr := ParseUdpAddr(addr)

	sock, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalln("listen udp=", udpAddr, err)
	}
	return sock
}

func ParseUdpAddr(addr string) *net.UDPAddr {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalln("ResolveUDPAddr ", err)
	}
	return udpAddr
}
