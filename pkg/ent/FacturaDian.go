package ent

import (
	"encoding/xml"
	"fmt"
	"time"
)

const (
	DianNSFe           string = "http://www.dian.gov.co/contratos/facturaelectronica/v1"
	DianNSCac          string = "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2"
	DianNSCbc          string = "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2"
	DianNSExt          string = "urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2"
	DianNSNs6          string = "http://www.w3.org/2000/09/xmldsig#"
	DianNSNs8          string = "http://uri.etsi.org/01903/v1.3.2#"
	DianNSSts          string = "http://www.dian.gov.co/contratos/facturaelectronica/v1/Structures"
	DianNSXsi          string = "http://www.w3.org/2001/XMLSchema-instance"
	DianSchemaLocation string = "http://www.dian.gov.co/contratos/facturaelectronica/v1 xsd/DIAN_UBL.xsd urn:un:unece:uncefact:data:specification:UnqualifiedDataTypesSchemaModule:2 xsd/UnqualifiedDataTypeSchemaModule-2.0.xsd urn:oasis:names:specification:ubl:schema:xsd:QualifiedDatatypes-2 xsd/UBL-QualifiedDatatypes-2.0.xsd"
	DianTextType       string = "cbc:TextType"
	CountrySchemeURI   string = "urn:oasis:names:specification:ubl:codelist:gc:CountryIdentificationCode-2.0"
	DianSchemeURI      string = "http://www.dian.gov.co/contratos/facturaelectronica/v1/InvoiceType"
	InvoiceDateFormat  string = "2006-01-02"
	InvoiceTimeFormat  string = "15:04:00"
)

func (invoice *InvoiceType) AgregarExtension(extension interface{}) {

	newExtension := UBLExtensionType{
		ExtensionContent: ExtensionContentType{
			Extension: extension,
		},
	}

	invoice.UBLExtensions.UBLExtensions =
		append(invoice.UBLExtensions.UBLExtensions, newExtension)
}

func NewInvoice() (*InvoiceType, error) {

	ublExtension := UBLExtensionType{
		ExtensionContent: ExtensionContentType{
			Extension: DianExtensionType{
				InvoiceControl: InvoiceControlType{
					InvoiceAuthorization: "18762010165197",
					AuthorizationPeriod: AuthorizationPeriodType{
						StartDate: invoiceDate{InvoiceDateFormat, time.Date(2018, 9, 10, 0, 0, 0, 0, time.UTC)},
						EndDate:   invoiceDate{InvoiceDateFormat, time.Date(2020, 3, 10, 0, 0, 0, 0, time.UTC)},
					},
					AuthorizedInvoices: AuthorizedInvoicesType{
						Prefix: PrefixType{
							Type: DianTextType,
							Data: "FV",
						},
						From: 1,
						To:   1000,
					},
				},
				InvoiceSource: InvoiceSourceType{
					IdentificationCode: newIdentificationCode(
						"6",
						"United Nations Economic Commission for Europe",
						CountrySchemeURI,
						"CO",
					),
				},
				SoftwareProvider: SoftwareProviderType{
					ProviderID: newProviderID(
						"195",
						"CO, DIAN (Direccion de Impuestos y Aduanas Nacionales)",
						"860028580",
					),
					SoftwareID: newSoftwareID(
						"195",
						"CO, DIAN (Direccion de Impuestos y Aduanas Nacionales)",
						"7095ba87-e8f1-4f66-be7b-ffedc0a6722f",
					),
				},
				SoftwareSecurityCode: newSoftwareSecurityCode(
					"195",
					"CO, DIAN (Direccion de Impuestos y Aduanas Nacionales)",
					"7e411f86923bed1e95b154686320d47c549f9c364452e276e2c579fa415e7e9346c9ef6a051fb8a28915476d0a5f3021",
				),
			},
		},
	}

	invoice := InvoiceType{
		XmlNSFe:           DianNSFe,
		XmlNSCac:          DianNSCac,
		XmlNSCbc:          DianNSCbc,
		XmlNSExt:          DianNSExt,
		XmlNSNs6:          DianNSNs6,
		XmlNSNs8:          DianNSNs8,
		XmlNSSts:          DianNSSts,
		XmlNSXsi:          DianNSXsi,
		XsiSchemaLocation: DianSchemaLocation,
		UBLExtensions: UBLExtensionsType{
			UBLExtensions: []UBLExtensionType{ublExtension},
		},
		UBLVersionID: "2.0",
		ProfileID:    "1.0",
		ID:           "FV81",
		UUID: newUUID(
			"195",
			"CO, DIAN (Direccion de Impuestos y Aduanas Nacionales)",
			"c56269b7dc3db40dc7d3a46d2715bfaf82200d7d",
		),
		IssueDate: invoiceDate{InvoiceDateFormat, time.Date(2018, 11, 23, 0, 0, 0, 0, time.UTC)},
		IssueTime: invoiceDate{InvoiceTimeFormat, time.Date(2018, 11, 23, 0, 0, 0, 0, time.UTC)},
		InvoiceTypeCode: newInvoiceTypeCode(
			"195",
			"CO, DIAN (Direccion de Impuestos y Aduanas Nacionales)",
			"1",
		),
		Note:                 "",
		DocumentCurrencyCode: "COP",
		AccountingSupplierParty: AccountingSupplierPartyType{
			AdditionalAccountID: "1",
			Party: PartyType{
				PartyIdentification: PartyIdentificationType{
					ID: newPartyID(
						"31",
						"195",
						"CO, DIAN (Direccion de Impuestos y Aduanas Nacionales)",
						"800194208",
					),
				},
				PartyName: PartyNameType{
					Name: "GESTION ENERGETICA S.A E.S.P",
				},
				PhysicalLocation: PhysicalLocationType{
					Address: AddressType{
						Department:          "CALDAS",
						CitySubdivisionName: "",
						CityName:            "MANIZALES",
						AddressLine: AddressLineType{
							Line: []LineType{
								LineType{
									Data: "CR 23 64 B 33 ED S XXI BRR LAURELES",
								},
							},
						},
						Country: CountryType{
							IdentificationCode: "CO",
						},
					},
				},
				PartyTaxScheme: PartyTaxSchemeType{
					TaxLevelCode: "2",
					TaxScheme:    "",
				},
				PartyLegalEntity: PartyLegalEntityType{
					RegistrationName: "GESTION ENERGETICA S.A E.S.P",
				},
			},
		},
		AccountingCustomerParty: AccountingCustomerPartyType{
			AdditionalAccountID: "1",
			Party: PartyType{
				PartyIdentification: PartyIdentificationType{
					ID: newPartyID(
						"31",
						"195",
						"CO, DIAN (Direccion de Impuestos y Aduanas Nacionales)",
						"811034077",
					),
				},
				PartyName: PartyNameType{
					Name: "EMPRESA DE ENERGIA ELECTRICA DE SERVICIOS PUBLICOS E.S.P DEL MUNICIPIO DE MURINDO",
				},
				PhysicalLocation: PhysicalLocationType{
					Address: AddressType{
						Department:          "ANTIOQUIA",
						CitySubdivisionName: "",
						CityName:            "MURINDO",
						AddressLine: AddressLineType{
							Line: []LineType{
								LineType{
									Data: "CALLE PRIMERA CABECERA MUNICIPAL",
								},
							},
						},
						Country: CountryType{
							IdentificationCode: "CO",
						},
					},
				},
				PartyTaxScheme: PartyTaxSchemeType{
					TaxLevelCode: "0",
					TaxScheme:    "",
				},
				PartyLegalEntity: PartyLegalEntityType{
					RegistrationName: "EMPRESA DE ENERGIA ELECTRICA DE SERVICIOS PUBLICOS E.S.P DEL MUNICIPIO DE MURINDO",
				},
			},
		},
		TaxTotal: TaxTotalType{
			TaxAmount: newTaxAmount(
				"COP",
				0.0,
			),
			TaxEvidenceIndicator: false,
			TaxSubtotal: TaxSubtotalType{
				TaxableAmount: newTaxableAmount(
					"COP",
					0.0,
				),
				TaxAmount: newTaxAmount(
					"COP",
					0.0,
				),
				Percent: 19.0,
				TaxCategory: TaxCategoryType{
					TaxScheme: TaxSchemeType{
						ID: "01",
					},
				},
			},
		},
		LegalMonetaryTotal: LegalMonetaryTotalType{
			LineExtensionAmount: newLineExtensionAmount(
				"COP",
				18135006.00,
			),
			TaxExclusiveAmount: newTaxExclusiveAmount(
				"COP",
				0.0,
			),
			PayableAmount: newPayableAmount(
				"COP",
				18135006.00,
			),
		},
		InvoiceLine: []InvoiceLineType{
			InvoiceLineType{
				ID:               "13",
				InvoicedQuantity: 1,
				LineExtensionAmount: newLineExtensionAmount(
					"COP",
					12564553.00,
				),
				Item: ItemType{
					Description: "REEMBOLSABLES POR CARGOS Y COSTOS ASOCIADOS A LA FRONTERA COMERCIAL RIOSUCIO - CAUCHERAS",
				},
				Price: PriceType{
					PriceAmount: newPriceAmountType(
						"COP",
						12564553.00,
					),
				},
			},
		},
	}

	return &invoice, nil
}

func newPriceAmountType(
	currencyID string,
	data float32) PriceAmountType {
	v := PriceAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func newPayableAmount(
	currencyID string,
	data float32) PayableAmountType {
	v := PayableAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func newTaxExclusiveAmount(
	currencyID string,
	data float32) TaxExclusiveAmountType {
	v := TaxExclusiveAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func newLineExtensionAmount(
	currencyID string,
	data float32) LineExtensionAmountType {
	v := LineExtensionAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v

}

func newTaxableAmount(
	currencyID string,
	data float32) TaxableAmountType {
	v := TaxableAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func newTaxAmount(
	currencyID string,
	data float32) TaxAmountType {
	v := TaxAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func newIdentificationCode(
	listAgencyID string,
	listAgencyName string,
	listSchemeURI string,
	data string) IdentificationCodeType {
	v := IdentificationCodeType{}
	v.ListAgencyID = listAgencyID
	v.ListAgencyName = listAgencyName
	v.ListSchemeURI = listSchemeURI
	v.Data = data
	return v
}

func newProviderID(
	schemeAgencyID string,
	schemeAgencyName string,
	data string) ProviderIDType {
	v := ProviderIDType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.Data = data
	return v
}

func newSoftwareID(
	schemeAgencyID string,
	schemeAgencyName string,
	data string) SoftwareIDType {
	v := SoftwareIDType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.Data = data
	return v
}

func newSoftwareSecurityCode(
	schemeAgencyID string,
	schemeAgencyName string,
	data string) SoftwareSecurityCodeType {
	v := SoftwareSecurityCodeType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.Data = data
	return v
}

func newUUID(
	schemeAgencyID string,
	schemeAgencyName string,
	data string) UUIDType {
	v := UUIDType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.Data = data
	return v
}
func newInvoiceTypeCode(
	listAgencyID string,
	listAgencyName string,
	data string) InvoiceTypeCodeType {
	v := InvoiceTypeCodeType{}
	v.ListAgencyID = listAgencyID
	v.ListAgencyName = listAgencyName
	v.ListSchemeURI = DianSchemeURI
	v.Data = data
	return v
}

func newPartyID(
	schemeID string,
	schemeAgencyID string,
	schemeAgencyName string,
	data string) PartyIDType {
	v := PartyIDType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.Data = data
	v.SchemeID = schemeID
	return v
}

type InvoiceType struct {
	XMLName                 xml.Name                    `xml:"fe:Invoice"`
	XmlNSFe                 string                      `xml:"xmlns:fe,attr"`
	XmlNSCac                string                      `xml:"xmlns:cac,attr"`
	XmlNSCbc                string                      `xml:"xmlns:cbc,attr"`
	XmlNSExt                string                      `xml:"xmlns:ext,attr"`
	XmlNSNs6                string                      `xml:"xmlns:ns6,attr"`
	XmlNSNs8                string                      `xml:"xmlns:ns8,attr"`
	XmlNSSts                string                      `xml:"xmlns:sts,attr"`
	XmlNSXsi                string                      `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation       string                      `xml:"xsi:schemaLocation,attr"`
	UBLExtensions           UBLExtensionsType           `xml:",omitempty"`
	UBLVersionID            string                      `xml:"cbc:UBLVersionID"`
	ProfileID               string                      `xml:"cbc:ProfileID"`
	ID                      string                      `xml:"cbc:ID"`
	UUID                    UUIDType                    `xml:",omitempty"`
	IssueDate               invoiceDate                 `xml:"cbc:IssueDate"`
	IssueTime               invoiceDate                 `xml:"cbc:IssueTime"`
	InvoiceTypeCode         InvoiceTypeCodeType         `xml:",omitempty"`
	Note                    string                      `xml:"cbc:Note"`
	DocumentCurrencyCode    string                      `xml:"cbc:DocumentCurrencyCode"`
	AccountingSupplierParty AccountingSupplierPartyType `xml:",omitempty"`
	AccountingCustomerParty AccountingCustomerPartyType `xml:",omitempty"`
	TaxTotal                TaxTotalType                `xml:",omitempty"`
	LegalMonetaryTotal      LegalMonetaryTotalType      `xml:",omitempty"`
	InvoiceLine             []InvoiceLineType           `xml:",omitempty"`
}

type UBLExtensionsType struct {
	XMLName       xml.Name           `xml:"ext:UBLExtensions"`
	UBLExtensions []UBLExtensionType `xml:",omitempty"`
}

type UBLExtensionType struct {
	XMLName          xml.Name             `xml:"ext:UBLExtension"`
	ExtensionContent ExtensionContentType `xml:",omitempty"`
}

type ExtensionContentType struct {
	XMLName   xml.Name    `xml:"ext:ExtensionContent"`
	Extension interface{} `xml:",omitempty"`
}

type DianExtensionType struct {
	XMLName              xml.Name                 `xml:"sts:DianExtensions"`
	InvoiceControl       InvoiceControlType       `xml:",omitempty"`
	InvoiceSource        InvoiceSourceType        `xml:",omitempty"`
	SoftwareProvider     SoftwareProviderType     `xml:",omitempty"`
	SoftwareSecurityCode SoftwareSecurityCodeType `xml:",omitempty"`
}

type InvoiceControlType struct {
	XMLName              xml.Name                `xml:"sts:InvoiceControl"`
	InvoiceAuthorization string                  `xml:"sts:InvoiceAuthorization"`
	AuthorizationPeriod  AuthorizationPeriodType `xml:",omiteempty"`
	AuthorizedInvoices   AuthorizedInvoicesType  `xml:",omiteempty"`
}

type AuthorizationPeriodType struct {
	XMLName   xml.Name    `xml:"sts:AuthorizationPeriod"`
	StartDate invoiceDate `xml:"cbc:StartDate"`
	EndDate   invoiceDate `xml:"cbc:EndDate"`
}

type AuthorizedInvoicesType struct {
	XMLName xml.Name   `xml:"sts:AuthorizedInvoices"`
	Prefix  PrefixType `xml:",omiteempty"`
	From    int        `xml:"sts:From"`
	To      int        `xml:"sts:To"`
}

type PrefixType struct {
	XMLName xml.Name `xml:"sts:Prefix"`
	Type    string   `xml:"xsi:type,attr"`
	Data    string   `xml:",chardata"`
}

type InvoiceSourceType struct {
	XMLName            xml.Name               `xml:"sts:InvoiceSource"`
	IdentificationCode IdentificationCodeType `xml:",omiteempty"`
}

type CodeType struct {
	ListAgencyID   string `xml:"listAgencyID,attr"`
	ListAgencyName string `xml:"listAgencyName,attr"`
	ListSchemeURI  string `xml:"listSchemeURI,attr"`
	Data           string `xml:",chardata"`
}

type IdentificationCodeType struct {
	XMLName xml.Name `xml:"cbc:IdentificationCode"`
	CodeType
}

type InvoiceTypeCodeType struct {
	XMLName xml.Name `xml:"cbc:InvoiceTypeCode"`
	CodeType
}

type SoftwareProviderType struct {
	XMLName    xml.Name       `xml:"sts:SoftwareProvider"`
	ProviderID ProviderIDType `xml:",omiteempty"`
	SoftwareID SoftwareIDType `xml:",omiteempty"`
}

type IDType struct {
	SchemeAgencyID   string `xml:"schemeAgencyID,attr"`
	SchemeAgencyName string `xml:"schemeAgencyName,attr"`
	Data             string `xml:",chardata"`
}

type ProviderIDType struct {
	XMLName xml.Name `xml:"sts:ProviderID"`
	IDType
}

type SoftwareIDType struct {
	XMLName xml.Name `xml:"sts:SoftwareID"`
	IDType
}

type SoftwareSecurityCodeType struct {
	XMLName xml.Name `xml:"sts:SoftwareSecurityCode"`
	IDType
}

type UUIDType struct {
	XMLName xml.Name `xml:"cbc:UUID"`
	IDType
}

type PartyIDType struct {
	XMLName  xml.Name `xml:"cbc:ID"`
	SchemeID string   `xml:"schemeID,attr"`
	IDType
}

type AccountingSupplierPartyType struct {
	XMLName             xml.Name  `xml:"fe:AccountingSupplierParty"`
	AdditionalAccountID string    `xml:"cbc:AdditionalAccountID"`
	Party               PartyType `xml:",omiteempty"`
}

type AccountingCustomerPartyType struct {
	XMLName             xml.Name  `xml:"fe:AccountingCustomerParty"`
	AdditionalAccountID string    `xml:"cbc:AdditionalAccountID"`
	Party               PartyType `xml:",omiteempty"`
}

type PartyType struct {
	XMLName             xml.Name                `xml:"fe:Party"`
	PartyIdentification PartyIdentificationType `xml:",omitempty"`
	PartyName           PartyNameType           `xml:",omitempty"`
	PhysicalLocation    PhysicalLocationType    `xml:",omitempty"`
	PartyTaxScheme      PartyTaxSchemeType      `xml:",omitempty"`
	PartyLegalEntity    PartyLegalEntityType    `xml:",omitempty"`
}

type PartyIdentificationType struct {
	XMLName xml.Name    `xml:"cac:PartyIdentification"`
	ID      PartyIDType `xml:",omiteempty"`
}

type PartyNameType struct {
	XMLName xml.Name `xml:"cac:PartyName"`
	Name    string   `xml:"cbc:Name"`
}

type PhysicalLocationType struct {
	XMLName xml.Name    `xml:"fe:PhysicalLocation"`
	Address AddressType `xml:",omiteempty"`
}

type AddressType struct {
	XMLName             xml.Name        `xml:"fe:Address"`
	Department          string          `xml:"cbc:Department"`
	CitySubdivisionName string          `xml:"cbc:CitySubdivisionName"`
	CityName            string          `xml:"cbc:CityName"`
	AddressLine         AddressLineType `xml:",omiteempty"`
	Country             CountryType     `xml:",omiteempty"`
}

type AddressLineType struct {
	XMLName xml.Name   `xml:"cac:AddressLine"`
	Line    []LineType `xml:",omiteempty"`
}

type LineType struct {
	XMLName xml.Name `xml:"cbc:Line"`
	Data    string   `xml:",chardata"`
}

type CountryType struct {
	XMLName            xml.Name `xml:"cac:Country"`
	IdentificationCode string   `xml:"cbc:IdentificationCode"`
}

type PartyTaxSchemeType struct {
	XMLName      xml.Name `xml:"fe:PartyTaxScheme"`
	TaxLevelCode string   `xml:"cbc:TaxLevelCode"`
	TaxScheme    string   `xml:"cac:TaxScheme"`
}

type PartyLegalEntityType struct {
	XMLName          xml.Name `xml:"fe:PartyLegalEntity"`
	RegistrationName string   `xml:"cbc:RegistrationName"`
}

type TaxTotalType struct {
	XMLName              xml.Name        `xml:"fe:TaxTotal"`
	TaxAmount            TaxAmountType   `xml:",omitempty"`
	TaxEvidenceIndicator bool            `xml:"cbc:TaxEvidenceIndicator"`
	TaxSubtotal          TaxSubtotalType `xml:",omitempty"`
}

type AmountType struct {
	CurrencyID string  `xml:"currencyID,attr"`
	Data       float32 `xml:",chardata"`
}

func (a AmountType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	strValue := fmt.Sprintf("%.2f", a.Data)
	attr := xml.Attr{xml.Name{"", "currencyID"}, a.CurrencyID}
	start.Attr = []xml.Attr{attr}
	err := e.EncodeElement(strValue, start)
	if err != nil {
		return err
	}
	return nil
}

type TaxAmountType struct {
	XMLName xml.Name `xml:"cbc:TaxAmount"`
	AmountType
}

type TaxSubtotalType struct {
	XMLName       xml.Name          `xml:"fe:TaxSubtotal"`
	TaxableAmount TaxableAmountType `xml:",omitempty"`
	TaxAmount     TaxAmountType     `xml:",omitempty"`
	Percent       float64           `xml:"cbc:Percent"`
	TaxCategory   TaxCategoryType   `xml:",omitempty"`
}

type TaxableAmountType struct {
	XMLName xml.Name `xml:"cbc:TaxableAmount"`
	AmountType
}

type TaxCategoryType struct {
	XMLName   xml.Name      `xml:"cac:TaxCategory"`
	TaxScheme TaxSchemeType `xml:",omitempty"`
}

type TaxSchemeType struct {
	XMLName xml.Name `xml:"cac:TaxScheme"`
	ID      string   `xml:"cbc:ID"`
}

type LegalMonetaryTotalType struct {
	XMLName             xml.Name                `xml:"fe:LegalMonetaryTotal"`
	LineExtensionAmount LineExtensionAmountType `xml:",omitempty"`
	TaxExclusiveAmount  TaxExclusiveAmountType  `xml:",omitempty"`
	PayableAmount       PayableAmountType       `xml:",omitempty"`
}

type LineExtensionAmountType struct {
	XMLName xml.Name `xml:"cbc:LineExtensionAmount"`
	AmountType
}

type TaxExclusiveAmountType struct {
	XMLName xml.Name `xml:"cbc:TaxExclusiveAmount"`
	AmountType
}

type PayableAmountType struct {
	XMLName xml.Name `xml:"cbc:PayableAmount"`
	AmountType
}

type InvoiceLineType struct {
	XMLName             xml.Name                `xml:"fe:InvoiceLine"`
	ID                  string                  `xml:"cbc:ID"`
	InvoicedQuantity    int                     `xml:"cbc:InvoicedQuantity"`
	LineExtensionAmount LineExtensionAmountType `xml:",omitempty"`
	Item                ItemType                `xml:",omitempty"`
	Price               PriceType               `xml:",omitempty"`
}

type ItemType struct {
	XMLName     xml.Name `xml:"fe:Item"`
	Description string   `xml:"cbc:Description"`
}

type PriceType struct {
	XMLName     xml.Name        `xml:"fe:Price"`
	PriceAmount PriceAmountType `xml:",omitempty"`
}

type PriceAmountType struct {
	XMLName xml.Name `xml:"cbc:PriceAmount"`
	AmountType
}

type invoiceDate struct {
	FormatString string
	time.Time
}

func (t invoiceDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	fmted := t.Format(t.FormatString)
	err := e.EncodeElement(fmted, start)
	if err != nil {
		return err
	}
	return nil
}
func (t *invoiceDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(t.FormatString, v)
	if err != nil {
		return err
	}
	*t = invoiceDate{t.FormatString, parse}
	return nil
}
