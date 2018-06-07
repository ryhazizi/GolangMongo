package repository

import (
	"golangmongo/src/modules/user/model"
	mgo "gopkg.in/mgo.v2"
)

type userRepositoryMongo struct {
	db         *mgo.Database
	collection string
}

func NewUserRepositoryMongo(db *mgo.Database, collection string) *userRepositoryMongo {
	return &userRepositoryMongo{
		db:         db,
		collection: collection,
	}
}

func (r *userRepositoryMongo) Insert(user *model.User) error {
	err := r.db.C(r.collection).Insert(user)
	return err
}
