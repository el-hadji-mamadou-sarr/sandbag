package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func backend() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from backend"))
	})
	go func() {
		http.ListenAndServe(":8081", nil)
	}()
}

func readRequest(buf *bufio.Reader, conn net.Conn, upstream string) (*http.Request, error) {
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	req, err := http.ReadRequest(buf)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	req = req.WithContext(ctx)
	go func() {
		<-req.Context().Done()
		cancel()
		conn.Close()
	}()
	req.URL.Scheme = "http"
	req.URL.Host = upstream
	req.Host = upstream
	req.RequestURI = ""
	return req, nil
}

func forwardRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Proxy-Version", "v1")

	return client.Do(req)
}

func shouldClose(req *http.Request, res *http.Response) bool {
	if req.Close || res.Close {
		return true
	}
	return false
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := bufio.NewReader(conn)
	upstreamClient := &http.Client{
		Timeout: 0,
	}
	for {
		req, err := readRequest(buf, conn, "localhost:8081")
		if err != nil {
			fmt.Println("error while reading request:", err)
			return
		}
		res, err := forwardRequest(upstreamClient, req)
		if err != nil {
			fmt.Println("error while forwarding response: ", err)
			return
		}
		if err := res.Write(conn); err != nil {
			return
		}
		if shouldClose(req, res) {
			return
		}
	}
}

func main() {
	go backend()
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}
