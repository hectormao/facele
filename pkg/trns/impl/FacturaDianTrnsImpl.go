package impl

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"strconv"
	"time"

	"github.com/hectormao/facele/pkg/ent"
	"github.com/hectormao/facele/pkg/ssl"
	trnsConfig "github.com/hectormao/facele/pkg/trns/cfg"

	"log"
	"text/template"
)

type FacturaDianTrnsImpl struct {
	Config      trnsConfig.FacturaDianTrnsConfig
	URLTemplate *template.Template
}

func (trns FacturaDianTrnsImpl) FacturaToInvoice(factura ent.FacturaType) (ent.InvoiceType, error) {

	invoice, err := trns.newInvoice(factura)
	if err != nil {
		return ent.InvoiceType{}, err
	}

	sign, err := ssl.FirmarXML(invoice)
	if err != nil {
		return ent.InvoiceType{}, err
	}

	invoice.AgregarExtension(sign)

	return invoice, nil
}

func (trns FacturaDianTrnsImpl) newInvoice(factura ent.FacturaType) (ent.InvoiceType, error) {

	idFactura := factura.Resolucion.Prefijo + strconv.Itoa(factura.CabezaFactura.Consecutivo)

	cufe := trns.generarCodigoCUFE(factura)
	factura.Cufe = cufe

	adquirientePrincipal := factura.CabezaFactura.ListaAdquirentes[0]

	ublExtension := ent.UBLExtensionType{
		ExtensionContent: ent.ExtensionContentType{
			Extension: ent.DianExtensionType{
				InvoiceControl: ent.InvoiceControlType{
					InvoiceAuthorization: factura.Resolucion.Numero,
					AuthorizationPeriod: ent.AuthorizationPeriodType{
						StartDate: ent.InvoiceDate{ent.InvoiceDateFormat, factura.Resolucion.Vigencia.Desde},
						EndDate:   ent.InvoiceDate{ent.InvoiceDateFormat, factura.Resolucion.Vigencia.Hasta},
					},
					AuthorizedInvoices: ent.AuthorizedInvoicesType{
						Prefix: factura.Resolucion.Prefijo,
						From:   factura.Resolucion.Rango.Inferior,
						To:     factura.Resolucion.Rango.Superior,
					},
				},
				InvoiceSource: ent.InvoiceSourceType{
					IdentificationCode: trns.newIdentificationCode(
						ent.ListAgencyID,
						ent.ListAgencyName,
						ent.CountrySchemeURI,
						factura.Empresa.Ubicacion.Pais.Codigo,
					),
				},
				SoftwareProvider: ent.SoftwareProviderType{
					ProviderID: trns.newProviderID(
						ent.AgencyID,
						ent.AgencyName,
						ent.SchemeID,
						ent.SchemeName,
						factura.Empresa.NumeroDocumento,
					),
					SoftwareID: trns.newSoftwareID(
						ent.AgencyID,
						ent.AgencyName,
						factura.Empresa.SoftwareFacturacion.Id,
					),
				},
				SoftwareSecurityCode: trns.newSoftwareSecurityCode(
					ent.AgencyID,
					ent.AgencyName,
					factura.Resolucion.ClaveTecnica,
				),
				AuthorizationProvider: trns.newAuthorizationProviderID(
					ent.AgencyID,
					ent.AgencyName,
					ent.SchemeID4,
					ent.SchemeName,
					factura.Empresa.NumeroDocumento,
				),
				QRCode: trns.qrCodeURL(factura),
			},
		},
	}

	invoice := ent.InvoiceType{
		XmlNS:             ent.DianNS,
		XmlNSCac:          ent.DianNSCac,
		XmlNSCbc:          ent.DianNSCbc,
		XmlNSExt:          ent.DianNSExt,
		XmlNSSts:          ent.DianNSSts,
		XmlNSXades:        ent.DianNSXades,
		XmlNSXades141:     ent.DianNSXades141,
		XmlNSXsi:          ent.DianNSXsi,
		XsiSchemaLocation: ent.DianSchemaLocation,
		UBLExtensions: ent.UBLExtensionsType{
			UBLExtensions: []ent.UBLExtensionType{ublExtension},
		},
		UBLVersionID:       ent.UBLVersion,
		CustomizationID:    ent.Customization,
		ProfileID:          ent.Profile,
		ProfileExecutionID: ent.ProfileExecution,
		ID:                 "BC100201",
		UUID: trns.newUUID(
			"2",
			"CUFE-SHA256",
			cufe,
		),
		IssueDate:            ent.InvoiceDate{ent.InvoiceDateFormat, factura.CabezaFactura.FechaFacturacion},
		IssueTime:            ent.InvoiceDate{ent.InvoiceTimeFormat, factura.CabezaFactura.FechaFacturacion},
		DueDate:              ent.InvoiceDate{ent.InvoiceTimeFormat, factura.CabezaFactura.Pago.FechaVencimiento},
		InvoiceTypeCode:      strconv.Itoa(factura.CabezaFactura.TipoDocumento),
		Note:                 trns.getCampoAdicionalPorNombre(factura, "OBSERVACIONES"),
		TaxPointDate:         ent.InvoiceDate{ent.InvoiceDateFormat, factura.CabezaFactura.Pago.FechaVencimiento},
		DocumentCurrencyCode: factura.CabezaFactura.Pago.Moneda,
		LineCountNumeric:     len(factura.CabezaFactura.ListaDetalles),
		OrderReference: ent.ReferenceType{
			ID: idFactura,
		},
		DespatchDocumentReference: ent.ReferenceType{
			ID: idFactura,
		},
		ReceiptDocumentReference: ent.ReferenceType{
			ID: idFactura,
		},
		AdditionalDocumentReference: ent.AdditionalDocumentReferenceType{
			ID:               idFactura,
			DocumentTypeCode: strconv.Itoa(factura.CabezaFactura.TipoDocumento),
		},
		AccountingSupplierParty: ent.AccountingSupplierPartyType{
			AdditionalAccountID: ent.IDType{
				SchemeAgencyID: ent.AgencyID,
				Data:           factura.Empresa.Tipo,
			},
			Party: ent.PartyType{
				PartyName: []ent.PartyNameType{
					ent.PartyNameType{
						Name: factura.Empresa.RazonSocial,
					},
				},
				PhysicalLocation: ent.PhysicalLocationType{
					Address: ent.AddressType{
						ID:                   factura.Empresa.Ubicacion.Municipio.Codigo,
						CityName:             factura.Empresa.Ubicacion.Municipio.Nombre,
						CountrySubentity:     factura.Empresa.Ubicacion.Departamento.Nombre,
						CountrySubentityCode: factura.Empresa.Ubicacion.Departamento.Codigo,
						AddressLine: ent.AddressLineType{
							Line: []ent.LineType{
								ent.LineType{
									Data: factura.Empresa.Ubicacion.Direccion,
								},
							},
						},
						Country: ent.CountryType{
							IdentificationCode: factura.Empresa.Ubicacion.Pais.Codigo,
							Name:               factura.Empresa.Ubicacion.Pais.Nombre,
						},
					},
				},
				PartyTaxScheme: ent.PartyTaxSchemeType{
					RegistrationName: factura.Empresa.RazonSocial,
					CompanyID: trns.newID(
						ent.AgencyID,
						ent.AgencyName,
						ent.SchemeID,
						ent.SchemeName,
						factura.Empresa.NumeroDocumento,
					),
					TaxLevelCode: ent.TaxLevelCodeType{
						ListName: ent.ListName,
						Data:     "0-11",
					},
					RegistrationAddress: ent.AddressType{
						ID:                   factura.Empresa.Ubicacion.Municipio.Codigo,
						CityName:             factura.Empresa.Ubicacion.Municipio.Nombre,
						CountrySubentity:     factura.Empresa.Ubicacion.Departamento.Nombre,
						CountrySubentityCode: factura.Empresa.Ubicacion.Departamento.Codigo,
						AddressLine: ent.AddressLineType{
							Line: []ent.LineType{
								ent.LineType{
									Data: factura.Empresa.Ubicacion.Direccion,
								},
							},
						},
						Country: ent.CountryType{
							IdentificationCode: factura.Empresa.Ubicacion.Pais.Codigo,
							Name:               factura.Empresa.Ubicacion.Pais.Nombre,
						},
					},
					TaxScheme: ent.NameType{
						ID:   "01",
						Name: "IVA",
					},
				},
				PartyLegalEntity: ent.PartyLegalEntityType{
					RegistrationName: factura.Empresa.RazonSocial,
					CompanyID: trns.newID(
						ent.AgencyID,
						ent.AgencyName,
						ent.SchemeID,
						ent.SchemeName,
						factura.Empresa.NumeroDocumento,
					),
					CorporateRegistrationScheme: ent.NameType{
						ID:   "BC",
						Name: "12345",
					},
				},
				Contact: ent.ContactType{
					Telephone:      factura.Empresa.Contacto.Telefonos[0],
					ElectronicMail: factura.Empresa.Contacto.Correos[0],
				},
			},
		},
		AccountingCustomerParty: ent.AccountingCustomerPartyType{
			AdditionalAccountID: ent.IDType{
				Data: adquirientePrincipal.TipoPersona,
			},
			Party: ent.PartyType{
				PartyName: []ent.PartyNameType{
					ent.PartyNameType{
						Name: adquirientePrincipal.NombreCompleto,
					},
				},
				PhysicalLocation: ent.PhysicalLocationType{
					Address: ent.AddressType{
						ID:                   adquirientePrincipal.Ciudad,
						CityName:             adquirientePrincipal.DescripcionCiudad,
						CountrySubentity:     adquirientePrincipal.Departamento,
						CountrySubentityCode: adquirientePrincipal.NombreDepartamento,
						AddressLine: ent.AddressLineType{
							Line: []ent.LineType{
								ent.LineType{
									Data: adquirientePrincipal.Direccion,
								},
							},
						},
						Country: ent.CountryType{
							IdentificationCode: adquirientePrincipal.Pais,
							Name:               adquirientePrincipal.NombrePais,
						},
					},
				},
			},
		},
		InvoiceLine: trns.getInvoiceLine(factura),
	}

	return invoice, nil
}
func (trns *FacturaDianTrnsImpl) qrCodeURL(factura ent.FacturaType) string {
	if trns.URLTemplate == nil {
		rawUrl := trns.Config.GetQRCodeUrl()
		template, err := template.New("URLQRCode").Parse(rawUrl)
		if err != nil {
			log.Printf("Error al parsear URL QRCODE %v", err)
			trns.URLTemplate = nil
			return ""
		}
		trns.URLTemplate = template
	}
	buf := new(bytes.Buffer)
	err := trns.URLTemplate.Execute(buf, factura)
	if err != nil {
		log.Printf("Error al parsear URL QRCODE (ejecutando template) %v", err)
		return ""
	}

	return buf.String()
}

func (trns FacturaDianTrnsImpl) getCampoAdicionalPorNombre(factura ent.FacturaType, nombreCampo string) string {
	camposAdicionales := factura.
		CabezaFactura.
		ListaCamposAdicionales

	for _, campo := range camposAdicionales {
		if campo.NombreCampo == nombreCampo {
			return campo.ValorCampo
		}
	}

	return ""

}

func (trns FacturaDianTrnsImpl) generarCodigoCUFE(factura ent.FacturaType) string {

	impuestos := make(map[string]ent.ImpuestoType)
	for _, impuesto := range factura.CabezaFactura.ListaImpuestos {
		impuestos[impuesto.CodigoImpuestoRetencion] = impuesto
	}

	fechaFormato := "20060102150405"

	sha := sha1.New()

	impuesto01, impuesto01OK := impuestos["01"]
	impuesto02, impuesto02OK := impuestos["02"]
	impuesto03, impuesto03OK := impuestos["03"]

	valorImpuesto01 := 0.0
	valorImpuesto02 := 0.0
	valorImpuesto03 := 0.0

	if impuesto01OK {
		valorImpuesto01 = impuesto01.ValorImpuestoRetencion
	}

	if impuesto02OK {
		valorImpuesto02 = impuesto02.ValorImpuestoRetencion
	}

	if impuesto03OK {
		valorImpuesto03 = impuesto03.ValorImpuestoRetencion
	}

	adquirientePrincipal := factura.CabezaFactura.ListaAdquirentes[0]

	numeroFactura := factura.Resolucion.Prefijo + strconv.Itoa(factura.CabezaFactura.Consecutivo)
	fechaFactura := factura.CabezaFactura.FechaFacturacion.Format(fechaFormato)
	valorFactura := fmt.Sprintf("%.2f", factura.CabezaFactura.Pago.TotalImporteBruto)
	codImp1 := "01"
	valImp1 := fmt.Sprintf("%.2f", valorImpuesto01)
	codImp2 := "02"
	valImp2 := fmt.Sprintf("%.2f", valorImpuesto02)
	codImp3 := "03"
	valImp3 := fmt.Sprintf("%.2f", valorImpuesto03)
	valImp := fmt.Sprintf("%.2f", factura.CabezaFactura.Pago.TotalFactura)
	nitOFE := factura.CabezaFactura.Nit
	tipAdq := adquirientePrincipal.TipoPersona
	numAdq := adquirientePrincipal.NumeroIdentificacion
	ciTec := factura.Resolucion.ClaveTecnica

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

func (trns FacturaDianTrnsImpl) calcularDiasVencimiento(inicio time.Time, fin time.Time) int64 {
	diaInicio := trns.inicioDia(inicio)
	diaFin := trns.inicioDia(fin)

	return int64(diaFin.Sub(diaInicio).Hours() / 24)
}

func (trns FacturaDianTrnsImpl) inicioDia(t time.Time) time.Time {
	anio, mes, dia := t.Date()
	return time.Date(anio, mes, dia, 0, 0, 0, 0, t.Location())
}

func (trns FacturaDianTrnsImpl) getInvoiceLine(factura ent.FacturaType) []ent.InvoiceLineType {

	lines := factura.CabezaFactura.ListaDetalles

	result := make([]ent.InvoiceLineType, len(lines))

	for idx, line := range lines {

		impuesto := line.ListaImpuestos[0]
		cargo := line.ListaCargos[0]

		result[idx] = ent.InvoiceLineType{
			ID: line.CodigoProducto,
			InvoicedQuantity: ent.QuantityType{
				UnitCode: line.UnidadMedida,
				Data:     float64(line.Cantidad),
			},
			LineExtensionAmount: ent.AmountType{
				CurrencyID: factura.CabezaFactura.Pago.Moneda,
				Data:       line.PrecioTotal,
			},
			FreeOfChargeIndicator: !(line.PrecioTotal > 0),
			PricingReference: ent.PricingReferenceType{
				AlternativeConditionPrice: ent.AlternativeConditionPriceType{
					PriceAmount: ent.AmountType{
						CurrencyID: factura.CabezaFactura.Pago.Moneda,
						Data:       line.PrecioSinImpuestos,
					},
					PriceTypeCode: line.CodigoTipo,
				},
			},
			TaxTotal: ent.TaxTotalType{
				TaxAmount: ent.AmountType{
					CurrencyID: factura.CabezaFactura.Pago.Moneda,
					Data:       impuesto.BaseImponible,
				},
				TaxEvidenceIndicator: impuesto.IsAutoRetenido,
				TaxSubtotal: ent.TaxSubtotalType{
					TaxableAmount: ent.AmountType{
						CurrencyID: factura.CabezaFactura.Pago.Moneda,
						Data:       impuesto.BaseImponible,
					},
					TaxAmount: ent.AmountType{
						CurrencyID: factura.CabezaFactura.Pago.Moneda,
						Data:       impuesto.ValorImpuestoRetencion,
					},
					TaxCategory: ent.TaxCategoryType{
						Percent: impuesto.Porcentaje,
						TaxScheme: ent.NameType{
							ID:   impuesto.CodigoImpuestoRetencion,
							Name: cargo.Descripcion,
						},
					},
				},
			},
			Item: ent.ItemType{
				Description: line.Descripcion,
				SellersItemIdentification: ent.IdentificationType{
					ID: line.CodigoProducto,
				},
				StandardItemIdentification: ent.IdentificationType{
					ID: line.TipoCodigoProducto,
				},
			},
			Price: ent.PriceType{
				PriceAmount: ent.AmountType{
					CurrencyID: factura.CabezaFactura.Pago.Moneda,
					Data:       line.ValorUnitario,
				},
				BaseQuantity: ent.QuantityType{
					UnitCode: line.UnidadMedida,
					Data:     float64(line.Cantidad),
				},
			},
		}
	}

	return result
}

func (trns FacturaDianTrnsImpl) newIdentificationCode(
	listAgencyID string,
	listAgencyName string,
	listSchemeURI string,
	data string) ent.IdentificationCodeType {
	v := ent.IdentificationCodeType{}
	v.ListAgencyID = listAgencyID
	v.ListAgencyName = listAgencyName
	v.ListSchemeURI = listSchemeURI
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newID(
	schemeAgencyID string,
	schemeAgencyName string,
	schemeID string,
	schemeName string,
	data string) ent.IDType {
	v := ent.IDType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.SchemeID = schemeID
	v.SchemeName = schemeName
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newProviderID(
	schemeAgencyID string,
	schemeAgencyName string,
	schemeID string,
	schemeName string,
	data string) ent.ProviderIDType {
	v := ent.ProviderIDType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.SchemeID = schemeID
	v.SchemeName = schemeName
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newSoftwareID(
	schemeAgencyID string,
	schemeAgencyName string,
	data string) ent.SoftwareIDType {
	v := ent.SoftwareIDType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newSoftwareSecurityCode(
	schemeAgencyID string,
	schemeAgencyName string,
	data string) ent.SoftwareSecurityCodeType {
	v := ent.SoftwareSecurityCodeType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newAuthorizationProviderID(
	schemeAgencyID string,
	schemeAgencyName string,
	schemeID string,
	schemeName string,
	data string) ent.AuthorizationProviderType {
	v := ent.AuthorizationProviderIDType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.SchemeID = schemeID
	v.SchemeName = schemeName
	v.Data = data
	return ent.AuthorizationProviderType{
		AuthorizationProviderID: v,
	}
}

func (trns FacturaDianTrnsImpl) newUUID(
	schemeID string,
	schemeName string,
	data string) ent.UUIDType {
	v := ent.UUIDType{}
	v.SchemeID = schemeID
	v.SchemeName = schemeName
	v.Data = data
	return v
}
