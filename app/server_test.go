package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

//write handshake test

func readSingleLineResponse(conn net.Conn, t *testing.T) RESPRecord {
	RESPParsed, err := readInput(conn)

	if err != nil {
		t.Error(err)
	}

	return RESPParsed.records[0]
}

func startServer() (*Server, net.Conn) {
	var server *Server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server = NewServer(WithPort("6379"))
	}()

	wg.Wait()

	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(fmt.Sprintf("connection err: %q", err))
	}

	return server, conn
}

func cleanup(server *Server, conn net.Conn) {
	m = map[string]CustomSetStore{}
	conn.Close()
	server.Close()
}

func write(conn net.Conn, data []byte) {
	_, err := conn.Write(data)
	if err != nil {
		// dont need to handle them gracefully, it will be tested separately, we only need simple interface
		panic(fmt.Sprintf("write err: %q", err))
	}
}

// TODO rewrite all test
// func TestCreateStream(t *testing.T) {
// 	server, conn := startServer()

// 	write(conn, []byte(BuildRESPArray([]string{"xadd", "key", "1-1", "foo", "bar"})))
// 	RESPParsed := readSingleLineResponse(conn, t)
// 	streamId := RESPParsed.data.([]string)[0]

// 	write(conn, []byte(BuildRESPArray([]string{"type", "key"})))
// 	RESPParsed = readSingleLineResponse(conn, t)
// 	dataType := RESPParsed.data.([]string)[0]

// 	if dataType != StreamType {
// 		t.Errorf("Expected type to be: %v, got: %v", StreamType, dataType)
// 	}
// 	if streamId != "1-1" {
// 		t.Errorf("expected steam id to be 1-1, got %v", streamId)
// 	}

// 	cleanup(server, conn)

// }

func TestG2(t *testing.T) {
	var server *Server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server = NewServer(WithPort("6379"))
	}()

	wg.Wait()

	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "1-1", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	RESPParsed, err := readInput(conn)
	data := RESPParsed.records[0].data.([]string)[0]
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	if data != "1-1" {
		t.Errorf("expected: %q, got:%q", "$3\r\n1-1\r\n", data)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "1-1", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	RESPParsed, err = readInput(conn)
	data = RESPParsed.records[0].data.([]string)[0]
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	if data != "-ERR The ID specified in XADD is equal or smaller than the target stream top item" {
		t.Errorf("expected: %q, got:%q", "-ERR The ID specified in XADD is equal or smaller than the target stream top item",
			data)
	}

	t.Cleanup(func() {
		m = map[string]CustomSetStore{}
		conn.Close()
		server.Close()
	})

}

func TestG3(t *testing.T) {
	var server *Server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server = NewServer(WithPort("6379"))
	}()

	wg.Wait()

	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "1-*", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	RESPParsed, err := readInput(conn)
	data := RESPParsed.records[0].data.([]string)[0]
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	if data != "1-0" {
		t.Errorf("expected: %q, got:%q", "$3\r\n1-0\r\n", data)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "1-*", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	RESPParsed, err = readInput(conn)
	data = RESPParsed.records[0].data.([]string)[0]
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	if data != "1-1" {
		t.Errorf("expected: %q, got:%q", "$3\r\n1-1\r\n", data)
	}

	t.Cleanup(func() {
		m = map[string]CustomSetStore{}
		conn.Close()
		server.Close()
	})

}

func TestG4(t *testing.T) {
	var server *Server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server = NewServer(WithPort("6379"))
	}()

	wg.Wait()
	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "0-*", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	RESPParsed, err := readInput(conn)
	data := RESPParsed.records[0].data.([]string)[0]
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	if data != "0-1" {
		t.Errorf("expected: %q, got:%q", "$3\r\n0-1\r\n", data)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "1-*", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	RESPParsed, err = readInput(conn)
	data = RESPParsed.records[0].data.([]string)[0]
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	if data != "1-0" {
		t.Errorf("expected: %q, got:%q", "$3\r\n1-0\r\n", data)
	}

	t.Cleanup(func() {
		m = map[string]CustomSetStore{}
		conn.Close()
		server.Close()
	})

}

func TestG5(t *testing.T) {
	var server *Server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server = NewServer(WithPort("6379"))
	}()

	wg.Wait()
	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "*", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	RESPParsed, err := readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}
	data := RESPParsed.records[0].data.([]string)[0]

	response := strings.Split(data, "-")
	timee := response[0]
	sequence := response[1]
	currentTimeMili := time.Now().UnixMilli()
	timeTwoHoursAgoMili := time.Now().Add(-2 * time.Hour).UnixMilli()

	if strconv.FormatInt(currentTimeMili, 10) < timee || strconv.FormatInt(timeTwoHoursAgoMili, 10) > timee {
		t.Errorf("time should be a bit before current and less than 2 hours ago, current time:%v, two hours ago:%v, got time:%v", currentTimeMili, timeTwoHoursAgoMili, timee)
	}

	if sequence != "0" {
		t.Errorf("expected sequence to be: %q, got:%q", "0", sequence)
	}

	t.Cleanup(func() {
		conn.Close()
		m = map[string]CustomSetStore{}
		server.Close()
	})

}

func TestG6(t *testing.T) {
	var server *Server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server = NewServer(WithPort("6379"))
	}()

	wg.Wait()

	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "0-1", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "0-2", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "0-3", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "0-4", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xrange", "key", "0-2", "0-4"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Errorf("got err: %q", err)
	}
	if string(buf[:n]) != "*3\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\n0-3\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\n0-4\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" {
		t.Fatalf("We are expeciting, %q, and got:%q", "3\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\n0-3\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\n0-4\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", string(buf[:n]))
	}

	t.Cleanup(func() {
		conn.Close()
		m = map[string]CustomSetStore{}
		server.Close()
	})

}

func TestG7(t *testing.T) {
	var server *Server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server = NewServer(WithPort("6379"))
	}()

	wg.Wait()

	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "0-1", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "0-2", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xrange", "key", "-", "0-2"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Errorf("got err: %q", err)
	}
	if string(buf[:n]) != "*2\r\n*2\r\n$3\r\n0-1\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" {
		t.Fatalf("We are expeciting, %q, and got:%q", "*2\r\n*2\r\n$3\r\n0-1\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", string(buf[:n]))
	}

	t.Cleanup(func() {
		conn.Close()
		m = map[string]CustomSetStore{}
		server.Close()
	})

}

func TestG8(t *testing.T) {
	var server *Server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server = NewServer(WithPort("6379"))
	}()

	wg.Wait()

	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "0-1", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "key", "0-2", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xrange", "key", "0-0", "+"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Errorf("got err: %q", err)
	}
	if string(buf[:n]) != "*2\r\n*2\r\n$3\r\n0-1\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" {
		t.Fatalf("We are expeciting, %q, and got:%q", "*2\r\n*2\r\n$3\r\n0-1\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", string(buf[:n]))
	}

	t.Cleanup(func() {
		conn.Close()
		m = map[string]CustomSetStore{}
		server.Close()
	})

}

func TestG9(t *testing.T) {
	var server *Server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server = NewServer(WithPort("6379"))
	}()

	wg.Wait()

	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xadd", "testStream", "0-1", "temperature", "96"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"xread", "streams", "testStream", "0-0"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Errorf("got err: %q", err)
	}
	if string(buf[:n]) != "*1\r\n*2\r\n$10\r\ntestStream\r\n*1\r\n*2\r\n$3\r\n0-1\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n96\r\n" {
		t.Fatalf("We are expeciting, %q, and got:%q", "*1\r\n*2\r\n$10\r\ntestStream\r\n*1\r\n*2\r\n$3\r\n0-1\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n96\r\n", string(buf[:n]))
	}

	t.Cleanup(func() {
		conn.Close()
		m = map[string]CustomSetStore{}
		server.Close()
	})

}

func TestRPush(t *testing.T) {
	server, conn := startServer()

	write(conn, []byte(BuildRESPArray([]string{"rpush", "newList", "aa"})))
	RESPParsed := readSingleLineResponse(conn, t)
	data := RESPParsed.data.(RESPIntData)

	if data != 1 {
		t.Errorf("Expected data to be 1, got :%v", data)
	}

	cleanup(server, conn)

}
