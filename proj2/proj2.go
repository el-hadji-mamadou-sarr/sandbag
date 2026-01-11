package proj2

import (
	"io"
	"log"
	"net"
	"sync"
)

func proxyConn(client net.Conn, upstreamAddr string) {
	defer client.Close()
	var wg sync.WaitGroup
	upstream, err := net.Dial("tcp", upstreamAddr)
	if err != nil {
		return
	}
	defer upstream.Close()
	wg.Go(func() {
		io.Copy(client, upstream)
		if tcp, ok := client.(*net.TCPConn); ok {
			tcp.Close()
		}
	})
	wg.Go(func() {
		io.Copy(upstream, client)
		if tcp, ok := upstream.(*net.TCPConn); ok {
			tcp.Close()
		}
	})
	wg.Wait()
}

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal()
	}

	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go proxyConn(conn, "localhost:9001")
	}

}
