package model

type User struct {
	ID       string `bson:"id"`
	FullName string `bson:"full_name"`
	Address  string `bson:"address"`
}

type Users []User
