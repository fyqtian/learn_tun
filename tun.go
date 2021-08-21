package learn_tun

import (
	"github.com/songgao/water"
	"log"
	"os/exec"
	"reflect"
	"runtime"
)

func CreateTun(name string) *water.Interface {
	conf := water.Config{DeviceType: water.TUN}
	if runtime.GOOS == "linux" && name != "" {
		// Use reflect to avoid separate build file for linux-only.
		reflect.ValueOf(&conf).Elem().FieldByName("Name").SetString(name)
	}
	iface, err := water.New(conf)
	if err != nil {
		log.Fatalf("create tun err %v", err)
	}
	log.Printf("created tun device: %v\n", iface.Name())
	return iface
}

func TunCommand(local, remote, name string) {
	switch runtime.GOOS {
	case "linux":
		if err := exec.Command("/sbin/ip", "link", "set", "dev", name, "mtu", "1300").Run(); err != nil {
			log.Fatalf("ip link error: %v", err)
		}
		if err := exec.Command("/sbin/ip", "addr", "add", local+"/24", "dev", name).Run(); err != nil {
			log.Fatalf("ip addr error: %v", err)
		}
		if err := exec.Command("/sbin/ip", "link", "set", "dev", name, "up").Run(); err != nil {
			log.Fatalf("ip link error: %v", err)
		}
	case "darwin":
		if err := exec.Command("/sbin/ifconfig", name, "mtu", "1300", local, remote, "up").Run(); err != nil {
			log.Fatalf("ifconfig error: %v", err)
		}
	default:
		log.Fatalf("no tun support for: %v", runtime.GOOS)
	}
}
