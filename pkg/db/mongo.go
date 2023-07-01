/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package db

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient struct {
	Client                  *mongo.Client
	Context                 context.Context
	ClustersCollection      *mongo.Collection
	GitopsCatalogCollection *mongo.Collection
	ServicesCollection      *mongo.Collection
}

type MongoDBClientParameters struct {
	HostType string
	Host     string
	Username string
	Password string
}

// Connect
func Connect(p *MongoDBClientParameters) *MongoDBClient {
	var connString string
	var clientOptions *options.ClientOptions

	ctx := context.Background()

	switch p.HostType {
	case "atlas":
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		connString = fmt.Sprintf("mongodb+srv://%s:%s@%s",
			p.Username,
			p.Password,
			p.Host,
		)
		clientOptions = options.Client().ApplyURI(connString).SetServerAPIOptions(serverAPI)
	case "local":
		connString = fmt.Sprintf("mongodb://%s:%s@%s/?authSource=admin",
			p.Username,
			p.Password,
			p.Host,
		)
		clientOptions = options.Client().ApplyURI(connString)
	default:
		log.Fatalf("could not create mongodb client: %s is not a valid host type option", p.HostType)
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("could not create mongodb client: %s", err)
	}

	cl := MongoDBClient{
		Client:                  client,
		Context:                 ctx,
		ClustersCollection:      client.Database("api").Collection("clusters"),
		GitopsCatalogCollection: client.Database("api").Collection("gitops-catalog"),
		ServicesCollection:      client.Database("api").Collection("services"),
	}

	return &cl
}

// TestDatabaseConnection
func (mdbcl *MongoDBClient) TestDatabaseConnection() error {
	err := mdbcl.Client.Database("admin").RunCommand(mdbcl.Context, bson.D{{Key: "ping", Value: 1}}).Err()
	if err != nil {
		log.Fatalf("error connecting to mongodb: %s", err)
	}

	return nil
}
