package repository

import (
	"log"
	"os"

	"github.com/marcusolsson/goddd/location"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type locationRepositoryMongoDB struct {
}

func withCollection(f func(c *mgo.Collection)) {
	session, err := mgo.Dial(os.Getenv("MONGOHQ_URL"))

	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("app30695645").C("locations")

	f(c)
}

func (r *locationRepositoryMongoDB) Store(l location.Location) {
	withCollection(func(c *mgo.Collection) {
		if _, err := c.Upsert(bson.M{"unlocode": l.UNLocode}, l); err != nil {
			log.Fatal(err)
		}
	})
}

func (r *locationRepositoryMongoDB) Find(locode location.UNLocode) (location.Location, error) {
	return location.Location{}, nil
}

func (r *locationRepositoryMongoDB) FindAll() []location.Location {
	var result []location.Location

	withCollection(func(c *mgo.Collection) {
		if err := c.Find(bson.M{}).All(&result); err != nil {
			log.Fatal(err)
		}
	})

	return result
}

func ensureUNLocodeIndex() {
	withCollection(func(c *mgo.Collection) {
		index := mgo.Index{
			Key:    []string{"unlocode"},
			Unique: true,
		}

		err := c.EnsureIndex(index)

		if err != nil {
			panic(err)
		}
	})
}

func NewLocationRepositoryMongoDB() location.Repository {

	r := &locationRepositoryMongoDB{}

	ensureUNLocodeIndex()

	r.Store(location.Stockholm)
	r.Store(location.Hamburg)
	r.Store(location.Chicago)

	return r
}
