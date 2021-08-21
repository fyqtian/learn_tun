package main

import (
	"flag"
	"learn_tun"
	"log"
	"net"
)

var tunnelAddress string
var serverAddress string

func init() {
	flag.StringVar(&tunnelAddress, "tunnelAddress", "10.0.0.1", "tunnel interface address")
	flag.StringVar(&serverAddress, "serverAddress", "0.0.0.0:8000", "server udp address")

	flag.Parse()
}

func main() {
	checkFlags()
	iface := learn_tun.CreateTun("")
	learn_tun.TunCommand(tunnelAddress, "", iface.Name())

	sock := learn_tun.ListenUdp(serverAddress)
	go func() {
		for {
			b := make([]byte, 1<<16)
			var n int
			var err error
			var rAddr *net.UDPAddr
			n, rAddr, err = sock.ReadFromUDP(b)
			if err != nil {
				log.Fatalln(err)
				continue
			}
			go func(r *net.UDPAddr) {
				for {
					b := make([]byte, 1<<16)
					n, err := iface.Read(b)
					if err != nil {
						log.Println("iface read", err)
						continue
					}
					log.Println("ifcae receive", b[:n])

					_, err = sock.WriteToUDP(b[:n], r)
					if err != nil {
						log.Println("sock write", err)
					}

				}
			}(rAddr)
			log.Println("udp receive", b[:n])
			_, err = iface.Write(b[:n])
			if err != nil {
				log.Println("write msg", err)
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

	if serverAddress == "" {
		log.Fatalln("serverAddress is empty")
	}
}
