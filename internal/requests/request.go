package requests

import (
	"errors"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ErrReadingMessage = errors.New("can not read the message")
var ErrParsingMessage = errors.New("can not parse the message")
var ErrHttpVersion = errors.New("unsupported http protocol version")
var ErrMethod = errors.New("unsupported method")

const SEPARATOR = "\r\n"

var methods = []string{
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"DELETE",
	"CONNECT",
	"OPTIONS",
	"TRACE",
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, ErrReadingMessage
	}
	str := string(data)
	requestLine, _, err := parseRequestLine(str)
	if err != nil {
		return nil, err
	}
	request := Request{*requestLine}
	return &request, nil
}

func parseRequestLine(data string) (*RequestLine, string, error) {
	idx := strings.Index(data, SEPARATOR)
	if idx == -1 {
		return nil, "", ErrParsingMessage
	}
	requestLine := data[:idx]
	restOfMessage := data[idx+len(SEPARATOR):]

	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, restOfMessage, ErrParsingMessage
	}

	method := parts[0]
	requestTarget := parts[1]
	protocol := parts[2]

	if !slices.Contains(methods, method) {
		return nil, restOfMessage, ErrMethod
	}

	if protocol != "HTTP/1.1" {
		return nil, restOfMessage, ErrHttpVersion
	}

	rl := RequestLine{"1.1", requestTarget, method}

	return &rl, restOfMessage, nil

}
