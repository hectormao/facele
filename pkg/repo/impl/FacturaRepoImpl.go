package impl

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/clbanning/mxj"
)

type FacturaRepoImpl struct {
	client *mongo.Client
	db     *mongo.Database
}

func (repo FacturaRepoImpl) AlmacenarFactura(cliente string, nombreArchivo string, contenido []byte) (string, error) {
	log.Printf("Almacenando Archivo - cliente: %s nombre: %s contenido: %s", cliente, nombreArchivo, contenido)

	if repo.db == nil {
		err := repo.conectar()
		if err != nil {
			return "", err
		}
	}
	facturaCl := repo.db.Collection("factura")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	mapValue, err := mxj.NewMapXml(contenido)

	if err != nil {
		return "", err
	}

	mapValue["_cliente"] = cliente

	res, err := facturaCl.InsertOne(ctx, mapValue)
	if err != nil {
		return "", err
	}

	id := res.InsertedID.(objectid.ObjectID)

	return id.Hex(), nil
}

func (repo FacturaRepoImpl) ActualizarFactura(id string, data map[string]interface{}) error {
	log.Printf("Actualizando Factura id: %v data: %v", id, data)

	if repo.db == nil {
		err := repo.conectar()
		if err != nil {
			return err
		}
	}
	facturaCl := repo.db.Collection("factura")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := objectid.FromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", objId}}
	update := bson.D{
		{
			"$set",
			data,
		},
	}

	_, err = facturaCl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *FacturaRepoImpl) conectar() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, "mongodb://172.17.0.3:27017")
	if err != nil {
		return err
	}

	database := client.Database("facele")

	repo.db = database
	repo.client = client

	return nil
}

func (repo *FacturaRepoImpl) terminar() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	repo.client.Disconnect(ctx)
}
