package requests

import (
	"io"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}


func RequestFromReader(reader io.Reader) (*Request, error){

	rl := RequestLine{"1.1", "/", "GET"}
	r := Request{rl}
	return &r, nil
}
