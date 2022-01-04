package hyper

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
)

var getenv = os.Getenv
var _ TLSDescriptor = environmentTLS{}

type TLSDescriptor interface {
	from() (*tls.Config, error)
}

// -----------------------------------
// TLS from environment

type environmentTLS struct {
	env EnvironmentTLSOpts
}

type EnvironmentTLSOpts struct {
	TLSCert          string
	TLSCertBlockType string
	TLSKey           string
	TLSKeyBlockType  string
}

func NewEnvironmentTLS(e EnvironmentTLSOpts) TLSDescriptor {
	return environmentTLS{env: e}
}

func (t environmentTLS) from() (*tls.Config, error) {
	cert := getenv(t.env.TLSCert)
	certBlockType := getenv(t.env.TLSCertBlockType)
	key := getenv(t.env.TLSKey)
	keyBlockName := getenv(t.env.TLSKeyBlockType)

	if cert == "" || certBlockType == "" || key == "" || keyBlockName == "" {
		return nil, errors.New("invalid certification setup")
	}

	certDecoded, err := base64.StdEncoding.DecodeString(cert)
	if err != nil {
		return nil, err
	}
	certBlock := pem.Block{
		Type:  certBlockType,
		Bytes: certDecoded,
	}
	certBuf := new(bytes.Buffer)
	if err := pem.Encode(certBuf, &certBlock); err != nil {
		return nil, err
	}

	keyDecoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	keyBlock := pem.Block{
		Type:  keyBlockName,
		Bytes: keyDecoded,
	}
	keyBuf := new(bytes.Buffer)
	if err := pem.Encode(keyBuf, &keyBlock); err != nil {
		return nil, err
	}

	certificate, err := tls.X509KeyPair(certBuf.Bytes(), keyBuf.Bytes())
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{certificate}}, nil
}

// -----------------------------------
// TLS from file

type fileTLS struct {
	env FileTLSOpts
}

type FileTLSOpts struct {
	TLSCertPath string
	TLSKeyPath  string
}

func NewFileTLS(e FileTLSOpts) TLSDescriptor {
	return fileTLS{env: e}
}

func (t fileTLS) from() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(os.Getenv(t.env.TLSCertPath), os.Getenv(t.env.TLSKeyPath))
	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}
