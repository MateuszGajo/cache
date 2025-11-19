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

func writeMany(conn net.Conn, data [][]byte) {
	for _, item := range data {
		write(conn, item)

		_, err := readInput(conn)
		if err != nil {
			panic(err)
		}
	}
}

func TestCreateStream(t *testing.T) {
	server, conn := startServer()

	write(conn, []byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "1-1", "foo", "bar"})))
	RESPParsed := readSingleLineResponse(conn, t)
	addStreamResp := RESPParsed.data.([]string)[0]

	write(conn, []byte(BuildPrimitiveRESPArray([]string{"type", "key"})))
	RESPParsed = readSingleLineResponse(conn, t)
	dataTypeResp := RESPParsed.data.([]string)[0]

	if dataTypeResp != StreamType {
		t.Errorf("Expected type to be: %v, got: %v", StreamType, dataTypeResp)
	}
	if addStreamResp != "1-1" {
		t.Errorf("expected steam id to be 1-1, got %v", addStreamResp)
	}

	cleanup(server, conn)

}

func TestCreateStreamWithDuplicate(t *testing.T) {
	server, conn := startServer()

	write(conn, []byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "1-1", "foo", "bar"})))
	RESPParsed := readSingleLineResponse(conn, t)
	addStreamResp := RESPParsed.data.([]string)[0]

	write(conn, []byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "1-1", "foo", "bar"})))
	RESPParsed = readSingleLineResponse(conn, t)
	addDuplicateStreamResp := RESPParsed.data.([]string)[0]

	if addStreamResp != "1-1" {
		t.Errorf("expected streamId to be: %q, got:%q", "1-1", addStreamResp)
	}
	if addDuplicateStreamResp != "-ERR The ID specified in XADD is equal or smaller than the target stream top item" {
		t.Errorf("expected: %q, got:%q", "-ERR The ID specified in XADD is equal or smaller than the target stream top item",
			addStreamResp)
	}

	cleanup(server, conn)

}

func TestAutoGenerateStreamSequenceNumber(t *testing.T) {
	server, conn := startServer()

	write(conn, []byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "1-*", "foo", "bar"})))
	RESPParsed := readSingleLineResponse(conn, t)
	addFirstStreamResp := RESPParsed.data.([]string)[0]

	write(conn, []byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "1-*", "foo", "bar"})))
	RESPParsed = readSingleLineResponse(conn, t)
	addSecondStreamResp := RESPParsed.data.([]string)[0]

	if addFirstStreamResp != "1-0" {
		t.Errorf("expected: %q, got:%q", "$3\r\n1-0\r\n", addFirstStreamResp)
	}
	if addSecondStreamResp != "1-1" {
		t.Errorf("expected: %q, got:%q", "$3\r\n1-1\r\n", addFirstStreamResp)
	}

	cleanup(server, conn)
}

func TestAutoGenerateSequenceNumberWith0Prefix(t *testing.T) {
	server, conn := startServer()

	write(conn, []byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "0-*", "foo", "bar"})))
	RESPParsed := readSingleLineResponse(conn, t)
	addFirstStreamResp := RESPParsed.data.([]string)[0]

	write(conn, []byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "1-*", "foo", "bar"})))
	RESPParsed = readSingleLineResponse(conn, t)
	addSecondStreamResp := RESPParsed.data.([]string)[0]

	if addFirstStreamResp != "0-1" {
		t.Errorf("expected: %q, got:%q", "$3\r\n0-1\r\n", addFirstStreamResp)
	}
	if addSecondStreamResp != "1-0" {
		t.Errorf("expected: %q, got:%q", "$3\r\n1-0\r\n", addFirstStreamResp)
	}

	cleanup(server, conn)
}

func TestAutoGenerateSequenceAndTime(t *testing.T) {
	server, conn := startServer()

	write(conn, []byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "*", "foo", "bar"})))
	RESPParsed := readSingleLineResponse(conn, t)
	addStreamResp := RESPParsed.data.([]string)[0]

	responseSplitted := strings.Split(addStreamResp, "-")
	respTime := responseSplitted[0]
	respSequence := responseSplitted[1]
	currentTimeMili := time.Now().UnixMilli()
	timeTwoSecAgoMili := time.Now().Add(-2 * time.Second).UnixMilli()

	if strconv.FormatInt(currentTimeMili, 10) < respTime || strconv.FormatInt(timeTwoSecAgoMili, 10) > respTime {
		t.Errorf("time should be a bit before current and less than 2 seconds ago, current time:%v, two hours ago:%v, got time:%v", currentTimeMili, timeTwoSecAgoMili, respTime)
	}
	if respSequence != "0" {
		t.Errorf("expected sequence to be: %q, got:%q", "0", respSequence)
	}

	cleanup(server, conn)
}

func TestReadingStreamRange(t *testing.T) {
	server, conn := startServer()
	streams := [][]byte{
		[]byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "0-1", "foo", "bar"})),
		[]byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "0-2", "foo", "bar"})),
		[]byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "0-3", "foo", "bar"})),
		[]byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "0-4", "foo", "bar"})),
	}

	writeMany(conn, streams)

	write(conn, []byte(BuildPrimitiveRESPArray([]string{"xrange", "key", "0-2", "0-4"})))
	RESPParsed := readSingleLineResponse(conn, t)
	respData := RESPParsed.data.([]string)
	fmt.Println(respData)

	// fmt.Printf("%+v\n", respData)
	// t.Error("fds")
	//Fix parser first

	// if string(respData) != "*3\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\n0-3\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\n0-4\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" {
	// 	t.Fatalf("We are expeciting, %q, and got:%q", "3\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\n0-3\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\n0-4\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", string(respData))
	// }

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

	_, err = conn.Write([]byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "0-1", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "0-2", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildPrimitiveRESPArray([]string{"xrange", "key", "-", "0-2"})))
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

	_, err = conn.Write([]byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "0-1", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildPrimitiveRESPArray([]string{"xadd", "key", "0-2", "foo", "bar"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildPrimitiveRESPArray([]string{"xrange", "key", "0-0", "+"})))
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

	_, err = conn.Write([]byte(BuildPrimitiveRESPArray([]string{"xadd", "testStream", "0-1", "temperature", "96"})))
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = readInput(conn)
	if err != nil {
		t.Errorf("got err: %q", err)
	}

	_, err = conn.Write([]byte(BuildPrimitiveRESPArray([]string{"xread", "streams", "testStream", "0-0"})))
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

	write(conn, []byte(BuildPrimitiveRESPArray([]string{"rpush", "newList", "aa"})))
	RESPParsed := readSingleLineResponse(conn, t)
	data := RESPParsed.data.(RESPIntData)

	if data != 1 {
		t.Errorf("Expected data to be 1, got :%v", data)
	}

	cleanup(server, conn)

}
