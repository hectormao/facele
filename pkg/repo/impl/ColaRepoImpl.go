// ColaRepoImpl.go
package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	repo "github.com/hectormao/facele/pkg/repo"
	"github.com/streadway/amqp"
)

const (
	ColaEnvio        string = "enviarFacturaQ"
	ColaNotificacion string = "notificarFacturaQ"
)

type ColaRepoImpl struct {
	connection        *amqp.Connection
	channel           *amqp.Channel
	queueEnvio        amqp.Queue
	queueNotificacion amqp.Queue
}

func (cola ColaRepoImpl) AgregarEnColaEnvioFacturas(objectId string, facturaXML []byte) error {
	return cola.agregarEnCola(ColaEnvio, objectId, facturaXML)
}

func (cola ColaRepoImpl) GetFacturasAEnviar() (<-chan repo.FacturaCola, error) {
	return cola.getFacturas(ColaEnvio)
}

func (cola ColaRepoImpl) AgregarEnColaNotificacion(objectId string, documentoElectronico []byte) error {
	return cola.agregarEnCola(ColaNotificacion, objectId, documentoElectronico)
}

func (cola ColaRepoImpl) GetFacturasANotificar() (<-chan repo.FacturaCola, error) {
	return cola.getFacturas(ColaNotificacion)
}

func (cola ColaRepoImpl) agregarEnCola(nombreCola string, objectId string, xml []byte) error {
	log.Printf("Agregando a la cola %s %s id objeto: %s xml: %s", nombreCola, objectId, xml)
	if cola.channel == nil {
		err := cola.conectar()
		if err != nil {
			return errors.New(fmt.Sprintf("Error al conectar a la cola %v", err))
		}
	}

	body := repo.FacturaCola{
		Id:  objectId,
		Xml: xml,
	}

	log.Printf("body %v", body)

	jsonValue, err := json.Marshal(body)
	if err != nil {
		return errors.New(fmt.Sprintf("Error al parsear json %v", err))
	}

	log.Printf("Agregando a la cola %s de facturas %s", nombreCola, jsonValue)
	cola.channel.Publish(
		"",
		nombreCola,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonValue,
		},
	)

	return nil
}

func (cola ColaRepoImpl) getFacturas(nombreCola string) (<-chan repo.FacturaCola, error) {
	log.Printf("Creando consumidor de la cola %v", nombreCola)
	if cola.channel == nil {
		err := cola.conectar()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error al conectar a la cola %v", err))
		}
	}

	facturaChannel := make(chan repo.FacturaCola)

	msgs, err := cola.channel.Consume(
		nombreCola,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error al crear consumidor la cola %v", err))
	}

	go func() {
		for msg := range msgs {
			data := msg.Body
			var factura repo.FacturaCola
			jsonErr := json.Unmarshal(data, &factura)
			if jsonErr != nil {
				log.Printf("%v", jsonErr)
			} else {
				facturaChannel <- factura
			}
		}
	}()

	return facturaChannel, nil
}

func (cola ColaRepoImpl) terminar() error {
	err := cola.channel.Close()
	if err != nil {
		return err
	}
	err = cola.connection.Close()
	if err != nil {
		return err
	}

	return nil
}

func (cola *ColaRepoImpl) conectar() error {

	conn, err := amqp.Dial("amqp://guest:guest@172.17.0.2:5672/")
	if err != nil {
		return err
	}
	cola.connection = conn
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	cola.channel = ch

	queueEnvio, err := ch.QueueDeclare(
		ColaEnvio,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	cola.queueEnvio = queueEnvio

	log.Printf("Declarando Cola: %v", cola.queueEnvio)

	queueNotificacion, err := ch.QueueDeclare(
		ColaNotificacion,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	cola.queueNotificacion = queueNotificacion

	log.Printf("Declarando Cola: %v", cola.queueNotificacion)

	return nil
}
