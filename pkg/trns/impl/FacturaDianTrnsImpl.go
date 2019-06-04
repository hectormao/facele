package impl

import (
	"strconv"
	"time"

	"github.com/hectormao/facele/pkg/ent"
	"github.com/hectormao/facele/pkg/ssl"
)

type FacturaDianTrnsImpl struct {
	genericodes map[string]map[string]ent.Genericode
}

func (trns FacturaDianTrnsImpl) FacturaToInvoice(factura ent.FacturaType, resolucion ent.ResolucionFacturacionType, vendedor ent.EmpresaType, genericodes map[string]map[string]ent.Genericode) (ent.InvoiceType, error) {

	invoice, err := trns.newInvoice(factura, vendedor, resolucion, genericodes)
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

func (trns FacturaDianTrnsImpl) newInvoice(factura ent.FacturaType, vendedor ent.EmpresaType, resolucion ent.ResolucionFacturacionType, genericodes map[string]map[string]ent.Genericode) (ent.InvoiceType, error) {

	impuestosFactura := trns.getImpuestosFactura(factura)
	totalImpuesto := trns.calcularTotalImpuestos(impuestosFactura)

	idFactura := resolucion.Prefijo + strconv.Itoa(factura.CabezaFactura.Consecutivo)

	ublExtension := ent.UBLExtensionType{
		ExtensionContent: ent.ExtensionContentType{
			Extension: ent.DianExtensionType{
				InvoiceControl: ent.InvoiceControlType{
					InvoiceAuthorization: resolucion.Numero,
					AuthorizationPeriod: ent.AuthorizationPeriodType{
						StartDate: ent.InvoiceDate{ent.InvoiceDateFormat, resolucion.Vigencia.Desde},
						EndDate:   ent.InvoiceDate{ent.InvoiceDateFormat, resolucion.Vigencia.Hasta},
					},
					AuthorizedInvoices: ent.AuthorizedInvoicesType{
						Prefix: resolucion.Prefijo,
						From:   resolucion.Rango.Inferior,
						To:     resolucion.Rango.Superior,
					},
				},
				InvoiceSource: ent.InvoiceSourceType{
					IdentificationCode: trns.newIdentificationCode(
						"6",
						"United Nations Economic Commission for Europe",
						ent.CountrySchemeURI,
						"CO",
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
					resolucion.ClaveTecnica,
				),
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
			trns.generarCodigoCUFE(resolucion, factura, impuestosFactura),
		),
		IssueDate:            ent.InvoiceDate{ent.InvoiceDateFormat, factura.CabezaFactura.FechaFacturacion},
		IssueTime:            ent.InvoiceDate{ent.InvoiceTimeFormat, factura.CabezaFactura.FechaFacturacion},
		InvoiceTypeCode:      strconv.Itoa(factura.CabezaFactura.TipoDocumento),
		Note:                 factura.CabezaFactura.Observaciones,
		TaxPointDate:         ent.InvoiceDate{ent.InvoiceDateFormat, factura.CabezaFactura.FechaVencimiento},
		DocumentCurrencyCode: factura.CabezaFactura.Moneda,
		LineCountNumeric:     len(factura.CabezaFactura.ListaDetalles.Detalles),
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
				Data:           vendedor.Tipo,
			},
			Party: ent.PartyType{
				PartyName: ent.NameType{
					Name: vendedor.RazonSocial,
				},
				PhysicalLocation: ent.PhysicalLocationType{
					Address: ent.AddressType{
						ID:                   "",
						CityName:             "",
						CountrySubentity:     factura.CabezaFactura.Ciudad,
						CountrySubentityCode: "",
						AddressLine: AddressLineType{
							Line: []LineType{
								LineType{
									Data: vendedor.Ubicacion.Direccion,
								},
							},
						},
						Country: CountryType{
							IdentificationCode: vendedor.Ubicacion.Pais,
							Name:               "",
						},
					},
				},
				PartyTaxScheme: PartyTaxSchemeType{
					RegistrationName: vendedor.RazonSocial,
				},
				PartyLegalEntity: PartyLegalEntityType{
					RegistrationName: vendedor.RazonSocial,
				},
			},
		},
		AccountingCustomerParty: AccountingCustomerPartyType{
			AdditionalAccountID: strconv.Itoa(factura.CabezaFactura.TipoPersona),
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
					TaxLevelCode: strconv.Itoa(factura.CabezaFactura.TipoCompra),
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
				impuestosFactura,
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

	return invoice, nil
}

func (trns FacturaDianTrnsImpl) generarSubtotalImpuestos(impuestos map[string]*ent.ImpuestosCabezaType, moneda string) []ent.TaxSubtotalType {

	result := make([]TaxSubtotalType, len(impuestos))
	idx := 0
	for _, impuesto := range impuestos {
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
		idx++
	}

	return result

}

func (trns FacturaDianTrnsImpl) generarCodigoCUFE(resolucion ent.ResolucionFacturacionType, factura ent.FacturaType, impuestos map[string]*ent.ImpuestosCabezaType) string {

	fechaFormato := "20060102150405"

	sha := sha1.New()

	impuesto01 := impuestos["01"]
	impuesto02 := impuestos["02"]
	impuesto03 := impuestos["03"]

	valorImpuesto01 := 0.0
	valorImpuesto02 := 0.0
	valorImpuesto03 := 0.0

	if impuesto01 != nil {
		valorImpuesto01 = impuesto01.ValorImpuestoRetencion
	}

	if impuesto02 != nil {
		valorImpuesto02 = impuesto02.ValorImpuestoRetencion
	}

	if impuesto03 != nil {
		valorImpuesto03 = impuesto03.ValorImpuestoRetencion
	}

	numeroFactura := resolucion.Prefijo + strconv.Itoa(factura.CabezaFactura.Consecutivo)
	fechaFactura := factura.CabezaFactura.FechaFacturacion.Format(fechaFormato)
	valorFactura := fmt.Sprintf("%.2f", factura.CabezaFactura.TotalImporteBruto)
	codImp1 := "01"
	valImp1 := fmt.Sprintf("%.2f", valorImpuesto01)
	codImp2 := "02"
	valImp2 := fmt.Sprintf("%.2f", valorImpuesto02)
	codImp3 := "03"
	valImp3 := fmt.Sprintf("%.2f", valorImpuesto03)
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

func (trns FacturaDianTrnsImpl) getValoresImpuestos(impuestos []ent.ImpuestosCabezaType) map[string]float64 {
	result := make(map[string]float64)

	for _, impuesto := range impuestos {
		result[impuesto.CodigoImpuestoRetencion] = impuesto.ValorImpuestoRetencion
	}

	return result
}

func (trns FacturaDianTrnsImpl) calcularDiasVencimiento(inicio time.Time, fin time.Time) int64 {
	diaInicio := inicioDia(inicio)
	diaFin := inicioDia(fin)

	return int64(diaFin.Sub(diaInicio).Hours() / 24)
}

func (trns FacturaDianTrnsImpl) inicioDia(t time.Time) time.Time {
	anio, mes, dia := t.Date()
	return time.Date(anio, mes, dia, 0, 0, 0, 0, t.Location())
}

func (trns FacturaDianTrnsImpl) getInvoiceLine(factura ent.FacturaType) []ent.InvoiceLineType {

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

func (trns FacturaDianTrnsImpl) getImpuestosFactura(factura ent.FacturaType) map[string]*ent.ImpuestosCabezaType {
	result := make(map[string]*ImpuestosCabezaType)
	for _, detalle := range factura.CabezaFactura.ListaDetalles.Detalles {
		for _, impuesto := range detalle.ListaImpuestos.Detalles {
			impuestoGeneral, ok := result[impuesto.CodigoImpuestoRetencion]
			if ok {
				impuestoGeneral.BaseImponible += impuesto.BaseImponible
				impuestoGeneral.ValorImpuestoRetencion += impuesto.ValorImpuestoRetencion
			} else {
				result[impuesto.CodigoImpuestoRetencion] = &ImpuestosCabezaType{
					BaseImponible:           impuesto.BaseImponible,
					CodigoImpuestoRetencion: impuesto.CodigoImpuestoRetencion,
					Porcentaje:              impuesto.Porcentaje,
					ValorImpuestoRetencion:  impuesto.ValorImpuestoRetencion,
				}
			}
		}
	}

	return result
}

func (trns FacturaDianTrnsImpl) calcularTotalImpuestos(impuestos map[string]*ent.ImpuestosCabezaType) float64 {
	total := 0.0
	for _, impuesto := range impuestos {
		total += impuesto.ValorImpuestoRetencion
	}
	return total
}

func (trns FacturaDianTrnsImpl) newPriceAmountType(
	currencyID string,
	data float64) ent.PriceAmountType {
	v := PriceAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newPayableAmount(
	currencyID string,
	data float64) ent.PayableAmountType {
	v := PayableAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newTaxExclusiveAmount(
	currencyID string,
	data float64) ent.TaxExclusiveAmountType {
	v := TaxExclusiveAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newLineExtensionAmount(
	currencyID string,
	data float64) ent.LineExtensionAmountType {
	v := LineExtensionAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v

}

func (trns FacturaDianTrnsImpl) newTaxableAmount(
	currencyID string,
	data float64) ent.TaxableAmountType {
	v := TaxableAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newTaxAmount(
	currencyID string,
	data float64) ent.TaxAmountType {
	v := TaxAmountType{}
	v.CurrencyID = currencyID
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newIdentificationCode(
	listAgencyID string,
	listAgencyName string,
	listSchemeURI string,
	data string) ent.IdentificationCodeType {
	v := IdentificationCodeType{}
	v.ListAgencyID = listAgencyID
	v.ListAgencyName = listAgencyName
	v.ListSchemeURI = listSchemeURI
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newProviderID(
	schemeAgencyID string,
	schemeAgencyName string,
	schemeID string,
	schemeName string,
	data string) ent.ProviderIDType {
	v := ProviderIDType{}
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
	v := SoftwareIDType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newSoftwareSecurityCode(
	schemeAgencyID string,
	schemeAgencyName string,
	data string) ent.SoftwareSecurityCodeType {
	v := SoftwareSecurityCodeType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newUUID(
	schemeID string,
	schemeName string,
	data string) ent.UUIDType {
	v := UUIDType{}
	v.SchemeID = schemeID
	v.SchemeName = schemeName
	v.Data = data
	return v
}
func (trns FacturaDianTrnsImpl) newInvoiceTypeCode(
	listAgencyID string,
	listAgencyName string,
	data string) ent.InvoiceTypeCodeType {
	v := InvoiceTypeCodeType{}
	v.ListAgencyID = listAgencyID
	v.ListAgencyName = listAgencyName
	v.ListSchemeURI = DianSchemeURI
	v.Data = data
	return v
}

func (trns FacturaDianTrnsImpl) newPartyID(
	schemeID string,
	schemeAgencyID string,
	schemeAgencyName string,
	data string) ent.PartyIDType {
	v := PartyIDType{}
	v.SchemeAgencyID = schemeAgencyID
	v.SchemeAgencyName = schemeAgencyName
	v.Data = data
	v.SchemeID = schemeID
	return v
}
