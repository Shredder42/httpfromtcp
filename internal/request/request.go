package request

import (
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("Error reading request: %s", err)
	}

	requestLine, err := parseRequestLine(request)
	if err != nil {
		return nil, err
	}

	newRequest := Request{
		RequestLine: *requestLine,
	}

	return &newRequest, nil
}

func parseRequestLine(request []byte) (*RequestLine, error) {
	stringRequest := string(request)
	lines := strings.Split(stringRequest, "\r\n")
	requestLine := lines[0]
	splitLine := strings.Split(requestLine, " ")

	// validate request line
	if len(splitLine) != 3 {
		return nil, fmt.Errorf("malformed request line: %s", splitLine)
	}

	for _, s := range splitLine[0] {
		if !unicode.IsUpper(s) {
			return nil, fmt.Errorf("method does not contain only capital letters: %s", splitLine[0])
		}
	}

	versionParts := strings.Split(splitLine[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", splitLine[2])
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        splitLine[0],
		RequestTarget: splitLine[1],
		HttpVersion:   version,
	}, nil
}
