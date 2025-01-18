package DB

import (
	"context"
	"os"
	"rest-auth/datamodel"
	"rest-auth/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// GetUsers gets all the users from the database
func (db *DB) GetUsers(emailId string) ([]datamodel.User, error) {
	var users []datamodel.User
	filter := bson.D{}
	if emailId != "" {
		filter = bson.D{{Key: "email", Value: emailId}}
	}
	cur, err := db.Client.Database(db.DBName).
		Collection(os.Getenv(utils.USERS_COLLECTION)).
		Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	err = cur.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (db *DB) GetUserByEmail(email string) (*datamodel.User, error) {
	var user *datamodel.User
	err := db.Client.Database(db.DBName).
		Collection(os.Getenv(utils.USERS_COLLECTION)).
		FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).
		Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) InsertUser(user *datamodel.User) error {
	_, err := db.Client.Database(db.DBName).Collection(os.Getenv(utils.USERS_COLLECTION)).InsertOne(context.TODO(), user)
	return err
}

func (db *DB) DeleteUser(email string) (int, error) {
	count, err := db.Client.Database(db.DBName).Collection(os.Getenv(utils.USERS_COLLECTION)).DeleteOne(context.TODO(), bson.D{{Key: "email", Value: email}})
	return int(count.DeletedCount), err
}

func (db *DB) UpdateUser(emailId string, user *datamodel.User) (int64, error) {
	count, err := db.Client.Database(db.DBName).Collection(os.Getenv(utils.USERS_COLLECTION)).UpdateOne(context.TODO(), bson.D{{Key: "email", Value: emailId}}, user)
	return count.ModifiedCount, err
}

func (db *DB) PatchUser(email, key string, value interface{}) (int64, error) {
	count, err := db.Client.Database(db.DBName).Collection(os.Getenv(utils.USERS_COLLECTION)).UpdateOne(context.TODO(), bson.D{{Key: "email", Value: email}}, bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: value}}}})
	return count.ModifiedCount, err

}
