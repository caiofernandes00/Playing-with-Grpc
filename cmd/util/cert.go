package util

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

func LoadCAPool() (*x509.CertPool, error) {
	pemCA, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, fmt.Errorf("cannot load TLS certificates: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemCA) {
		return nil, fmt.Errorf("cannot add CA certificate")
	}

	return certPool, nil

}
