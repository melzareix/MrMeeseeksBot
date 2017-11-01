package Database

import (
	"os"
	"gopkg.in/mgo.v2"
	"github.com/melzareix/MrMeeseeksBot/Models"
	"gopkg.in/mgo.v2/bson"
)

var (
	DB *Mongo
)

const (
	USERS_COLLECTION = "users"
	DB_NAME = "MrMeseeks"
)

type Mongo struct {
	session *mgo.Session
}

// Connect to MongoDB Server
func (db *Mongo) Connect() error{
	MONGO_URL := os.Getenv("MONGO_URL")
	session, err := mgo.Dial(MONGO_URL)

	if err != nil {
		return err
	}
	db.session = session
	return nil
}

// Create A New User
func (db *Mongo) CreateUser(u *Models.User) (error){
	session := db.session.Copy()
	defer session.Close()

	c := session.DB(DB_NAME).C(USERS_COLLECTION)
	err := c.Insert(&u)

	return err
}

func (db *Mongo) GetUser(id string) (*Models.User, error) {
	session := db.session.Copy()
	defer session.Close()

	c := session.DB(DB_NAME).C(USERS_COLLECTION)

	user := &Models.User{}
	err := c.Find(bson.M{"uuid": id}).One(user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *Mongo) SaveUser(u *Models.User) error {
	session := db.session.Copy()
	defer session.Close()

	c := session.DB(DB_NAME).C(USERS_COLLECTION)

	err := c.Update(bson.M{"uuid": u.Uuid}, u)

	if err != nil {
		return err
	}
	return nil
}

// Bind the connection to global DB variable
func Init() {
	DB = &Mongo{}
	DB.Connect()
}