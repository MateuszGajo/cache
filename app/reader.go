package main

import (
	"net"
)

type Reader struct {
	parser Parser
	conn   net.Conn
}

func NewReader(conn net.Conn) *Reader {
	return &Reader{
		parser: *NewParser(),
		conn:   conn,
	}
}

func (reader *Reader) readInput() ([]ParseResultRecord, error) {
	buf := make([]byte, 1024)
	n, err := reader.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	input := string(buf[:n])

	if len(input) == 0 {
		return nil, nil
	}
	parsed := reader.parser.parseStream(input)

	if parsed.err != nil {
		return nil, parsed.err
	}

	return parsed.records, nil
}
