package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"regexp"
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
		URI:  "mongodb://dATeRMinCOgM:1pO2xHFkyR9S@localhost:27017",
		DB:   "challenge_1",
		Coll: "users",
	})
	if err != nil {
		log.Fatalln(err)
	}

	// initiate server
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/create-user", createUser)
	r.Get("/get-user/{email}", getUser)
	err = http.ListenAndServe(":3001", r)
	if err != nil {
		panic(err)
	}
}

const emailRegex = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"

func createUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match, err := regexp.MatchString(emailRegex, user.Email)
	if !match {
		http.Error(w, errors.New("email format is not valid").Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:3000/generate-salt", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		http.Error(w, "status was not 200", http.StatusInternalServerError)
		return
	}

	var salt string
	err = json.NewDecoder(res.Body).Decode(&salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	oldPassword := user.Password
	user.Password = hashPassword(oldPassword, salt)

	if match := doPasswordsMatch(user.Password, oldPassword, salt); !match {
		http.Error(w, errors.New("password hash and salt encountered an error").Error(), http.StatusInternalServerError)
		return
	}

	user.Salt = salt
	create, err := repository.Create(ctx, &user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			http.Error(w, errors.New("duplicate email").Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, create)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := chi.URLParam(r, "email")
	if email == "" {
		http.Error(w, "email not valid", http.StatusBadRequest)
		return
	}

	user, err := repository.Get(ctx, email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func hashPassword(password string, salt string) string {

	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Create sha-512 hasher
	md5Hasher := md5.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	// Write password bytes to the hasher
	md5Hasher.Write(passwordBytes)

	// Get the md5 hashed password
	hashedPasswordBytes := md5Hasher.Sum(nil)

	// Convert the hashed password to a hex string
	hashedPasswordHex := hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}

func doPasswordsMatch(hashedPassword, currPassword, salt string) bool {
	return hashedPassword == hashPassword(currPassword, salt)
}
