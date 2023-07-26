package cert

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	TYPE_PUBLIC  = "PUBLIC KEY"
	TYPE_PRIVATE = "PRIVATE KEY"
)

var ErrInvalidKeyPair = errors.New("invalid key pair")

func MarshalED25519PKCS8(priv ed25519.PrivateKey, password []byte) (*pem.Block, *pem.Block, error) {
	keyPUB, err := x509.MarshalPKIXPublicKey(priv.Public())
	if err != nil {
		return nil, nil, err
	}

	pemPUB := &pem.Block{Type: TYPE_PUBLIC, Bytes: keyPUB}

	keyPRIV, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}

	pemPRIV := &pem.Block{Type: TYPE_PRIVATE, Bytes: keyPRIV}

	// although EncryptPEMBlock is deprecated because of padding oracle
	// it's usable here for private key encryption
	//nolint:staticcheck
	pemPRIVEncrypted, err := x509.EncryptPEMBlock(rand.Reader, pemPRIV.Type, pemPRIV.Bytes, password, x509.PEMCipherAES256)
	if err != nil {
		return nil, nil, err
	}

	return pemPUB, pemPRIVEncrypted, nil
}

func UnmarshalED25519PKCS8(pemPRIVEncrypted *pem.Block, password []byte) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	//nolint:staticcheck
	keyPRIVBytes, err := x509.DecryptPEMBlock(pemPRIVEncrypted, password)
	if err != nil {
		return nil, nil, ErrInvalidKeyPair
	}

	keyPRIV, err := x509.ParsePKCS8PrivateKey(keyPRIVBytes)
	if err != nil {
		return nil, nil, ErrInvalidKeyPair
	}

	priv, ok := keyPRIV.(ed25519.PrivateKey)
	if !ok {
		return nil, nil, ErrInvalidKeyPair
	}

	pub, ok := priv.Public().(ed25519.PublicKey)
	if !ok {
		return nil, nil, ErrInvalidKeyPair
	}

	return pub, priv, nil
}

func GenerateED25519Pair(path string, password []byte) (*pem.Block, *pem.Block, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	pemPUB, pemPRIV, err := MarshalED25519PKCS8(priv, password)
	if err != nil {
		return nil, nil, err
	}

	certWriter, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open %s for writing: %w", path, err)
	}

	if err := pem.Encode(certWriter, pemPUB); err != nil {
		return nil, nil, fmt.Errorf("failed to write data to %s: %w", path, err)
	}

	if err := pem.Encode(certWriter, pemPRIV); err != nil {
		return nil, nil, fmt.Errorf("failed to write data to %s: %w", path, err)
	}

	if err := certWriter.Close(); err != nil {
		return nil, nil, fmt.Errorf("failed to close %s: %w", path, err)
	}

	return pemPUB, pemPRIV, nil
}

type ED25519JWK interface {
	GetCrv() string
	GetKty() string
	GetD() string
}

func ED25519FromJWK(jwk ED25519JWK) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	if jwk.GetCrv() != "Ed25519" || jwk.GetKty() != "OKP" {
		return nil, nil, ErrInvalidKeyPair
	}

	d := jwk.GetD()
	if d == "" {
		return nil, nil, ErrInvalidKeyPair
	}

	b, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(d, "="))
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %q", ErrInvalidKeyPair, err)
	}

	priv := ed25519.NewKeyFromSeed(b)
	pub, ok := priv.Public().(ed25519.PublicKey)
	if !ok {
		return nil, nil, ErrInvalidKeyPair
	}

	return pub, priv, nil
}

func LoadED25519Pair(b []byte, password []byte) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	var pemPUB, pemPRIV *pem.Block

	for {
		block, rest := pem.Decode(b)
		if block == nil {
			break
		}

		b = rest

		if block.Type == TYPE_PUBLIC {
			pemPUB = block
			continue
		}

		if block.Type == TYPE_PRIVATE {
			pemPRIV = block
			continue
		}
	}

	if pemPUB == nil || pemPRIV == nil {
		return nil, nil, ErrInvalidKeyPair
	}

	return UnmarshalED25519PKCS8(pemPRIV, password)
}
