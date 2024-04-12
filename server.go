package main

import (
	"errors"
	"net"
	"strconv"
	"time"
)

type method string

type server struct {
	addr              string
	port              int64
	ln                net.Listener
	read_timeout      time.Duration
	total_connections int
	methods           map[method]int
}

type response struct {
	message string
	code    uint
}

const (
	SERVER_METHOD_GET method = "GET"
	SERVER_METHOD_SET method = "SET"
	SERVER_METHOD_DEL method = "DEL"
)

func (s *server) serve() error {
	remoteHost := net.JoinHostPort(s.addr, strconv.FormatInt(s.port, 10))
	listener, err := net.Listen("tcp", remoteHost)

	if err != nil {
		return err
	}

	s.ln = listener
	return nil
}

func NewServer(
	addr string,
	port int64,
	read_timeout time.Duration,
) (*server, error) {
	server := &server{}
	server.addr = addr
	server.port = port
	server.read_timeout = read_timeout
	server.total_connections = 0

	err := server.serve()

	return server, err
}

func (s server) writeResponse(c net.Conn, res response) {
	c.Write([]byte(res.message))
}

func (s server) handleConnection(c net.Conn, store *store) error {
	defer c.Close()

	deadline := time.Now().Add(s.read_timeout * time.Second)
	c.SetReadDeadline(deadline)

	buffer := make([]byte, 1024)

	var total []byte

	// Log("Processing new connection...")

	errorResponse := response{"SERVER_ERROR", 500}

	for {
		// Log("Reading Request into Buffer")
		n, err := c.Read(buffer)

		if err != nil || n == 0 {
			s.writeResponse(c, errorResponse)
			break
		}

		is_end := isEnd(buffer[:n])
		total = append(total, buffer[:n]...)

		if !is_end {
			continue
		}

		// Log("Parsing Buffer...")

		parser := parser{}
		parser.parse(total).method().subject().value()

		if parser.error != nil {
			errorResponse := response{parser.error.Error(), 500}
			s.writeResponse(c, errorResponse)
			break
		}

		var res response

		switch parser.payload.method {

		case SERVER_METHOD_SET:
			store.set(parser.payload.subject, parser.payload.value)
			res = response{"OK", 201}
		case SERVER_METHOD_GET:
			node, err := store.get(parser.payload.subject)
			if err != nil {
				res = response{"NOT_FOUND", 404}
			} else {
				res = response{node.value, 200}
			}
		case SERVER_METHOD_DEL:
			_, err := store.delete(parser.payload.subject)
			if err != nil {
				res = response{"NOT_FOUND", 404}
			} else {
				res = response{"OK", 200}
			}
		default:
			s.writeResponse(c, response{"SERVER_ERROR", 500})
			return errors.New("Unable to handle request")
		}

		s.writeResponse(c, res)
		return nil
	}
	return nil
}
