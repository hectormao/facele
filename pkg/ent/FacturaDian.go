package ent

import (
	"encoding/xml"
	"fmt"
	"time"
)

const (
	DianNS             string = "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2"
	DianNSCac          string = "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2"
	DianNSCbc          string = "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2"
	DianNSExt          string = "urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2"
	DianNSNs6          string = "http://www.w3.org/2000/09/xmldsig#"
	DianNSNs8          string = "http://uri.etsi.org/01903/v1.3.2#"
	DianNSSts          string = "dian:gov:co:facturaelectronica:Structures-2-1"
	DianNSXades        string = "http://uri.etsi.org/01903/v1.3.2#"
	DianNSXades141     string = "http://uri.etsi.org/01903/v1.4.1#"
	DianNSXsi          string = "http://www.w3.org/2001/XMLSchema-instance"
	DianSchemaLocation string = "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2     http://docs.oasis-open.org/ubl/os-UBL-2.1/xsd/maindoc/UBL-Invoice-2.1.xsd"
	DianTextType       string = "cbc:TextType"
	CountrySchemeURI   string = "urn:oasis:names:specification:ubl:codelist:gc:CountryIdentificationCode-2.1"
	DianSchemeURI      string = "http://www.dian.gov.co/contratos/facturaelectronica/v1/InvoiceType"
	InvoiceDateFormat  string = "2006-01-02"
	InvoiceTimeFormat  string = "15:04:00"
	AgencyID           string = "195"
	AgencyName         string = "CO, DIAN (Direccion de Impuestos y Aduanas Nacionales)"
	SchemeID           string = "9"
	SchemeID4          string = "9"
	SchemeName         string = "31"
	UBLVersion         string = "UBL 2.1"
	Customization      string = "1"
	Profile            string = "Dian 2.0"
	ProfileExecution   string = "2"
	ListAgencyID       string = "6"
	ListAgencyName     string = "United Nations Economic Commission for Europe"
	ListName           string = "05"
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

type InvoiceType struct {
	XMLName                     xml.Name                        `xml:"Invoice"`
	XmlNS                       string                          `xml:"xmlns,attr"`
	XmlNSCac                    string                          `xml:"xmlns:cac,attr"`
	XmlNSCbc                    string                          `xml:"xmlns:cbc,attr"`
	XmlNSExt                    string                          `xml:"xmlns:ext,attr"`
	XmlNSSts                    string                          `xml:"xmlns:sts,attr"`
	XmlNSXades                  string                          `xml:"xmlns:xades,attr"`
	XmlNSXades141               string                          `xml:"xmlns:xades141,attr"`
	XmlNSXsi                    string                          `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation           string                          `xml:"xsi:schemaLocation,attr"`
	UBLExtensions               UBLExtensionsType               `xml:",omitempty"`
	UBLVersionID                string                          `xml:"cbc:UBLVersionID"`
	CustomizationID             string                          `xml:"cbc:CustomizationID"`
	ProfileID                   string                          `xml:"cbc:ProfileID"`
	ProfileExecutionID          string                          `xml:cbc:ProfileExecutionID`
	ID                          string                          `xml:"cbc:ID"`
	UUID                        UUIDType                        `xml:",omitempty"`
	IssueDate                   InvoiceDate                     `xml:"cbc:IssueDate"`
	IssueTime                   InvoiceDate                     `xml:"cbc:IssueTime"`
	DueDate                     InvoiceDate                     `xml:"cbc:DueDate"`
	InvoiceTypeCode             string                          `xml:"cbc:InvoiceTypeCode"`
	Note                        string                          `xml:"cbc:Note"`
	TaxPointDate                InvoiceDate                     `xml:"cbc:TaxPointDate"`
	DocumentCurrencyCode        string                          `xml:"cbc:DocumentCurrencyCode"`
	LineCountNumeric            int                             `xml:"cbc:LineCountNumeric"`
	InvoicePeriod               InvoicePeriodType               `xml:",omitempty"`
	OrderReference              ReferenceType                   `xml:"cac:OrderReference,omitempty"`
	DespatchDocumentReference   ReferenceType                   `xml:"cac:DespatchDocumentReference,omitempty"`
	ReceiptDocumentReference    ReferenceType                   `xml:"cac:ReceiptDocumentReference,omitempty"`
	AdditionalDocumentReference AdditionalDocumentReferenceType `xml:"cac:AdditionalDocumentReference,omitempty"`
	AccountingSupplierParty     AccountingSupplierPartyType     `xml:",omitempty"`
	AccountingCustomerParty     AccountingCustomerPartyType     `xml:",omitempty"`
	TaxRepresentativeParty      TaxRepresentativePartyType      `xml:"cac:TaxRepresentativeParty"`
	Delivery                    DeliveryType                    `xml:"cac:Delivery"`
	PaymentMeans                PaymentMeansType                `xml:"cac:PaymentMeans"`
	PaymentTerms                PaymentTermsType                `xml:"cac:PaymentTerms"`
	AllowanceCharge             []AllowanceChargeType           `xml:"cac:AllowanceCharge"`
	PaymentExchangeRate         PaymentExchangeRateType         `xml:"cac:PaymentExchangeRate"`
	TaxTotal                    []TaxTotalType                  `xml:",omitempty"`
	LegalMonetaryTotal          LegalMonetaryTotalType          `xml:",omitempty"`
	InvoiceLine                 []InvoiceLineType               `xml:",omitempty"`
}

type AdditionalDocumentReferenceType struct {
	ID               string `xml:"cbc:ID"`
	DocumentTypeCode string `xml:"cbc:DocumentTypeCode"`
}

type PaymentExchangeRateType struct {
	SourceCurrencyCode     string      `xml:"cbc:SourceCurrencyCode"`
	SourceCurrencyBaseRate float64     `xml:"cbc:SourceCurrencyBaseRate"`
	TargetCurrencyCode     string      `xml:"cbc:TargetCurrencyCode"`
	TargetCurrencyBaseRate float64     `xml:"cbc:TargetCurrencyBaseRate"`
	CalculationRate        float64     `xml:"cbc:CalculationRate"`
	Date                   InvoiceDate `xml:"cbc:Date"`
}

type AllowanceChargeType struct {
	ID                        string     `xml:"cbc:ID"`
	ChargeIndicator           bool       `xml:"cbc:ChargeIndicator"`
	AllowanceChargeReasonCode int        `xml:"cbc:AllowanceChargeReasonCode,omitempty"`
	AllowanceChargeReason     string     `xml:"cbc:AllowanceChargeReason"`
	MultiplierFactorNumeric   float64    `xml:"cbc:MultiplierFactorNumeric"`
	Amount                    AmountType `xml:"cbc:Amount"`
	BaseAmount                AmountType `xml:"cbc:BaseAmount"`
}

type PaymentTermsType struct {
	ReferenceEventCode int                  `xml:"cbc:ReferenceEventCode"`
	SettlementPeriod   SettlementPeriodType `xml:"cac:SettlementPeriod"`
}

type SettlementPeriodType struct {
	DurationMeasure DurationMeasureType `xml:",omitempty"`
}

type PaymentMeansType struct {
	ID               string      `xml:"cbc:ID"`
	PaymentMeansCode string      `xml:"cbc:PaymentMeansCode"`
	PaymentDueDate   InvoiceDate `xml:"cbc:PaymentDueDate"`
}

type DeliveryType struct {
	ActualDeliveryDate InvoiceDate `xml:"cbc:ActualDeliveryDate"`
	ActualDeliveryTime InvoiceDate `xml:"cbc:ActualDeliveryTime"`
	DeliveryAddress    AddressType `xml:"cac:DeliveryAddress"`
	DeliveryParty      PartyType   `xml:"cac:DeliveryParty"`
}

type TaxRepresentativePartyType struct {
	PartyIdentification IDType   `xml:"cac:PartyIdentification"`
	PartyName           NameType `xml:"cac:PartyName"`
}

type ReferenceType struct {
	ID               string `xml:"cbc:ID"`
	DocumentTypeCode string `xml:"cbc:DocumentTypeCode,omitempty"`
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
	XMLName               xml.Name                  `xml:"sts:DianExtensions"`
	InvoiceControl        InvoiceControlType        `xml:",omitempty"`
	InvoiceSource         InvoiceSourceType         `xml:",omitempty"`
	SoftwareProvider      SoftwareProviderType      `xml:",omitempty"`
	SoftwareSecurityCode  SoftwareSecurityCodeType  `xml:",omitempty"`
	AuthorizationProvider AuthorizationProviderType `xml:",omitempty"`
	QRCode                string                    `xml:"sts:QRCode"`
}

type AuthorizationProviderType struct {
	XMLName                 xml.Name                    `xml:"sts:AuthorizationProvider"`
	AuthorizationProviderID AuthorizationProviderIDType `xml:",omitempty`
}

type AuthorizationProviderIDType struct {
	XMLName xml.Name `xml:"sts:AuthorizationProviderID"`
	IDType
}

type InvoiceControlType struct {
	XMLName              xml.Name                `xml:"sts:InvoiceControl"`
	InvoiceAuthorization string                  `xml:"sts:InvoiceAuthorization"`
	AuthorizationPeriod  AuthorizationPeriodType `xml:",omiteempty"`
	AuthorizedInvoices   AuthorizedInvoicesType  `xml:",omiteempty"`
}

type AuthorizationPeriodType struct {
	XMLName   xml.Name    `xml:"sts:AuthorizationPeriod"`
	StartDate InvoiceDate `xml:"cbc:StartDate"`
	EndDate   InvoiceDate `xml:"cbc:EndDate"`
}

type AuthorizedInvoicesType struct {
	XMLName xml.Name `xml:"sts:AuthorizedInvoices"`
	Prefix  string   `xml:sts:Prefix",omiteempty"`
	From    int64    `xml:"sts:From"`
	To      int64    `xml:"sts:To"`
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

type InvoicePeriodType struct {
	XMLName   xml.Name    `xml:"cac:InvoicePeriod"`
	StartDate InvoiceDate `xml:"cbc:StartDate"`
	StartTime InvoiceDate `xml:"cbc:StartTime"`
	EndDate   InvoiceDate `xml:"cbc:EndDate"`
	EndTime   InvoiceDate `xml:"cbc:EndTime"`
}

type DurationMeasureType struct {
	XMLName  xml.Name `xml:"cbc:DurationMeasure"`
	UnitCode string   `xml:"unitCode,attr"`
	Data     int64    `xml:",chardata"`
}

type SoftwareProviderType struct {
	XMLName    xml.Name       `xml:"sts:SoftwareProvider"`
	ProviderID ProviderIDType `xml:",omiteempty"`
	SoftwareID SoftwareIDType `xml:",omiteempty"`
}

type IDType struct {
	SchemeAgencyID   string `xml:"schemeAgencyID,attr,omiteempty"`
	SchemeAgencyName string `xml:"schemeAgencyName,attr,omiteempty"`
	SchemeID         string `xml:"schemeID,attr,omiteempty"`
	SchemeName       string `xml:"schemeName,attr,omiteempty"`
	Data             string `xml:",chardata"`
}

type ProviderIDType struct {
	XMLName xml.Name `xml:"sts:ProviderID"`
	IDType
}

type SoftwareIDType struct {
	XMLName xml.Name `xml:"sts:SoftwareIDs"`
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
	XMLName             xml.Name  `xml:"cac:AccountingSupplierParty"`
	AdditionalAccountID IDType    `xml:"cbc:AdditionalAccountID"`
	Party               PartyType `xml:"cac:Party"`
}

type AccountingCustomerPartyType struct {
	XMLName             xml.Name  `xml:"cac:AccountingCustomerParty"`
	AdditionalAccountID IDType    `xml:"cbc:AdditionalAccountID"`
	Party               PartyType `xml:"cac:Party"`
}

type PartyType struct {
	PartyName        []PartyNameType      `xml:",omitempty"`
	PhysicalLocation PhysicalLocationType `xml:",omitempty"`
	PartyTaxScheme   PartyTaxSchemeType   `xml:",omitempty"`
	PartyLegalEntity PartyLegalEntityType `xml:",omitempty"`
	Contact          ContactType          `xml:"cac:Contact,omitempty"`
}

type ContactType struct {
	Telephone      string `xml:"cbc:Telephone"`
	ElectronicMail string `xml:"cbc:ElectronicMail"`
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
	XMLName xml.Name    `xml:"cac:PhysicalLocation"`
	Address AddressType `xml:"cac:Address"`
}

type AddressType struct {
	XMLName              xml.Name        `xml:"cac:Address"`
	ID                   string          `xml:"cbc:ID"`
	CityName             string          `xml:"cbc:CityName"`
	CountrySubentity     string          `xml:"cbc:CountrySubentity"`
	CountrySubentityCode string          `xml:"cbc:CountrySubentityCode"`
	AddressLine          AddressLineType `xml:",omiteempty"`
	Country              CountryType     `xml:",omiteempty"`
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
	Name               string   `xml:"cbc:Name"`
}

type PartyTaxSchemeType struct {
	XMLName             xml.Name         `xml:"cac:PartyTaxScheme"`
	RegistrationName    string           `xml:"cbc:RegistrationName"`
	CompanyID           IDType           `xml:"cbc:CompanyID"`
	TaxLevelCode        TaxLevelCodeType `xml:",omitempty"`
	RegistrationAddress AddressType      `xml:"cac:RegistrationAddress"`
	TaxScheme           NameType         `xml:"cac:TaxScheme"`
}

type NameType struct {
	ID   string `xml:"cbc:ID,omitempty"`
	Name string `xml:"cbc:Name"`
}

type TaxLevelCodeType struct {
	XMLName  xml.Name `xml:"cbc:TaxLevelCode"`
	ListName string   `xml:"listName,attr"`
	Data     string   `xml:",chardata"`
}

type PartyLegalEntityType struct {
	XMLName                     xml.Name `xml:"cac:PartyLegalEntity"`
	RegistrationName            string   `xml:"cbc:RegistrationName"`
	CompanyID                   IDType   `xml:"cbc:CompanyID"`
	CorporateRegistrationScheme NameType `xml:"cac:CorporateRegistrationScheme"`
}

type TaxTotalType struct {
	TaxAmount            AmountType      `xml:"cbc:TaxAmount"`
	TaxEvidenceIndicator bool            `xml:"cbc:TaxEvidenceIndicator"`
	TaxSubtotal          TaxSubtotalType `xml:"cac:TaxSubtotal"`
}

type AmountType struct {
	CurrencyID string  `xml:"currencyID,attr"`
	Data       float64 `xml:",chardata"`
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
	TaxableAmount AmountType      `xml:"cbc:TaxableAmount"`
	TaxAmount     AmountType      `xml:"cbc:TaxAmount"`
	TaxCategory   TaxCategoryType `xml:",omitempty"`
}

type TaxableAmountType struct {
	XMLName xml.Name `xml:"cbc:TaxableAmount"`
	AmountType
}

type TaxCategoryType struct {
	XMLName   xml.Name `xml:"cac:TaxCategory"`
	Percent   float64  `xml:"cbc:Percent"`
	TaxScheme NameType `xml:"cac:TaxScheme"`
}

type LegalMonetaryTotalType struct {
	XMLName              xml.Name   `xml:"cac:LegalMonetaryTotal"`
	LineExtensionAmount  AmountType `xml:"cbc:LineExtensionAmount"`
	TaxExclusiveAmount   AmountType `xml:"cbc:TaxExclusiveAmount"`
	TaxInclusiveAmount   AmountType `xml:"cbc:TaxInclusiveAmount"`
	AllowanceTotalAmount AmountType `xml:"cbc:AllowanceTotalAmount"`
	ChargeTotalAmount    AmountType `xml:"cbc:ChargeTotalAmount"`
	PayableAmount        AmountType `xml:"cbc:PayableAmount"`
}

type InvoiceLineType struct {
	XMLName               xml.Name            `xml:"cac:InvoiceLine"`
	ID                    string              `xml:"cbc:ID"`
	InvoicedQuantity      QuantityType        `xml:"cbc:InvoicedQuantity"`
	LineExtensionAmount   AmountType          `xml:"cbc:LineExtensionAmount"`
	FreeOfChargeIndicator bool                `xml:"cbc:FreeOfChargeIndicator"`
	AllowanceCharge       AllowanceChargeType `xml:"cac:AllowanceCharge"`
	TaxTotal              TaxTotalType        `xml:"cac:TaxTotal"`
	Item                  ItemType            `xml:",omitempty"`
	Price                 PriceType           `xml:",omitempty"`
}

type QuantityType struct {
	UnitCode string  `xml:"unitCode,attr"`
	Data     float64 `xml:"xml:",chardata""`
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

type InvoiceDate struct {
	FormatString string
	time.Time
}

func (t InvoiceDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	fmted := t.Format(t.FormatString)
	err := e.EncodeElement(fmted, start)
	if err != nil {
		return err
	}
	return nil
}
func (t *InvoiceDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(t.FormatString, v)
	if err != nil {
		return err
	}
	*t = InvoiceDate{t.FormatString, parse}
	return nil
}
