package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup
var udpRemoteAddr *net.UDPAddr
var isUdp bool = false

func ConnectAndWork(mode string, address string, port string) {

	addr := fmt.Sprintf("%s:%s", address, port)

	if mode == "udp" {
		isUdp = true
		udpLocalAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
		if err != nil {
			log.Fatalln(err)
			return
		}

		udpRemoteAddr, err = net.ResolveUDPAddr("udp", addr)
		if err != nil {
			log.Fatalln(err)
			return
		}
		listenAndServeUdp(udpLocalAddr)
	} else {
		//tcpLocalAddr := net.ResolveTCPAddr("tcp", addr)
		tcpRemoteAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			log.Fatalln(err)
			return
		}
		connectTcpAndWork(tcpRemoteAddr)
	}

}

func connectTcpAndWork(addr *net.TCPAddr) {

	log.Println("Connecting to ", addr)

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected")

	wg.Add(1)

	go func(c net.Conn) {

		defer c.Close()
		defer wg.Done()

		reader := bufio.NewReader(c)

		for {

			buff, _, err := reader.ReadLine()
			if err == io.EOF {
				log.Println("client closed")
				return
			} else if err != nil {
				log.Fatalln("err", err)
				return
			}

			line := string(buff)

			fmt.Println(line)
		}

	}(conn)

	go interact(conn)

	wg.Wait()

}

func listenAndServeUdp(laddr *net.UDPAddr) {

	log.Println("UDP Listening :", laddr)

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln(err)
	}

	wg.Add(1)

	var n int

	go func() {

		var buff [1024]byte

		for {
			n, _, _ = conn.ReadFromUDP(buff[:])
			fmt.Print(string(buff[:n]))
		}
	}()

	go interact(conn)

	wg.Wait()

}

func interact(c net.Conn) {

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(c)

	var conn *net.UDPConn

	if isUdp {
		conn = c.(*net.UDPConn)
	}

	for {

		buff, _, err := reader.ReadLine()
		if err == io.EOF {
			log.Println("bye~")
			return
		}

		buff = append(buff, []byte("\r\n")...)

		if isUdp {
			conn.WriteTo(buff, udpRemoteAddr)
		} else {
			writer.Write(buff)
			writer.Flush()
		}

	}

}
