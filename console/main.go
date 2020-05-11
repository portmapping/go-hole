package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/goextension/tool"
	"github.com/portmapping/lurker"
)

func main() {
	lurker.DefaultTCP = 16004
	lurker.DefaultUDP = 16005
	rnd := tool.GenerateRandomString(16)
	address := ""
	list := sync.Map{}
	if len(os.Args) > 2 {
		address = os.Args[2]
	}
	if len(os.Args) > 3 {
	}
	l := lurker.New()
	listener, err := l.Listen()
	if err != nil {
		panic(err)
		return
	}
	go func() {
		for source := range listener {
			go func(s lurker.Source) {
				_, ok := list.Load(s.Service().ID)
				if ok {
					return
				}
				err := s.TryConnect()
				fmt.Println("reverse connected:", err)
				if err != nil {
					list.Store(s.Service().ID, s)
				}
			}(source)
		}
	}()

	if len(os.Args) > 2 {
		addr, i := lurker.ParseAddr(address)

		internalAddress, err := l.NAT().GetInternalAddress()
		if err != nil {
			return
		}
		fmt.Println("remote addr:", addr.String(), i)
		fmt.Println("your connect id:", rnd)
		s := lurker.NewSource(lurker.Service{
			ID:       rnd,
			ISP:      internalAddress,
			PortUDP:  l.PortUDP(),
			PortHole: l.PortHole(),
			PortTCP:  l.PortTCP(),
			ExtData:  nil,
		}, lurker.Addr{
			Protocol: "tcp",
			IP:       addr,
			Port:     16004,
		})
		go func() {
			b := s.TryConnect()
			fmt.Println("target connected:", b)
		}()

	}
	fmt.Println("ready for waiting")
	time.Sleep(30 * time.Minute)
}
