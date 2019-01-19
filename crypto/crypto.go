package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"log"
	mathrand "math/rand"
)

// SetHostname Sets the Hostname of the machine to a
// random string with the length of 16, encrypts the outcome and
// writes it to Disk
func SetHostname(path string) error {
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	length := 16
	hostname := make([]byte, length)
	for i := 0; i < length; i++ {
		num := mathrand.Intn(len(all))
		hostname[i] = all[num]
	}
	return EnctoFile(hostname, path)
}

func GetHostname(path string) []byte {
	log.Println("getHostname")
	content, err := ioutil.ReadFile(path)
	if err != nil {
		if SetHostname(path) == nil {
			content, err = ioutil.ReadFile(path)
			if err != nil {
				return nil
			}
			return content
		}
		return nil
	}
	return content
}

func loadRsaKey() (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyRsa))
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("Key is not a *rsa.PublicKey")
	}
}

func encRsa(data []byte) []byte {
	rand := rand.Reader
	RsaPublicKey, err := loadRsaKey()
	if err != nil {
		log.Println("Error loading RsaKey", err)
		return nil
	}
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand, RsaPublicKey, data, nil)
	if err != nil {
		log.Println("Error from encryption:", err)
		return nil
	}
	return ciphertext
}

func EnctoFile(data []byte, path string) error {
	aeskey, err := genAesKey()
	if err != nil {
		return err
	}
	encKey := encRsa(aeskey)
	encData, err := encAes(data, aeskey)
	if err != nil {
		return err
	}
	enc := append(encKey, encData...)
	err = ioutil.WriteFile(path, enc, 0666)
	if err != nil {
		return err
	}
	return nil

}

// genAesKey generates a 256bit AES Key
func genAesKey() ([]byte, error) {
	AesKey := make([]byte, 32)
	_, err := rand.Read(AesKey)
	if err != nil {
		return nil, err
	}
	return AesKey, nil
}

func encAes(data []byte, AesKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(AesKey)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	encData := append(nonce, ciphertext...)
	return encData, nil
}
