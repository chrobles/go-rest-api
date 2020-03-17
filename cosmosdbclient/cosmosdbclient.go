package cosmosdbclient

import (
	"context"
	"errors"
	"time"

	"github.com/chrobles/go-rest-api/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client : client for interacting with cosmosdb
type Client struct {
	Connstr string
	Client  *mongo.Client
}

// Record : data to be indexed in cosmosdb
type Record struct {
	ID       uuid.UUID
	Listings types.MarketListings
}

// Configure : apply client configuration
func (client *Client) Configure(cfg types.Config) error {
	var (
		err error
	)

	client.Connstr = cfg.CosmosDB.Connstr
	if client.Connstr == "" {
		return errors.New("cosmsosdb requires CDB_CONN_STR")
	}

	client.Client, err = mongo.NewClient(options.Client().ApplyURI(client.Connstr))
	if err != nil {
		return err
	}

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
		uid        uuid.UUID
		record     *Record
		res        *mongo.InsertOneResult
	)

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	uid, err = uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	record = new(Record)
	record.ID = uid
	record.Listings = *listings

	collection = client.Client.Database("base").Collection("container")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err = collection.InsertOne(ctx, record)
	if err != nil {
		return nil, err
	}
	spew.Dump(res)

	return id, nil
}
