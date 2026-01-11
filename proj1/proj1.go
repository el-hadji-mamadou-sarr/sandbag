package proj1

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func worker(ctx context.Context, id int, job <-chan string) {
	for {

		select {
		case <-ctx.Done():
			fmt.Println("worker", id, "closed")
			return
		case msg, ok := <-job:
			if !ok {
				fmt.Println("worker", id, "channel is closed")
				return
			}
			fmt.Println("worker id: ", id, "message: ", msg)
		}

	}
}

func test_channels() {
	ctx, cancel := context.WithCancel(context.Background())

	jobs := make(chan string, 10)
	for i := 1; i <= 3; i++ {
		go worker(ctx, i, jobs)
	}
	go func() {
		jobs <- "a"
		jobs <- "b"
		jobs <- "c"
		close(jobs)
	}()

	time.Sleep(2 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
}

type Handler interface {
	serverHandler(http.ResponseWriter, *http.Request)
	reverseProxyHandler(http.ResponseWriter, *http.Request)
}

type MyHandler struct {
}

func (h *MyHandler) serverHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("response from backend profile"))
}
func (h *MyHandler) reverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request arrived to the proxy server")
	backendUrl := "http://localhost:8080" + "/backend" + r.URL.Path
	fmt.Println(backendUrl)
	ctx, cancel := context.WithTimeout(r.Context(), 100*time.Millisecond)
	defer cancel()
	request, error := http.NewRequestWithContext(ctx, r.Method, backendUrl, r.Body)
	if error != nil {
		fmt.Println("error occured")
		http.Error(w, "internal error occured", http.StatusInternalServerError)
		return
	}
	request.Header = r.Header
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("error occured while calling the backend")
		http.Error(w, "backend unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	fmt.Println("backend responded with status", resp.StatusCode)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func run_server() {
	serverHandler := MyHandler{}
	http.HandleFunc("/backend/profile", serverHandler.serverHandler)
	http.ListenAndServe(":8080", nil)
}

func run_reverse_proxy() {
	serverHandler := MyHandler{}
	http.HandleFunc("/profile", serverHandler.reverseProxyHandler)
	http.ListenAndServe(":9090", nil)
}
func rev() {
	for {
		go run_server()
		go run_reverse_proxy()
		select {}
	}
}

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func (s *Store) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

func (s *Store) Del(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.data[key]; !ok {
		return false
	}
	delete(s.data, key)
	return true
}

func handleSplit(line string) []string {
	resp := strings.Fields(line)
	return resp
}

const maxFrameSize = 1024

func readFrame(conn net.Conn) ([]byte, error) {
	var lenBuf [4]byte
	// trying to read the first 4 bytes and store in length
	if _, err := io.ReadFull(conn, lenBuf[:]); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenBuf[:])

	if length > maxFrameSize {
		fmt.Println("the payload is too big")
		conn.Close()
		return nil, errors.New("the payload is too big")
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func writeFrame(conn net.Conn, payload []byte) error {
	if err := binary.Write(conn, binary.BigEndian, uint32(len(payload))); err != nil {
		return err
	}
	_, err := conn.Write(payload)
	return err
}

func handleConn(ctx context.Context, store *Store, conn net.Conn) {
	defer conn.Close()
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		buf, err := readFrame(conn)
		if err != nil {
			if me, ok := err.(net.Error); ok && me.Timeout() {
				writeFrame(conn, []byte("ERROR timeout"))
			}
			return
		}

		line := string(buf)
		resp := "ERROR unknown command"
		badArg := "BAD ARGUMENTS"
		if line == "PING" {
			resp = "PONG"
		}
		splitLine := handleSplit(line)
		if len(splitLine) == 0 {
			continue
		}
		action := splitLine[0]
		switch action {
		case "GET":
			if len(splitLine) > 1 {
				value, ok := store.Get(splitLine[1])
				if !ok {
					resp = "NOT_FOUND"
				} else {
					resp = "VALUE " + value
				}
			} else {
				resp = badArg
			}
		case "SET":
			if len(splitLine) > 2 {
				store.Set(splitLine[1], splitLine[2])
				resp = "OK"
			} else {
				resp = badArg
			}
		case "DEL":
			if len(splitLine) > 1 {

				ok := store.Del(splitLine[1])
				if !ok {
					resp = "NOT_FOUND"
				} else {
					resp = "OK"
				}
			} else {
				resp = badArg
			}
		}
		conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
		respBuf := []byte(resp)
		if err := writeFrame(conn, respBuf); err != nil {
			return
		}

	}
}
func server(ctx context.Context) {
	var wg sync.WaitGroup
	data := make(map[string]string)
	store := Store{data: data}
	connLimit := make(chan struct{}, 5)
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-ctx.Done()
		ln.Close() // stop listening
	}()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					wg.Wait() // wait connexions finishes
					fmt.Println("server shut down cleanly")
					return
				default:
					continue
				}
			}
			select {
			case connLimit <- struct{}{}:
			default:
				conn.Write([]byte("server busy"))
				conn.Close()
				continue
			}

			wg.Go(func() {
				defer func() { <-connLimit }()
				handleConn(ctx, &store, conn)
			})
		}
	}()
}
func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	server(ctx)
	time.Sleep(100 * time.Millisecond)
	var wg sync.WaitGroup

	wg.Go(func() {
		conn, err := net.Dial("tcp", "localhost:9000")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		writeFrame(conn, []byte("PING"))
		writeFrame(conn, []byte("SET hello world"))
		writeFrame(conn, []byte("DEL hello"))

		for range 3 {
			data, err := readFrame(conn)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			fmt.Println("echo: ", string(data))
		}
	})
	wg.Wait()
}
