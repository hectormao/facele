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
