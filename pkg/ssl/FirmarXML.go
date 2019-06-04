package ssl

import (
	"crypto/tls"
	"encoding/pem"
	"io/ioutil"

	"github.com/amdonov/xmlsig"
	"golang.org/x/crypto/pkcs12"
)

func FirmarXML(obj interface{}) (*xmlsig.Signature, error) {

	p12Data, err := ioutil.ReadFile("../../resources/ssl/certificado.p12")

	if err != nil {
		return nil, err
	}

	p12Block, err := pkcs12.ToPEM(p12Data, "Auditoria000")

	if err != nil {
		return nil, err
	}

	var pemData []byte
	for _, block := range p12Block {
		pemData = append(pemData, pem.EncodeToMemory(block)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)

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
