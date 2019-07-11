// ColaRepoImpl.go
package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hectormao/facele/pkg/ent"
	repoCfg "github.com/hectormao/facele/pkg/repo/cfg"
	"github.com/streadway/amqp"
)

const (
	ColaEnvio        string = "enviarFacturaQ"
	ColaNotificacion string = "notificarFacturaQ"
)

type ColaRepoImpl struct {
	Config            repoCfg.RabbitConfig
	connection        *amqp.Connection
	channel           *amqp.Channel
	queueEnvio        amqp.Queue
	queueNotificacion amqp.Queue
}

func (cola ColaRepoImpl) AgregarEnColaEnvioFacturas(factura ent.FacturaType) error {
	return cola.agregarEnCola(ColaEnvio, factura)
}

func (cola ColaRepoImpl) GetFacturasAEnviar() (<-chan ent.FacturaType, error) {
	return cola.getFacturas(ColaEnvio)
}

func (cola ColaRepoImpl) AgregarEnColaNotificacion(factura ent.FacturaType) error {
	return cola.agregarEnCola(ColaNotificacion, factura)
}

func (cola ColaRepoImpl) GetFacturasANotificar() (<-chan ent.FacturaType, error) {
	return cola.getFacturas(ColaNotificacion)
}

func (cola ColaRepoImpl) agregarEnCola(nombreCola string, factura ent.FacturaType) error {
	log.Printf("Agregando a la cola %s factura: %s", nombreCola, factura)

	if cola.channel == nil {
		err := cola.conectar()
		if err != nil {
			return errors.New(fmt.Sprintf("Error al conectar a la cola %v", err))
		}
	}

	facturaJson, err := json.Marshal(factura)
	if err != nil {
		return errors.New(fmt.Sprintf("Error al parsear la factura %v", err))
	}

	cola.channel.Publish(
		"",
		nombreCola,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        facturaJson,
		},
	)

	return nil
}

func (cola ColaRepoImpl) getFacturas(nombreCola string) (<-chan ent.FacturaType, error) {
	log.Printf("Creando consumidor de la cola %v", nombreCola)
	if cola.channel == nil {
		err := cola.conectar()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error al conectar a la cola %v", err))
		}
	}

	facturaChannel := make(chan ent.FacturaType)

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
			var factura ent.FacturaType
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
	urlRabbit := cola.Config.GetUrl()
	conn, err := amqp.Dial(urlRabbit)
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
		cola.Config.
			GetEnvioDianQueue().
			GetName(),
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
		cola.Config.
			GetNotificacionQueue().
			GetName(),
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
