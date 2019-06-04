package impl

import (
	"github.com/hectormao/facele/pkg/repo/cfg"
	"github.com/hectormao/facele/pkg/trns"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type GenericodeRepoImpl struct {
	client     *mongo.Client
	db         *mongo.Database
	Config     cfg.MongoConfig
	Translator trns.GeneriCodeTrans
}

func (repo GenericodeRepoImpl) GetGenericodes(filename string) []ent.Genericode {
	xmlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	codes, err := repo.
		Translator.
		XMLToGenericode(xmlData)

	return codes, nil

}

func (repo GenericodeRepoImpl) saveGenericode(code ent.Genericode) (string, error) {
	log.Printf("Almacenando genericode %v", code)

	if repo.db == nil {
		err := repo.conectar()
		if err != nil {
			return "", err
		}
	}
	genericodeCl := repo.db.Collection("genericode")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	res, err := genericodeCl.InsertOne(ctx, code)
	if err != nil {
		return "", err
	}

	id := res.InsertedID.(objectid.ObjectID)

	if err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (repo GenericodeRepoImpl) GetGenericodesMap() (map[string]map[string]ent.Genericode, error) {
	log.Printf("Get Genericodes")

	if repo.db == nil {
		err := repo.conectar()
		if err != nil {
			return nil, err
		}
	}

	genericodeCl := repo.db.Collection("genericode")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cur, err := genericodeCl.Find()
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []ent.Genericode
	for cur.Next(ctx) {
		var genericode ent.Genericode
		err = cur.Decode(&genericode)
		if err != nil {
			return nil, err
		}
		result = append(result, genericode)
	}

	return repo.Translator.GenericodeListToMap(result)
}

func (repo *GenericodeRepoImpl) conectar() error {
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

func (repo *GenericodeRepoImpl) terminar() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	repo.client.Disconnect(ctx)
}
