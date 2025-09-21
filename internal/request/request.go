package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	initialized requestState = iota
	done
)

const bufferSize = 8

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case initialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = done
		return n, nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	request := &Request{
		state: initialized,
	}

	for request.state != done {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		n, err := reader.Read(buf[readToIndex:])

		if err != nil {
			if err == io.EOF {
				request.state = done
				break
			}
			return nil, fmt.Errorf("error reading to buffer: %s", err)
		}
		readToIndex += n

		bytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error parsing request: %s", err)
		}

		newBuf := buf[bytesParsed:]
		copy(buf, newBuf)
		readToIndex -= bytesParsed
	}

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, nil
	}
	stringRequest := string(data[:idx])
	splitLine := strings.Split(stringRequest, " ")

	// validate request line
	if len(splitLine) != 3 {
		return nil, 0, fmt.Errorf("malformed request line: %s", splitLine)
	}

	method := splitLine[0]
	for _, s := range splitLine[0] {
		if s < 'A' || s > 'Z' {
			return nil, 0, fmt.Errorf("method does not contain only capital letters: %s", splitLine[0])
		}
	}

	requestTarget := splitLine[1]

	versionParts := strings.Split(splitLine[2], "/")
	if len(versionParts) != 2 {
		return nil, 0, fmt.Errorf("malformed start-line: %s", splitLine[2])
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, 0, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	version := versionParts[1]
	if version != "1.1" {
		return nil, 0, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   version,
	}, idx + 2, nil
}
