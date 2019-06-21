package ent

import (
	"encoding/xml"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type EmpresaType struct {
	ObjectId               objectid.ObjectID           `bson:"_id" json:"_id"`
	Tipo                   string                      `bson:"tipo" json:"tipo"`
	TipoFacturaElectronica string                      `bson:"tipo_factura_electronica" json:"tipo_factura_electronica"`
	NumeroDocumento        string                      `bson:"numero_documento" json:"numero_documento"`
	RazonSocial            string                      `bson:"razon_social" json:"razon_social"`
	Ubicacion              UbicacionType               `bson:"ubicacion" json:"ubicacion"`
	Contacto               ContactoType                `bson:"contacto" json:"contacto"`
	Resoluciones           []ResolucionFacturacionType `bson:"resoluciones" json:"resoluciones"`
	Dominio                []DominioType               `bson:"dominio" json:"dominio"`
	SoftwareFacturacion    SoftwareFacturacionType     `bson:"software_facturacion" json:"software_facturacion"`
}

type FacturaType struct {
	XMLName       xml.Name          `xml:"listaDocumentos" bson:"-" json:"-"`
	CabezaFactura CabezaFacturaType `xml:"cabezaDocumento" bson:"cabeza_documento" json:"cabeza_documento"`
	EmpresaID     string            `xml:"-" bson:"_empresa_id" json:"_empresa_id"`
	Empresa       EmpresaType       `xml:"-" bson:"-" json:"_empresa"`
	ObjectId      string            `xml:"-" bson:"-" json:"_id"`
	Cufe          string            `xml:"-" bson:"-" json:"_cufe"`
}

type CabezaFacturaType struct {
	TipoDocumento          int                  `xml:"Tipodocumento" bson:"tipo_documento" json:"tipo_documento"`
	Prefijo                string               `xml:"Prefijo" bson:"prefijo" json:"prefijo"`
	Consecutivo            int                  `xml:"Consecutivo" bson:"consecutivo" json:"consecutivo"`
	FechaFacturacion       time.Time            `xml:"Fechafacturacion" bson:"fecha_facturacion" json:"fecha_facturacion"`
	Nit                    string               `xml:"Nit" bson:"nit" json:"nit"`
	AplicaFel              bool                 `xml:"Aplicafel" bson:"aplica_fel" json:"aplica_fel"`
	Pago                   PagoType             `xml:"Pago" bson:"pago" json:"pago"`
	ListaDetalles          []DetalleType        `xml:"listaDetalle" bson:"lista_detalle" json:"lista_detalle"`
	ListaImpuestos         []ImpuestoType       `xml:"listaImpuestos" bson:"lista_impuestos" json:"lista_impuestos"`
	ListaAdquirentes       []AdquirienteType    `xml:"listaAdquirentes" bson:"lista_adquirentes" json:"lista_adquirentes"`
	ListaCamposAdicionales []CampoAdicionalType `xml:"listaCamposAdicionales" bson:"lista_adicionales" json:"lista_adicionales"`
	ListaCargos            []CargoType          `xml:"listaCargos" bson:"lista_cargos" json:"lista_cargos"`
	ListaMediosPagos       []MedioPagoType      `xml:"listaMediosPagos" bson:"lista_medios_pagos" json:"lista_medios_pagos"`
}

type MedioPagoType struct {
	XMLName   xml.Name `xml:"MedioPago" bson:"-" json:"-"`
	MedioPago string   `xml:"medioPago" bson:"medio_pago" json:"medio_pago"`
}

type AdquirienteType struct {
	XMLName                 xml.Name `xml:"Adquiriente" bson:"-" json:"-"`
	TipoPersona             string   `xml:"tipopersona" bson:"tipo_persona" json:"tipo_persona"`
	NombreCompleto          string   `xml:"nombrecompleto" bson:"nombre_completo" json:"nombre_completo"`
	TipoIdentificacion      string   `xml:"tipoidentificacion" bson:"tipo_identificacion" json:"tipo_identificacion"`
	NumeroIdentificacion    string   `xml:"numeroidentificacion" bson:"numero_identificacion" json:"numero_identificacion"`
	Regimen                 string   `xml:"regimen" bson:"regimen" json:"regimen"`
	Email                   string   `xml:"Email" bson:"email" json:"email"`
	Pais                    string   `xml:"pais" bson:"pais" json:"pais"`
	Departamento            string   `xml:"departamento" bson:"departamento" json:"departamento"`
	NombreDepartamento      string   `xml:"nombredepartamento" bson:"nombre_departamento" json:"nombre_departamento"`
	Ciudad                  string   `xml:"ciudad" bson:"ciudad" json:"ciudad"`
	DescripcionCiudad       string   `xml:"descripcionciudad" bson:"descripcion_ciudad" json:"descripcion_ciudad"`
	Direccion               string   `xml:"direccion" bson:"pdireccion" json:"direccion"`
	Telefono                string   `xml:"telefono" bson:"telefono" json:"telefono"`
	EnvioPorEmailPlataforma string   `xml:"envioPorEmailPlataforma" bson:"envio_por_email_plataforma" json:"envio_por_email_plataforma"`
	NitProveedorTecnologico string   `xml:"nitProveedorTecnologico" bson:"nit_proveedor_tecnologico" json:"nit_proveedor_tecnologico"`
}

type PagoType struct {
	Moneda               string    `xml:"moneda" bson:"moneda" json:"moneda"`
	TotalImporteBruto    float64   `xml:"totalimportebruto" bson:"total_importe_bruto" json:"total_importe_bruto"`
	TotalBaseImponible   float64   `xml:"totalbaseimponible" bson:"total_base_imponible" json:"total_base_imponible"`
	TotalBaseConImpuesto float64   `xml:"totalbaseconimpuesto" bson:"total_base_con_impuesto" json:"total_base_con_impuesto"`
	TotalFactura         float64   `xml:"totalfactura" bson:"total_factura" json:"total_factura"`
	TipoCompra           int       `xml:"tipocompra" bson:"tipo_compra" json:"tipo_compra"`
	PeriodoDePagoA       string    `xml:"periododepagoa" bson:"periodo_de_pago_a" json:"periodo_de_pago_a"`
	FechaVencimiento     time.Time `xml:"fechavencimiento" bson:"fecha_vencimiento" json:"fecha_vencimiento"`
}

type DetalleType struct {
	XMLName            xml.Name             `xml:"Detalle" bson:"-" json:"-"`
	CodigoProducto     string               `xml:"Codigoproducto" bson:"codigo_producto" json:"codigo_producto"`
	TipoCodigoProducto string               `xml:"tipocodigoproducto" bson:"tipo_codigo_producto" json:"tipo_codigo_producto"`
	NombreProducto     string               `xml:"nombreproducto" bson:"nombre_producto" json:"nombre_producto"`
	Descripcion        string               `xml:"descripcion" bson:"descripcion" json:"descripcion"`
	Cantidad           int                  `xml:"cantidad" bson:"cantidad" json:"cantidad"`
	ValorUnitario      float64              `xml:"valorunitario" bson:"valor_unitario" json:"valor_unitario"`
	PrecioSinImpuestos float64              `xml:"preciosinimpuestos" bson:"precio_sin_impuestos" json:"precio_sin_impuestos"`
	PrecioTotal        float64              `xml:"preciototal" bson:"precio_total" json:"precio_total"`
	TipoImpuesto       string               `xml:"tipoImpuesto" bson:"tipo_impuesto" json:"tipo_impuesto"`
	ListaImpuestos     []ImpuestoType       `xml:"listaImpuestos" bson:"lista_impuestos" json:"lista_impuestos"`
	ListaCargos        []CargoType          `xml:"listaCargos" bson:"lista_cargos" json:"lista_cargos"`
	ListaAdicionales   []CampoAdicionalType `xml:"listaAdicionales" bson:"Lista_adicionales" json:"Lista_adicionales"`
}

type CampoAdicionalType struct {
	XMLName     xml.Name `xml:"CampoAdicional" bson:"-" json:"-"`
	NombreCampo float64  `xml:"nombreCampo" bson:"nombre_campo" json:"nombre_campo"`
	ValorCampo  float64  `xml:"valorCampo" bson:"valor_campo" json:"valor_campo"`
}

type CargoType struct {
	XMLName         xml.Name `xml:"Cargo" bson:"-" json:"-"`
	ValorCargo      float64  `xml:"valorCargo" bson:"valor_cargo" json:"valor_cargo"`
	PorcentajeCargo float64  `xml:"porcentajeCargo" bson:"porcentaje_cargo" json:"porcentaje_cargo"`
	Descripcion     float64  `xml:"descripcion" bson:"descripcion" json:"descripcion"`
}

type ImpuestoType struct {
	XMLName                 xml.Name `xml:"Impuesto" bson:"-" json:"-"`
	CodigoImpuestoRetencion string   `xml:"codigoImpuestoRetencion" bson:"codigo_impuesto_retencion" json:"codigo_impuesto_retencion"`
	Porcentaje              float64  `xml:"porcentaje" bson:"porcentaje" json:"porcentaje"`
	ValorImpuestoRetencion  float64  `xml:"valorImpuestoRetencion" bson:"valor_impuesto_retencion" json:"valor_impuesto_retencion"`
	BaseImponible           float64  `xml:"baseimponible" bson:"base_imponible" json:"base_imponible"`
	IsAutoRetenido          bool     `xml:"isAutoRetenido" bson:"is_auto_retenido" json:"is_auto_retenido"`
}

type TerceroType struct {
	Id              string        `bson:"id" json:"id"`
	Tipo            string        `bson:"tipo" json:"tipo"`
	NumeroDocumento string        `bson:"numero_documento" json:"numero_documento"`
	Nombre          string        `bson:"nombre" json:"nombre"`
	Ubicacion       UbicacionType `bson:"ubicacion" json:"ubicacion"`
	Contacto        ContactoType  `bson:"contacto" json:"contacto"`
	Descripcion     string        `bson:"descripcion" json:"descripcion"`
	SiActivo        bool          `bson:"si_activo" json:"si_activo"`
}

type UbicacionType struct {
	Pais         ItemUbicacionType `bson:"pais" json:"pais"`
	Departamento ItemUbicacionType `bson:"departamento" json:"departamento"`
	Municipio    ItemUbicacionType `bson:"municipio" json:"municipio"`
	Direccion    string            `bson:"direccion" json:"direccion"`
}

type ItemUbicacionType struct {
	Codigo string `bson:"codigo" json:"codigo"`
	Nombre string `bson:"nombre" json:"nombre"`
}

type ContactoType struct {
	Telefonos []string `bson:"telefonos" json:"telefonos"`
	Correos   []string `bson:"correos" json:"correos"`
}

type ResolucionFacturacionType struct {
	Numero       string       `bson:"numero" json:"numero"`
	Fecha        time.Time    `bson:"fecha" json:"fecha"`
	Prefijo      string       `bson:"prefijo" json:"prefijo"`
	Rango        RangoType    `bson:"rango" json:"rango"`
	ClaveTecnica string       `bson:"clave_tecnica" json:"clave_tecnica"`
	Vigencia     VigenciaType `bson:"vigencia" json:"vigencia"`
	SiActivo     bool         `bson:"si_activo" json:"si_activo"`
}

type RangoType struct {
	Inferior int64 `bson:"inferior" json:"inferior"`
	Superior int64 `bson:"superior" json:"superior"`
}

type VigenciaType struct {
	Desde time.Time `bson:"desde" json:"desde"`
	Hasta time.Time `bson:"hasta" json:"hasta"`
}

type DominioType struct {
	Nombre      string `bson:"nombre" json:"nombre"`
	Valor       string `bson:"valor" json:"valor"`
	Descripcion string `bson:"descripcion" json:"descripcion"`
	SiActivo    bool   `bson:"si_activo" json:"si_activo"`
}

type SoftwareFacturacionType struct {
	Id       string `bson:"id" json:"id"`
	Pin      string `bson:"pin" json:"pin"`
	Nombre   string `bson:"nombre" json:"nombre"`
	Usuario  string `bson:"usuario" json:"usuario"`
	Password string `bson:"password" json:"password"`
	SiActivo bool   `bson:"si_activo" json:"si_activo"`
}
