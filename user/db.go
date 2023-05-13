package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ConnProps struct {
	URI  string
	DB   string
	Coll string
}

func NewMongo(props ConnProps) (*Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	clientOptions := options.Client().ApplyURI(props.URI)
	c, err := mongo.NewClient(clientOptions)
	err = c.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = c.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := c.Database(props.DB)
	coll := db.Collection(props.Coll)

	_, err = coll.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)

	return &Repo{
		client: c,
		db:     db,
		coll:   coll,
	}, nil
}

func (r *Repo) Get(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.coll.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) Create(ctx context.Context, user *User) (*User, error) {
	_, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repo) Disconnect(ctx context.Context) error {
	err := r.client.Disconnect(ctx)
	if err != nil {
		return err
	}

	return nil
}
