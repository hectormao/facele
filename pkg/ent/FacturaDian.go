package ent

import (
	"bytes"
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
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
	AgencyID           string = "195"
	AgencyName         string = "CO, DIAN (Direccion de Impuestos y Aduanas Nacionales)"
)

var NoResolucionesError error = errors.New("Resolucion Activa no existente")

func (invoice *InvoiceType) AgregarExtension(extension interface{}) {

	newExtension := UBLExtensionType{
		ExtensionContent: ExtensionContentType{
			Extension: extension,
		},
	}

	invoice.UBLExtensions.UBLExtensions =
		append(invoice.UBLExtensions.UBLExtensions, newExtension)
}

func NewInvoice(factura FacturaType, vendedor EmpresaType) (*InvoiceType, error) {

	resoluciones := factura.Empresa.Resoluciones
	var resolucion *ResolucionFacturacionType
	for _, res := range resoluciones {
		if res.SiActivo {
			resolucion = &res
			break
		}
	}
	if resolucion == nil {
		return nil, NoResolucionesError
	}

	totalImpuesto := calcularTotalImpuestos(factura)

	ublExtension := UBLExtensionType{
		ExtensionContent: ExtensionContentType{
			Extension: DianExtensionType{
				InvoiceControl: InvoiceControlType{
					InvoiceAuthorization: resolucion.Numero,
					AuthorizationPeriod: AuthorizationPeriodType{
						StartDate: invoiceDate{InvoiceDateFormat, resolucion.Vigencia.Desde},
						EndDate:   invoiceDate{InvoiceDateFormat, resolucion.Vigencia.Hasta},
					},
					AuthorizedInvoices: AuthorizedInvoicesType{
						Prefix: PrefixType{
							Type: DianTextType,
							Data: resolucion.Prefijo,
						},
						From: resolucion.Rango.Inferior,
						To:   resolucion.Rango.Superior,
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
						AgencyID,
						AgencyName,
						factura.Empresa.NumeroDocumento,
					),
					SoftwareID: newSoftwareID(
						AgencyID,
						AgencyName,
						factura.Empresa.SoftwareFacturacion.Id,
					),
				},
				SoftwareSecurityCode: newSoftwareSecurityCode(
					AgencyID,
					AgencyName,
					resolucion.ClaveTecnica,
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
		ID:           resolucion.Prefijo + strconv.Itoa(factura.CabezaFactura.Consecutivo),
		UUID: newUUID(
			AgencyID,
			AgencyName,
			generarCodigoCUFE(*resolucion, factura),
		),
		IssueDate: invoiceDate{InvoiceDateFormat, factura.CabezaFactura.FechaFacturacion},
		IssueTime: invoiceDate{InvoiceTimeFormat, factura.CabezaFactura.FechaFacturacion},
		InvoiceTypeCode: newInvoiceTypeCode(
			AgencyID,
			AgencyName,
			"1",
		),
		Note:                 factura.CabezaFactura.Observaciones,
		DocumentCurrencyCode: factura.CabezaFactura.Moneda,
		InvoicePeriod: InvoicePeriodType{
			DurationMeasure: DurationMeasureType{
				UnitCode: "DAY",
				Data: calcularDiasVencimiento(
					factura.CabezaFactura.FechaFacturacion,
					factura.CabezaFactura.FechaVencimiento,
				),
			},
			Description: factura.CabezaFactura.FechaVencimiento.Format(InvoiceDateFormat),
		},
		AccountingSupplierParty: AccountingSupplierPartyType{
			AdditionalAccountID: "1",
			Party: PartyType{
				PartyIdentification: PartyIdentificationType{
					ID: newPartyID(
						"31",
						AgencyID,
						AgencyName,
						vendedor.NumeroDocumento,
					),
				},
				PartyName: PartyNameType{
					Name: vendedor.RazonSocial,
				},
				PhysicalLocation: PhysicalLocationType{
					Address: AddressType{
						Department:          vendedor.Ubicacion.Departamento,
						CitySubdivisionName: "",
						CityName:            vendedor.Ubicacion.Municipio,
						AddressLine: AddressLineType{
							Line: []LineType{
								LineType{
									Data: vendedor.Ubicacion.Direccion,
								},
							},
						},
						Country: CountryType{
							IdentificationCode: vendedor.Ubicacion.Pais,
						},
					},
				},
				PartyTaxScheme: PartyTaxSchemeType{
					TaxLevelCode: strconv.Itoa(factura.CabezaFactura.TipoCompra),
					TaxScheme:    "",
				},
				PartyLegalEntity: PartyLegalEntityType{
					RegistrationName: vendedor.RazonSocial,
				},
			},
		},
		AccountingCustomerParty: AccountingCustomerPartyType{
			AdditionalAccountID: "1",
			Party: PartyType{
				PartyIdentification: PartyIdentificationType{
					ID: newPartyID(
						"31",
						AgencyID,
						AgencyName,
						factura.CabezaFactura.Nit,
					),
				},
				PartyName: PartyNameType{
					Name: factura.CabezaFactura.RazonSocial,
				},
				PhysicalLocation: PhysicalLocationType{
					Address: AddressType{
						Department:          factura.CabezaFactura.Departamento,
						CitySubdivisionName: "",
						CityName:            factura.CabezaFactura.Ciudad,
						AddressLine: AddressLineType{
							Line: []LineType{
								LineType{
									Data: factura.CabezaFactura.Direccion,
								},
							},
						},
						Country: CountryType{
							IdentificationCode: factura.CabezaFactura.Pais,
						},
					},
				},
				PartyTaxScheme: PartyTaxSchemeType{
					TaxLevelCode: "0",
					TaxScheme:    "",
				},
				PartyLegalEntity: PartyLegalEntityType{
					RegistrationName: factura.CabezaFactura.RazonSocial,
				},
			},
		},
		TaxTotal: TaxTotalType{
			TaxAmount: newTaxAmount(
				factura.CabezaFactura.Moneda,
				totalImpuesto,
			),
			TaxEvidenceIndicator: false,
			TaxSubtotal: generarSubtotalImpuestos(
				factura.CabezaFactura.ListaImpuestos.ImpuestosCabeza,
				factura.CabezaFactura.Moneda,
			),
		},
		LegalMonetaryTotal: LegalMonetaryTotalType{
			LineExtensionAmount: newLineExtensionAmount(
				factura.CabezaFactura.Moneda,
				factura.CabezaFactura.TotalImporteBruto,
			),
			TaxExclusiveAmount: newTaxExclusiveAmount(
				factura.CabezaFactura.Moneda,
				0.0,
			),
			PayableAmount: newPayableAmount(
				factura.CabezaFactura.Moneda,
				factura.CabezaFactura.TotalFactura,
			),
		},
		InvoiceLine: getInvoiceLine(factura),
	}

	return &invoice, nil
}

func generarSubtotalImpuestos(impuestos []ImpuestosCabezaType, moneda string) []TaxSubtotalType {

	result := make([]TaxSubtotalType, len(impuestos))

	for idx, impuesto := range impuestos {
		result[idx] = TaxSubtotalType{
			TaxableAmount: newTaxableAmount(
				moneda,
				impuesto.BaseImponible,
			),
			TaxAmount: newTaxAmount(
				moneda,
				impuesto.ValorImpuestoRetencion,
			),
			Percent: impuesto.Porcentaje,
			TaxCategory: TaxCategoryType{
				TaxScheme: TaxSchemeType{
					ID: impuesto.CodigoImpuestoRetencion,
				},
			},
		}
	}

	return result

}

func generarCodigoCUFE(resolucion ResolucionFacturacionType, factura FacturaType) string {

	fechaFormato := "20060102150405"

	mapImpuestos := getValoresImpuestos(factura.CabezaFactura.ListaImpuestos.ImpuestosCabeza)

	sha := sha1.New()

	numeroFactura := resolucion.Prefijo + strconv.Itoa(factura.CabezaFactura.Consecutivo)
	fechaFactura := factura.CabezaFactura.FechaFacturacion.Format(fechaFormato)
	valorFactura := fmt.Sprintf("%.2f", factura.CabezaFactura.TotalImporteBruto)
	codImp1 := "01"
	valImp1 := fmt.Sprintf("%.2f", mapImpuestos["01"])
	codImp2 := "02"
	valImp2 := fmt.Sprintf("%.2f", mapImpuestos["02"])
	codImp3 := "03"
	valImp3 := fmt.Sprintf("%.2f", mapImpuestos["03"])
	valImp := fmt.Sprintf("%.2f", factura.CabezaFactura.TotalFactura)
	nitOFE := factura.CabezaFactura.NumeroIdentificacion
	tipAdq := strconv.Itoa(factura.CabezaFactura.TipoPersona)
	numAdq := factura.CabezaFactura.Nit
	ciTec := resolucion.ClaveTecnica

	var buffer bytes.Buffer
	buffer.WriteString(numeroFactura)
	buffer.WriteString(fechaFactura)
	buffer.WriteString(valorFactura)
	buffer.WriteString(codImp1)
	buffer.WriteString(valImp1)
	buffer.WriteString(codImp2)
	buffer.WriteString(valImp2)
	buffer.WriteString(codImp3)
	buffer.WriteString(valImp3)
	buffer.WriteString(valImp)
	buffer.WriteString(nitOFE)
	buffer.WriteString(tipAdq)
	buffer.WriteString(numAdq)
	buffer.WriteString(ciTec)

	return fmt.Sprintf("%X", sha.Sum(buffer.Bytes()))
}

func getValoresImpuestos(impuestos []ImpuestosCabezaType) map[string]float64 {
	result := make(map[string]float64)

	for _, impuesto := range impuestos {
		result[impuesto.CodigoImpuestoRetencion] = impuesto.ValorImpuestoRetencion
	}

	return result
}

func calcularDiasVencimiento(inicio time.Time, fin time.Time) int64 {
	diaInicio := inicioDia(inicio)
	diaFin := inicioDia(fin)

	return int64(diaFin.Sub(diaInicio).Hours() / 24)
}

func inicioDia(t time.Time) time.Time {
	anio, mes, dia := t.Date()
	return time.Date(anio, mes, dia, 0, 0, 0, 0, t.Location())
}

func getInvoiceLine(factura FacturaType) []InvoiceLineType {

	lines := factura.CabezaFactura.ListaDetalles.Detalles

	result := make([]InvoiceLineType, len(lines))

	for idx, line := range lines {
		result[idx] = InvoiceLineType{
			ID:                  line.CodigoProducto,
			InvoicedQuantity:    line.Cantidad,
			LineExtensionAmount: newLineExtensionAmount(factura.CabezaFactura.Moneda, line.PrecioSinImpuestos),
			Item: ItemType{
				Description: line.Descripcion,
			},
			Price: PriceType{
				PriceAmount: newPriceAmountType(factura.CabezaFactura.Moneda, line.ValorUnitario),
			},
		}
	}

	return result
}

func calcularTotalImpuestos(factura FacturaType) float64 {
	total := 0.0
	for _, impuesto := range factura.CabezaFactura.ListaImpuestos.ImpuestosCabeza {
		total += impuesto.ValorImpuestoRetencion
	}
	return total
}

func newPriceAmountType(
	currencyID string,
	data float64) PriceAmountType {
	v := PriceAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func newPayableAmount(
	currencyID string,
	data float64) PayableAmountType {
	v := PayableAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func newTaxExclusiveAmount(
	currencyID string,
	data float64) TaxExclusiveAmountType {
	v := TaxExclusiveAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func newLineExtensionAmount(
	currencyID string,
	data float64) LineExtensionAmountType {
	v := LineExtensionAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v

}

func newTaxableAmount(
	currencyID string,
	data float64) TaxableAmountType {
	v := TaxableAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func newTaxAmount(
	currencyID string,
	data float64) TaxAmountType {
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
	InvoicePeriod           InvoicePeriodType           `xml:",omitempty"`
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
	From    int64      `xml:"sts:From"`
	To      int64      `xml:"sts:To"`
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

type InvoicePeriodType struct {
	XMLName         xml.Name            `xml:"cac:InvoicePeriod"`
	DurationMeasure DurationMeasureType `xml:",omiteempty"`
	Description     string              `xml:"cbc:Description,omiteempty"`
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
	XMLName              xml.Name          `xml:"fe:TaxTotal"`
	TaxAmount            TaxAmountType     `xml:",omitempty"`
	TaxEvidenceIndicator bool              `xml:"cbc:TaxEvidenceIndicator"`
	TaxSubtotal          []TaxSubtotalType `xml:",omitempty"`
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
