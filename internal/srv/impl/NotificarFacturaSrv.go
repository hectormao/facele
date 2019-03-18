package impl

import (
	"bytes"
	"html/template"
	"log"

	"errors"
	"fmt"

	"github.com/hectormao/facele/pkg/repo"
	"github.com/scorredoira/email"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"

	"net/mail"
	"net/smtp"
)

type NotificarFacturaSrvImpl struct {
	Repo                 repo.FacturaRepo
	ColaNotificacionRepo repo.ColaNotificacionRepo
}

func (srv NotificarFacturaSrvImpl) IniciarConsumidorCola() error {
	facturasANotificar, err := srv.ColaNotificacionRepo.GetFacturasANotificar()
	if err != nil {
		return errors.New(fmt.Sprintf("Error al crear consumidor: %v", err))
	}
	for facturaANotificar := range facturasANotificar {
		id := facturaANotificar["_id"].(string)
		log.Printf("Recibiendo factura: %v", id)
		pdfBytes, err := srv.crearPdf(id)
		if err != nil {
			log.Printf("Error al crear PDF: %v", err)
			continue
		}

		err = srv.enviarCorreo("ing.ogonzalez@gmail.com", pdfBytes)
		if err != nil {
			log.Printf("Error al enviar correo: %v", err)
			continue
		}
	}
	return nil
}

func (srv NotificarFacturaSrvImpl) crearPdf(id string) ([]byte, error) {
	html := `
	<html>
		<p>Esta es la factura: <strong>{{.Id}}</strong></p><img src="https://www.popsci.com/sites/popsci.com/files/styles/1000_1x_/public/images/2018/08/nasa-logo-web-rgb.png?itok=Iczdwajo&amp;fc=50,50" alt="" width="100" height="50" />
		<table>
		  <tr>
		  	<th>Item</th>
		    <th>Cantidad</th>
		    <th>Valor Unitario</th>
		    <th>Valor</th>
		  </tr>
		  <tr>
		  	<td>Mi Articulo</td>
		    <td>10</td>
		    <td>10</td>
		    <td>100</td>
		  </tr>
		</table>
	<html>`
	data := struct {
		Id string
	}{
		Id: id,
	}

	t, err := template.New("factura").Parse(html)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	err = t.Execute(&buffer, data)
	if err != nil {
		return nil, err
	}

	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pdfg.AddPage(wkhtml.NewPageReader(&buffer))
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	pdfBytes := pdfg.Bytes()
	return pdfBytes, nil
}

func (srv NotificarFacturaSrvImpl) enviarCorreo(correo string, pdfBytes []byte) error {
	log.Printf("Enviando Factura por correo: %s pdf: %d", correo, len(pdfBytes))

	from := "hectormao.gonzalez@gmail.com"
	password := "lcdgp2012$"

	m := email.NewHTMLMessage("Factura FecEle", "<html><h1>Este es el mensaje del correo</h1></html>")
	m.From = mail.Address{Name: "From", Address: from}
	m.To = []string{correo}
	err := m.AttachBuffer("factura.pdf", pdfBytes, false)

	if err != nil {
		return err
	}

	err = email.Send("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		m)

	if err != nil {
		return err
	}

	return nil
}
