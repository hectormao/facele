package ent

import (
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
	XMLName       xml.Name          `xml:"listaFacturas" bson:"-" json:"-"`
	CabezaFactura CabezaFacturaType `xml:"cabezaFactura" bson:"cabezaFactura" json:"cabezaFactura"`
}

type CabezaFacturaType struct {
	Prefijo                string                `xml:"prefijo" bson:"prefijo" json:"prefijo"`
	Consecutivo            int                   `xml:"consecutivo" bson:"consecutivo" json:"consecutivo"`
	FechaFacturacion       time.Time             `xml:"fechafacturacion" bson:"fecha_facturacion" json:"fecha_facturacion"`
	TipoDocumento          int                   `xml:"tipodocumento" bson:"tipo_documento" json:"tipo_documento"`
	TipoIdentificacion     int                   `xml:"tipoidentificacion" bson:"tipo_identificacion" json:"tipo_identificacion"`
	TipoPersona            int                   `xml:"tipopersona" bson:"tipo_persona" json:"tipo_persona"`
	Telefono               string                `xml:"telefono" bson:"telefono" json:"telefono"`
	Pais                   string                `xml:"pais" bson:"pais" json:"pais"`
	Ciudad                 string                `xml:"ciudad" bson:"ciudad" json:"ciudad"`
	Departamento           string                `xml:"departamento" bson:"departamento" json:"departamento"`
	Direccion              string                `xml:"direccion" bson:"direccion" json:"direccion"`
	NombreVendedor         string                `xml:"nombrevendedor" bson:"nombre_vendedor" json:"nombre_vendedor"`
	NumeroIdentificacion   string                `xml:"numeroidentificacion" bson:"numero_identificacion" json:"numero_identificacion"`
	TipoCompra             int                   `xml:"tipocompra" bson:"tipo_compra" json:"tipo_compra"`
	FechaFacturaModificada time.Time             `xml:"fechafacturamodificada" bson:"fecha_factura_modificada" json:"fecha_factura_modificada"`
	FechaVencimiento       time.Time             `xml:"fechavencimiento" bson:"fecha_vencimiento" json:"fecha_vencimiento"`
	RazonSocial            string                `xml:"razonsocial" bson:"razon_social" json:"razon_social"`
	Observaciones          string                `xml:"observaciones" bson:"observaciones" json:"observaciones"`
	Moneda                 string                `xml:"moneda" bson:"moneda" json:"moneda"`
	TotalBaseImponible     float64               `xml:"totalbaseimponible" bson:"total_base_imponible" json:"total_base_imponible"`
	TotalFactura           float64               `xml:"totalfactura" bson:"total_factura" json:"total_factura"`
	TotalImporteBruto      float64               `xml:"totalimportebruto" bson:"total_importe_bruto" json:"total_importe_bruto"`
	Nit                    string                `xml:"nit" bson:"nit" json:"nit"`
	ListaDetalles          []DetalleType         `xml:"listaDetalles" bson:"lista_detalles" json:"lista_detalles"`
	ListaImpuestos         []ImpuestosCabezaType `xml:"listaImpuestos" bson:"lista_impuestos" json:"lista_impuestos"`
}

type DetalleType struct {
	XMLName            xml.Name              `xml:"Detalle" bson:"-" json:"-"`
	Cantidad           int                   `xml:"cantidad" bson:"cantidad" json:"cantidad"`
	CodigoProducto     int                   `xml:"codigoproducto" bson:"codigo_producto" json:"codigo_producto"`
	Descripcion        string                `xml:"descripcion" bson:"descripcion" json:"descripcion"`
	PrecioSinImpuestos float64               `xml:"preciosinimpuestos" bson:"precio_sin_impuestos" json:"precio_sin_impuestos"`
	PrecioTotal        float64               `xml:"preciototal" bson:"precio_total" json:"precio_total"`
	ValorUnitario      float64               `xml:"valorunitario" bson:"valor_unitario" json:"valor_unitario"`
	ListaImpuestos     []ImpuestoDetalleType `xml:"listaImpuestos" bson:"lista_impuestos" json:"lista_impuestos"`
}

type ImpuestoDetalleType struct {
	XMLName                 xml.Name `xml:"impuestoDetalle" bson:"-" json:"-"`
	CodigoProducto          string   `xml:"codigoproducto" bson:"codigo_producto" json:"codigo_producto"`
	BaseImponible           float64  `xml:"baseimponible" bson:"base_imponible" json:"base_imponible"`
	CodigoImpuestoRetencion string   `xml:"codigoImpuestoRetencion" bson:"codigo_impuesto_retencion" json:"codigo_impuesto_retencion"`
	Porcentaje              float64  `xml:"porcentaje" bson:"porcentaje" json:"porcentaje"`
	ValorImpuestoRetencion  float64  `xml:"valorImpuestoRetencion" bson:"valor_impuesto_retencion" json:"valor_impuesto_retencion"`
}

type ImpuestosCabezaType struct {
	XMLName                 xml.Name `xml:"impuestosCabeza" bson:"-" json:"-"`
	BaseImponible           float64  `xml:"baseimponible"`
	CodigoImpuestoRetencion string   `xml:"codigoImpuestoRetencion"`
	Porcentaje              float64  `xml:"porcentaje"`
	ValorImpuestoRetencion  float64  `xml:"valorImpuestoRetencion"`
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
	Pais         string `bson:"pais" json:"pais"`
	Departamento string `bson:"departamento" json:"departamento"`
	Municipio    string `bson:"municipio" json:"municipio"`
	Direccion    string `bson:"direccion" json:"direccion"`
}

type ContactoType struct {
	Telefono string `bson:"telefono" json:"telefono"`
	Correo1  string `bson:"correo1" json:"correo1"`
	Correo2  string `bson:"correo2" json:"correo2"`
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
