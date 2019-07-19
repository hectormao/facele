package soap

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"time"

	soap "github.com/hooklift/gowsdl/soap"
)

type WSSSecurityHeader struct {
	XMLName   xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ wsse:Security"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	MustUnderstand string `xml:"mustUnderstand,attr,omitempty"`

	Token *WSSUsernameToken `xml:",omitempty"`
}

type WSSUsernameToken struct {
	XMLName   xml.Name `xml:"wsse:UsernameToken"`
	XmlNSWsu  string   `xml:"xmlns:wsu,attr"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	Id string `xml:"wsu:Id,attr,omitempty"`

	Username *soap.WSSUsername `xml:",omitempty"`
	Password *soap.WSSPassword `xml:",omitempty"`
	Created  *WSSCreated       `xml:",omitempty"`
	Nonce    *WSSNonce         `xml:",omitempty"`
}

type WSSCreated struct {
	XMLName  xml.Name `xml:"wsse:Created"`
	XmlNSWsu string   `xml:"xmlns:wsu,attr"`

	Data time.Time `xml:",chardata"`
}

type WSSNonce struct {
	XMLName   xml.Name `xml:"wsu:Nonce"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	Data string `xml:",chardata"`
}

// NewWSSSecurityHeader creates WSSSecurityHeader instance
func NewWSSSecurityHeader(user string, pass string, nonce string, created time.Time, mustUnderstand string) *WSSSecurityHeader {
	hdr := &WSSSecurityHeader{XmlNSWsse: soap.WssNsWSSE, MustUnderstand: mustUnderstand}
	hdr.Token = &WSSUsernameToken{XmlNSWsu: soap.WssNsWSU, XmlNSWsse: soap.WssNsWSSE}
	hdr.Token.Username = &soap.WSSUsername{XmlNSWsse: soap.WssNsWSSE, Data: user}

	crypt := sha256.New()
	crypt.Write([]byte(pass))
	encryptedPass := fmt.Sprintf("%x", crypt.Sum(nil))

	hdr.Token.Password = &soap.WSSPassword{XmlNSWsse: soap.WssNsWSSE, XmlNSType: soap.WssNsType, Data: encryptedPass}
	encodedNonce := base64.StdEncoding.EncodeToString([]byte(nonce))
	hdr.Token.Nonce = &WSSNonce{XmlNSWsse: soap.WssNsWSSE, Data: encodedNonce}
	hdr.Token.Created = &WSSCreated{XmlNSWsu: soap.WssNsWSU, Data: created}
	return hdr
}
