package impl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hectormao/facele/pkg/ent"
	repoCfg "github.com/hectormao/facele/pkg/repo/cfg"

	"github.com/hectormao/facele/pkg/cfg"
	"github.com/hectormao/facele/pkg/trns"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type FacturaRepoImpl struct {
	client *mongo.Client
	db     *mongo.Database
	Config repoCfg.MongoConfig
}

func (repo FacturaRepoImpl) AlmacenarFactura(factura ent.FacturaType) (string, error) {
	log.Printf("Almacenando factura %v", factura)

	if repo.db == nil {
		err := repo.conectar()
		if err != nil {
			return "", err
		}
	}
	facturaCl := repo.db.Collection("factura")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	res, err := facturaCl.InsertOne(ctx, factura)
	if err != nil {
		return "", err
	}

	id := res.InsertedID.(objectid.ObjectID)

	if err != nil {
		return "", err
	}

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

func (repo FacturaRepoImpl) GetEmpresa(numeroDocumento string) (*ent.EmpresaType, error) {

	log.Printf("Get Empresa %s", numeroDocumento)

	if repo.db == nil {
		err := repo.conectar()
		if err != nil {
			return nil, err
		}
	}

	empresaCl := repo.db.Collection("empresa")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.D{{"numero_documento", numeroDocumento}}
	cur, err := empresaCl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	if cur.Next(ctx) {
		var empresa ent.EmpresaType
		err = cur.Decode(&empresa)
		if err != nil {
			return nil, err
		}
		return &empresa, nil
	} else {
		return nil, errors.New(fmt.Sprintf("No se encontro Empresa: %s", numeroDocumento))
	}

}

func (repo FacturaRepoImpl) GetEmpresaPorId(id string) (*ent.EmpresaType, error) {

	log.Printf("Get Empresa por id: %s", id)

	if repo.db == nil {
		err := repo.conectar()
		if err != nil {
			return nil, err
		}
	}

	empresaCl := repo.db.Collection("empresa")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := objectid.FromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", objId}}
	cur, err := empresaCl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	if cur.Next(ctx) {
		var empresa ent.EmpresaType
		err = cur.Decode(&empresa)
		if err != nil {
			return nil, err
		}
		return &empresa, nil
	} else {
		return nil, errors.New(fmt.Sprintf("No se encontro Empresa con id: %s", id))
	}

}

func (repo *FacturaRepoImpl) conectar() error {
	ctx, _ := context.WithTimeout(context.Background(), repo.Config.getTimeout()*time.Second)
	client, err := mongo.Connect(ctx, repo.Config.getURL())
	if err != nil {
		return err
	}

	database := client.Database(repo.Config.getDatabase())

	repo.db = database
	repo.client = client

	return nil
}

func (repo *FacturaRepoImpl) terminar() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	repo.client.Disconnect(ctx)
}
