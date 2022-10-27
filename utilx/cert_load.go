package utilx

import (
    "crypto"
    "crypto/ecdsa"
    "crypto/rsa"
    "crypto/tls"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "github.com/potakhov/loge"
    "io/ioutil"
    "path/filepath"
    "sort"
)

// LoadCertficateAndKeyFromFile reads file, divides into key and certificates
func LoadCertficateAndKeyFromFile(path string) (*tls.Certificate, error) {
    raw, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cert tls.Certificate
    for {
        block, rest := pem.Decode(raw)
        if block == nil {
            break
        }
        if block.Type == "CERTIFICATE" {
            cert.Certificate = append(cert.Certificate, block.Bytes)
        } else {
            cert.PrivateKey, err = parsePrivateKey(block.Bytes)
            if err != nil {
                return nil, fmt.Errorf("failure reading private key from \"%s\": %s", path, err)
            }
        }
        raw = rest
    }

    if len(cert.Certificate) == 0 {
        return nil, fmt.Errorf("no certificate found in \"%s\"", path)
    } else if cert.PrivateKey == nil {
        return nil, fmt.Errorf("no private key found in \"%s\"", path)
    }

    return &cert, nil
}

// LoadCertificateDirectory globs all .pem files in given directory, parses them
// for certs (and private keys) and returns them
func LoadCertificateDirectory(dir string) ([]tls.Certificate, error) {
    // read certificate files
    certficateFiles, err := filepath.Glob(filepath.Join(dir, "*.pem"))
    if err != nil {
        return nil, fmt.Errorf("Failed to scan certificate dir \"%s\": %s", dir, err)
    }
    sort.Strings(certficateFiles)
    certs := make([]tls.Certificate, 0)
    for _, file := range certficateFiles {
        cert, err := LoadCertficateAndKeyFromFile(file)
        if err != nil {
            loge.Info("common(tls): %s", err)
        } else {
            certs = append(certs, *cert)
        }
    }
    return certs, nil
}

func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
    if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
        return key, nil
    }
    if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
        switch key := key.(type) {
        case *rsa.PrivateKey, *ecdsa.PrivateKey:
            return key, nil
        default:
            return nil, fmt.Errorf("Found unknown private key type in PKCS#8 wrapping")
        }
    }
    if key, err := x509.ParseECPrivateKey(der); err == nil {
        return key, nil
    }
    return nil, fmt.Errorf("Failed to parse private key")
}
