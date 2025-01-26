package DB

import (
	"context"
	"os"
	"rest-auth/datamodel"
	"rest-auth/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func (db *DB) InsertSession(session *datamodel.Session) error {
	_, err := db.Client.Database(db.DBName).Collection(os.Getenv(utils.SESSIONS_COLLECTION)).InsertOne(context.TODO(), session)
	return err
}

func (db *DB) DeleteSession(email string) error {
	_, err := db.Client.Database(db.DBName).Collection(os.Getenv(utils.SESSIONS_COLLECTION)).DeleteOne(context.TODO(), bson.D{{Key: "email", Value: email}})
	return err
}

func (db *DB) GetSession(sessionId string) (datamodel.Session, error) {
	session := datamodel.Session{}
	err := db.Client.Database(db.DBName).Collection(os.Getenv(utils.SESSIONS_COLLECTION)).FindOne(context.TODO(), bson.D{{Key: "sessionid", Value: sessionId}}).Decode(&session)
	return session, err
}
