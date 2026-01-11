package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
)

func backend() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from backend"))
	})
	http.ListenAndServe(":8081", nil)
}

func readRequest(conn net.Conn, upstream string) (*http.Request, error) {
	buf := bufio.NewReader(conn)
	req, err := http.ReadRequest(buf)
	if err != nil {
		return nil, err
	}
	req.URL.Scheme = "http"
	req.URL.Host = upstream
	req.Method = "GET"
	req.RequestURI = ""
	return req, nil
}

func forwardRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Proxy-Version", "v1")

	client := &http.Client{}
	return client.Do(req)
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	req, err := readRequest(conn, "localhost:8081")
	if err != nil {
		fmt.Println("error while reading request:", err)
		return
	}
	res, err := forwardRequest(req)
	if err != nil {
		fmt.Println("error while forwarding response: ", err)
		return
	}
	defer res.Body.Close()
	res.Write(conn)
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
