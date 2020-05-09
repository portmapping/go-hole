package lurker

import (
	"encoding/json"
	"github.com/portmapping/go-reuse"
	"net"
	"strconv"
	"sync"
)

// Source ...
type Source interface {
	TryConnect() error
	//Support() Support
}

// Addr ...
type Addr struct {
	Protocol string
	IP       net.IP
	Port     int
}

// Service ...
type Service struct {
	ID       string
	ISP      net.IP
	PortUDP  int
	PortHole int
	PortTCP  int
	ExtData  []byte
}

type source struct {
	addr    Addr
	service Service
	support Support
}

// Network ...
func (addr Addr) Network() string {
	return addr.Protocol
}

// Network ...
func (addr Addr) String() string {
	return net.JoinHostPort(addr.IP.String(), strconv.Itoa(addr.Port))
}

// PortUDP ...
func (addr Addr) UDP() *net.UDPAddr {
	return &net.UDPAddr{
		IP:   addr.IP,
		Port: addr.Port,
	}
}

// PortTCP ...
func (addr Addr) TCP() *net.TCPAddr {
	return &net.TCPAddr{
		IP:   addr.IP,
		Port: addr.Port,
	}
}

// JSON ...
func (s Service) JSON() []byte {
	marshal, err := json.Marshal(s)
	if err != nil {
		return nil
	}
	return marshal
}

// ParseService ...
func ParseService(data []byte) (service Service, err error) {
	err = json.Unmarshal(data, &service)
	return
}

func NewSource(service Service) Source {
	return &source{
		service: service,
	}
}

// String ...
func (s source) String() string {
	return s.addr.String()
}

// TryConnect ...
func (s *source) TryConnect() error {
	//remote := s.String()
	//localPort := LocalPort(s.Network(), s.mappingPort)
	//local := LocalAddr(localPort)
	//var dial net.Conn
	var err error
	//fmt.Println("ping", "local", local, "remote", remote, "network", s.Network(), "mapping", s.mappingPort)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := tryReverseUDP(s); err != nil {
			log.Errorw("tryReverseUDP|error", "error", err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		if err := tryReverseTCP(s); err != nil {
			log.Errorw("tryReverseTCP|error", "error", err)
			return
		}
	}()
	wg.Wait()
	//if s.mappingPort == localPort {
	//	dial, err = reuse.Dial(s.Network(), local, remote)
	//} else {
	//	if IsUDP(s.Network()) {
	//		udp, err := net.DialUDP(s.Network(), &net.UDPAddr{}, ParseUDPAddr(remote))
	//		if err != nil {
	//			return err
	//		}
	//		err = udp.SetDeadline(time.Now().Add(3 * time.Second))
	//		if err != nil {
	//			fmt.Println("debug|Ping|SetDeadline", err)
	//			return err
	//		}
	//		_, err = udp.Write([]byte("hello world"))
	//		if err != nil {
	//			fmt.Println("debug|Ping|Write", err)
	//			return err
	//		}
	//		data := make([]byte, maxByteSize)
	//		read, _, err := udp.ReadFromUDP(data)
	//		if err != nil {
	//			fmt.Println("debug|Ping|Read", err)
	//			return err
	//		}
	//		fmt.Println("received: ", string(data[:read]))
	//		return err
	//	}
	//	//dial, err = net.Dial(s.Network(), remote)
	//}

	//if err != nil {
	//	fmt.Println("debug|Ping|Dial", err)
	//	return err
	//}
	//_, err = dial.Write([]byte("hello world"))
	//if err != nil {
	//	fmt.Println("debug|Ping|Write", err)
	//	return err
	//}
	//data := make([]byte, maxByteSize)
	//read, err := dial.Read(data)
	//if err != nil {
	//	fmt.Println("debug|Ping|Read", err)
	//	return err
	//}
	//fmt.Println("received: ", string(data[:read]))
	return err
}

func tryReverseTCP(s *source) error {
	tcp, err := reuse.DialTCP("tcp", LocalTCPAddr(s.service.PortTCP), s.addr.TCP())
	if err != nil {
		log.Debugw("debug|tryReverse|DialTCP", err)
		return err
	}
	_, err = tcp.Write(s.service.JSON())
	if err != nil {
		log.Debugw("debug|tryReverse|Write", err)
		return err
	}
	data := make([]byte, maxByteSize)
	n, err := tcp.Read(data)
	if err != nil {
		log.Debugw("debug|tryReverse|ReadFromUDP", err)
		return err
	}
	log.Infow("received", "address", string(data[:n]))
	return nil
}

func tryReverseUDP(s *source) error {
	udp, err := net.DialUDP("udp", LocalUDPAddr(s.service.PortHole), s.addr.UDP())
	if err != nil {
		log.Debugw("debug|tryReverse|DialUDP", err)
		return err
	}

	_, err = udp.Write(s.service.JSON())
	if err != nil {
		log.Debugw("debug|tryReverse|Write", err)
		return err
	}
	data := make([]byte, maxByteSize)
	n, _, err := udp.ReadFromUDP(data)
	if err != nil {
		log.Debugw("debug|tryReverse|ReadFromUDP", err)
		return err
	}
	log.Infow("received", "address", string(data[:n]))
	return err
}

func tryUDP(addr *Addr) error {
	return nil
}

// ParseSourceAddr ...
func ParseSourceAddr(network string, ip net.IP, port int) *Addr {
	return &Addr{
		Protocol: network,
		IP:       ip,
		Port:     port,
	}
}
