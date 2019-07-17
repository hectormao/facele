package cfg

type FacturaDianTrnsConfig interface {
	GetQRCodeUrl() string
	GetSignCert() SignCertConfig
}

type SignCertConfig interface {
	GetPath() string
	GetPassword() string
}
