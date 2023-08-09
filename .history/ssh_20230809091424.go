package gox

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// SshKey
type SshKey struct {
	PubKey  ssh.PublicKey
	Key     []byte
	Comment string
	KeyType string
	Options []string
	Rest    []byte
}

func (e *SshKey) String() string {
	return fmt.Sprintf("SshKey KeyType: %v comment: %v", e.KeyType, e.Comment)
}

// LoadAuthorizedKeys load keys from authorized_keys file and puts the marshalled public key into the map as the key.
func LoadAuthorizedKeys(authorizedKeyFile string) ([]*SshKey, error) {
	authorizedKeysBytes, err := ioutil.ReadFile(authorizedKeyFile)
	if err != nil {
		loge.Printf("Failed to load authorized_keys, err: %v", err)
		fmt.Printf("Failed to load authorized_keys, err: %v", err)
		return nil, err
	}

	authorizedKeys := make([]*SshKey, 0)
	for len(authorizedKeysBytes) > 0 {
		//fmt.Printf("authorizedKeysBytes: %d\n", len(authorizedKeysBytes))
		pubKey, comment, options, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)

		//fmt.Printf("pubKey: %v\n", pubKey)
		//fmt.Printf("comment: %v\n", comment)
		//fmt.Printf("options: %v\n", options)
		//fmt.Printf("rest: %v\n", string(rest))
		//fmt.Printf("err: %v\n", err)

		if err != nil && len(rest) == 0 {
			//fmt.Printf("tes1: %v\n", err)
			break
			//return nil, err
		}

		key := &SshKey{
			PubKey:  pubKey,
			KeyType: pubKey.Type(),
			Key:     pubKey.Marshal(),
			Comment: comment,
			Options: options,
			Rest:    rest,
		}

		authorizedKeys = append(authorizedKeys, key)
		authorizedKeysBytes = rest
		loge.Printf("authorizedKeysMap add key type: %v - %v\n", key.KeyType, key.Comment)
	}

	return authorizedKeys, err
}

// GeneratePrivateKey creates a RSA Private Key of specified byte size
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	loge.Println("Private Key generated")
	return privateKey, nil
}

// EncodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// GeneratePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func GeneratePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	loge.Println("Public key generated")
	return pubKeyBytes, nil
}

// WriteKeyToFile writes keys to a file
func WriteKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	loge.Printf("Key saved to: %s", saveFileTo)
	return nil
}
