package Database

import (
	"os"
	"gopkg.in/mgo.v2"
	"log"
	"github.com/melzareix/MrMeseeksBot/Models"
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
func (db *Mongo) Connect(){
	MONGO_URL := os.Getenv("MONGO_URL")
	session, err := mgo.Dial(MONGO_URL)

	if err != nil {
		log.Fatal(err)
	}
	db.session = session

}

// Create A New User
func (db *Mongo) CreateUser(u *Models.User) (error){
	session := db.session.Copy()
	defer session.Close()

	c := session.DB(DB_NAME).C(USERS_COLLECTION)
	err := c.Insert(&u)

	return err
}


// Bind the connection to global DB variable
func Init() {
	DB = &Mongo{}
	DB.Connect()
}