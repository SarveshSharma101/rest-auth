package DB

import (
	"context"
	"log"
	"os"
	"rest-auth/utils"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client *mongo.Client
	DBName string
}

var once sync.Once
var Db *DB = &DB{}

// ConnectToDB connects to the MongoDB database
// Uses sync.Once to ensure that the connection is established only
func ConnectToDB() error {
	var err error
	Db.DBName = os.Getenv(utils.DB_NAME)
	once.Do(func() {
		uri := os.Getenv(utils.DB_URI)
		log.Println("Connecting to DB: ", uri)
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err == nil {
			Db.Client = client
		}
	})
	return err
}
