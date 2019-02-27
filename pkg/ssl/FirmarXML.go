package ssl

import (
	"crypto/tls"

	"github.com/amdonov/xmlsig"
)

func FirmarXML(obj interface{}) (*xmlsig.Signature, error) {
	cert, err := tls.LoadX509KeyPair(
		"../../resources/ssl/certificate.pem",
		"../../resources/ssl/key.pem",
	)
	if err != nil {
		return nil, err
	}

	signer, err := xmlsig.NewSigner(cert)
	if err != nil {
		return nil, err
	}

	sign, err := signer.CreateSignature(obj)
	if err != nil {
		return nil, err
	}

	return sign, nil
}
