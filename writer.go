package quiclog

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"unsafe"
)

type writer struct {
	addr  string
	index string
	typ   string
	url   string
}

var defalutClient *http.Client

func init() {
	defalutClient = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10000,
			MaxIdleConnsPerHost: 10000,
			IdleConnTimeout:     600 * time.Second,
			Dial: func(netw, addr string) (net.Conn, error) {
				return net.DialTimeout(netw, addr, time.Second*5)
			},

			ResponseHeaderTimeout: 5 * time.Second,
		},
	}

	if addr := os.Getenv("ES_ADDR"); addr != "" {
		hostname, _ := os.Hostname()
		w := &writer{
			index: "quic-log",
			typ:   hostname,
			addr:  addr,
			url:   addr + "/quic-log/" + hostname + "/",
		}

		log.SetOutput(w)
	}
}

type Body struct {
	Content   string `json:"content,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

func (w *writer) Write(p []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("nil")
	}

	body := &Body{
		Content:   bytes2str(p[27:]),
		Timestamp: bytes2str(p[:19]),
	}

	b, _ := json.Marshal(body)

	resp, err := defalutClient.Post(w.url, "application/json", bytes.NewReader(b))
	if err != nil {
		return 0, errors.New("resp")
	}

	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}

	return len(p), err
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
