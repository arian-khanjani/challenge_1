package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type User struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Salt     string `json:"salt" bson:"salt"`
}

type Repo struct {
	client *mongo.Client
	db     *mongo.Database
	coll   *mongo.Collection
}

var repository *Repo

func main() {
	// initiate mongodb
	var err error
	repository, err = NewMongo(ConnProps{
		URI:  "",
		DB:   "",
		Coll: "",
	})
	if err != nil {
		log.Fatalln(err)
	}

	// initiate server
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/create-user", createUser)
	r.Get("/get-user/{email:^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$}", getUser)
	err = http.ListenAndServe(":3001", r)
	if err != nil {
		panic(err)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create user"))
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get user"))
}
