package DB

import (
	"context"
	"os"
	"rest-auth/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client *mongo.Client
	DBName string
}

func ConnectToDB() (*DB, error) {
	var db *DB = &DB{
		DBName: os.Getenv(utils.DB_NAME),
	}
	uri := os.Getenv(utils.DB_URI)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return db, err
	}
	db.Client = client
	return db, err
}
