package cosmosdbclient

import (
	"context"
	"errors"
	"time"

	"github.com/chrobles/go-rest-api/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client : client for interacting with cosmosdb
type Client struct {
	Connstr string
	Opts    *options.ClientOptions
}

// Configure : apply client configuration
func (client *Client) Configure(cfg types.Config) error {
	client.Connstr = cfg.CosmosDB.Connstr
	if client.Connstr == "" {
		return errors.New("cosmsosdb requires CDB_CONN_STR")
	}

	client.Opts = options.Client().ApplyURI(client.Connstr)
	return nil
}

// Index : index market listing data in cosmosdb
func (client *Client) Index(listings *types.MarketListings) (interface{}, error) {
	var (
		cancel     context.CancelFunc
		collection *mongo.Collection
		ctx        context.Context
		err        error
		id         interface{}
		mdbclient  *mongo.Client
		res        *mongo.InsertOneResult
	)

	mdbclient, err = mongo.Connect(context.TODO(), client.Opts)
	if err != nil {
		return nil, err
	}

	ctx, _ = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection = mdbclient.Database("base").Collection("container")

	res, err = collection.InsertOne(ctx, *listings)
	if err != nil {
		return nil, err
	}
	id = res.InsertedID

	return id, nil
}
