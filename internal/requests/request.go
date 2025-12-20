package requests

import (
	"bytes"
	"errors"
	"io"
)

type parserState string

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ErrReadingMessage = errors.New("can not read the message")
var ErrParsingMessage = errors.New("can not parse the message")
var ErrHttpVersion = errors.New("unsupported http protocol version")
var ErrRequestInErrorState = errors.New("request in error state")
var ErrMethod = errors.New("unsupported method")

var SEPARATOR = []byte("\r\n")

var methods = [][]byte{
	[]byte("GET"),
	[]byte("HEAD"),
	[]byte("POST"),
	[]byte("PUT"),
	[]byte("PATCH"),
	[]byte("DELETE"),
	[]byte("CONNECT"),
	[]byte("OPTIONS"),
	[]byte("TRACE"),
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}
func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.state = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			read += n
			r.RequestLine = *rl

			r.state = StateDone

		case StateDone:
			break outer
		case StateError:
			return 0, ErrRequestInErrorState
		}

	}

	return read, nil

}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n

		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])

		bufLen -= readN
	}

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, SEPARATOR)
	parsedBytes := 0
	if idx == -1 {
		return nil, parsedBytes, nil
	}
	requestLine := data[:idx]
	parsedBytes = idx + len(SEPARATOR)

	parts := bytes.Split(requestLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrParsingMessage
	}

	method := parts[0]
	requestTarget := parts[1]
	protocol := parts[2]

	if !containByteSlice(methods, method) {
		return nil, 0, ErrMethod
	}

	protocolParts := bytes.Split(protocol, []byte("/"))
	if len(protocolParts) != 2 || !bytes.Equal(protocolParts[0], []byte("HTTP")) || !bytes.Equal(protocolParts[1], []byte("1.1")) {
		return nil, parsedBytes, ErrHttpVersion
	}
	httpVersion := protocolParts[1]

	rl := RequestLine{string(httpVersion), string(requestTarget), string(method)}

	return &rl, parsedBytes, nil

}

func containByteSlice(slices [][]byte, target []byte) bool {

	for _, slice := range slices {
		if bytes.Equal(slice, target) {
			return true
		}
	}

	return false
}
