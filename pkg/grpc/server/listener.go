package server

import (
	"crypto/tls"
	"net"
)

// NewListener creates a new instance of a TCP listener.
func NewListener(addr string, tlsConfig ...*tls.Config) (net.Listener, error) {
	var (
		listener net.Listener
		err      error
	)

	if len(tlsConfig) > 0 {
		listener, err = tls.Listen("tcp", addr, tlsConfig[0])
	} else {
		listener, err = net.Listen("tcp", addr)
	}

	if err != nil {
		return nil, err
	}

	return listener, nil
}
