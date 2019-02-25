package client

const (
	// serverDomain needs to be changed to your onion address
	// also works with ip-address and domain
	serverDomain = "youronionadresshere.onion"
	serverPort   = ":1337"
	serverAddr   = serverDomain + serverPort
)

// serverCert needs to be changed to the TLS certificate of the server
// intendation breaks the certificate
const serverCert = `-----BEGIN CERTIFICATE-----
____ CERTIFICATE GOES HERE | DO NOT INDENT ____
-----END CERTIFICATE-----`
