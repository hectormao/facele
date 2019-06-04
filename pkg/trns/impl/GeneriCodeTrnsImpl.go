package impl

import (
	"log"

	"github.com/beevik/etree"
	"github.com/hectormao/facele/pkg/ent"
)

type GeneriCodeTrnsImpl struct {
}

func (trns GeneriCodeTrnsImpl) XMLToGenericode(xml []byte) ([]ent.Genericode, error) {
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(xml)
	if err != nil {
		return nil, err
	}
	root := doc.SelectElement("gc:CodeList")
	name := trns.getNombreColleccion(root)
	log.Printf("Nombre Colleccion: %s", name)

	codigos := trns.getCodigos(name, root)
	log.Printf("Codigos: %v", codigos)
	return codigos, nil
}

func (GeneriCodeTrnsImpl) GenericodeListToMap(genericodes []ent.Genericode) map[string]map[string]ent.Genericode {
	result := make(map[string]map[string]ent.Genericode)
	for _, genericode := range genericodes {
		generi, ok := result[genericode.Genericode]
		if !ok {
			codeMap := make(map[string]ent.Genericode)
			result[genericode.Genericode] = codeMap
			generi = codeMap
		}
		generi[genericode.Code] = genericode
	}
	return result
}

func (GeneriCodeTrnsImpl) getNombreColleccion(root *etree.Element) string {
	identification := root.SelectElement("Identification")
	name := identification.
		SelectElement("ShortName").
		Text()
	return name
}

func (GeneriCodeTrnsImpl) getCodigos(name string, root *etree.Element) []ent.Genericode {
	codeList := root.SelectElement("SimpleCodeList")
	var result []map[string]string
	for _, row := range codeList.SelectElements("Row") {
		code := ent.Genericode{}
		code.Genericode = name
		for _, value := range row.SelectElements("Value") {
			columnRef := value.SelectAttr("ColumnRef").Value
			value := value.SelectElement("SimpleValue").Text()
			if columnRef == "code" {
				code.Code = value
			} else if columnRef == "name" {
				code.Name = value
			}
		}
		result = append(result, code)
	}
	return result
}
