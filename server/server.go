package server

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
var raddr *net.UDPAddr
var isUdp bool = true

func ListenAndServe(mode string, address string, port string) {

	addr := fmt.Sprintf("%s:%s", address, port)

	if mode == "tcp" {
		isUdp = false
		laddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			log.Fatalln(err)
			return
		}
		listenAndServeTcp(laddr)
	} else {
		laddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			log.Fatalln(err)
			return
		}
		listenAndServeUdp(laddr)
	}

}

func listenAndServeTcp(addr *net.TCPAddr) {

	log.Println("TCP Server Mode :", addr)

	l, err := net.ListenTCP("tcp", addr)

	if err != nil {
		log.Fatalln(err)
	}

	defer l.Close()

	conn, err := l.Accept()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("client connected")

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

func listenAndServeUdp(addr *net.UDPAddr) {

	log.Println("UDP Server Mode :", addr)

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	wg.Add(1)

	var n int

	go func() {

		var once sync.Once
		var buff [1024]byte

		for {
			n, raddr, _ = conn.ReadFromUDP(buff[:])
			fmt.Print(string(buff[:n]))
			once.Do(func() {
				wg.Done()
			})
		}
	}()

	wg.Wait()
	wg.Add(1)

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
			conn.WriteTo(buff, raddr)
		} else {
			writer.Write(buff)
			writer.Flush()
		}

	}

}
