package signalr

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type clientSSEConnection struct {
	ConnectionBase
	reqURL    string
	sseReader io.Reader
	sseWriter io.Writer
}

func newClientSSEConnection(parentContext context.Context, address string, connectionID string, body io.ReadCloser) (*clientSSEConnection, error) {
	// Setup request
	reqUrl, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	q := reqUrl.Query()
	q.Set("id", connectionID)
	reqUrl.RawQuery = q.Encode()
	c := clientSSEConnection{
		ConnectionBase: ConnectionBase{
			ctx:          parentContext,
			connectionID: connectionID,
		},
		reqURL: reqUrl.String(),
	}
	c.sseReader, c.sseWriter = io.Pipe()
	go func() {
		p := make([]byte, 1<<15)
	loop:
		for {
			n, err := body.Read(p)
			if err != nil {
				break loop
			}
			lines := strings.Split(string(p[:n]), "\n")
			for _, line := range lines {
				line = strings.Trim(line, "\r\t ")
				// Ignore everything but data
				if strings.Index(line, "data:") != 0 {
					continue
				}
				json := strings.Replace(strings.Trim(line, "\r"), "data:", "", 1)
				// Spec says: If it starts with Space, remove it
				if len(json) > 0 && json[0] == ' ' {
					json = json[1:]
				}
				_, err = c.sseWriter.Write([]byte(json))
				if err != nil {
					break loop
				}
			}
		}
		_ = body.Close()
	}()
	return &c, nil
}

func (c *clientSSEConnection) Read(p []byte) (n int, err error) {
	return c.sseReader.Read(p)
}

func (c *clientSSEConnection) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest("POST", c.reqURL, bytes.NewReader(p))
	if err != nil {
		return 0, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("POST %v -> %v", c.reqURL, resp.Status)
	}
	_ = resp.Body.Close()
	return len(p), err
}
