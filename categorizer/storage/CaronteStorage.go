package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// RowID : type used to represent the ID of a connection in the database
// Service : type used to represent the service active on the vulnerable machine
// Connection : type used to represent a connection, it contains all the information related to the connection itself
// those are all necessary to interact with caronte, which defines the following types and organizes data in the database accordingly
type OrderedDocument = bson.D
type UnorderedDocument = bson.M
type Entry = bson.E
type RowID = primitive.ObjectID

type Service struct {
	Port  uint16 `json:"port" bson:"_id"`
	Name  string `json:"name" binding:"min=3" bson:"name"`
	Color string `json:"color" binding:"hexcolor" bson:"color"`
	Notes string `json:"notes" bson:"notes"`
}

type Connection struct {
	ID              RowID     `json:"id" bson:"_id"`
	SourceIP        string    `json:"ip_src" bson:"ip_src"`
	DestinationIP   string    `json:"ip_dst" bson:"ip_dst"`
	SourcePort      uint16    `json:"port_src" bson:"port_src"`
	DestinationPort uint16    `json:"port_dst" bson:"port_dst"`
	StartedAt       time.Time `json:"started_at" bson:"started_at"`
	ClosedAt        time.Time `json:"closed_at" bson:"closed_at"`
	ClientBytes     int       `json:"client_bytes" bson:"client_bytes"`
	ServerBytes     int       `json:"server_bytes" bson:"server_bytes"`
	ClientDocuments int       `json:"client_documents" bson:"client_documents"`
	ServerDocuments int       `json:"server_documents" bson:"server_documents"`
	ProcessedAt     time.Time `json:"processed_at" bson:"processed_at"`
	MatchedRules    []RowID   `json:"matched_rules" bson:"matched_rules"`
	Hidden          bool      `json:"hidden" bson:"hidden,omitempty"`
	Marked          bool      `json:"marked" bson:"marked,omitempty"`
	Comment         string    `json:"comment" bson:"comment,omitempty"`
	Service         Service   `json:"service" bson:"-"`
}

const (
	Connections       = "connections"
	ConnectionStreams = "connection_streams"
	ImportingSessions = "importing_sessions"
	Rules             = "rules"
	Searches          = "searches"
	Settings          = "settings"
	Services          = "services"
	Statistics        = "statistics"
)

// Storage : interface, defines the general methods to interact with the database
type Storage interface {
	Insert(collectionName string) InsertOperation
	Update(collectionName string) UpdateOperation
	Find(collectionName string) FindOperation
	Delete(collectionName string) DeleteOperation
}

// MongoStorage : implements Storage for MongoDB, it is used to interact with the database
type MongoStorage struct {
	client      *mongo.Client
	collections map[string]*mongo.Collection
}

// NewMongoStorage : creates a new MongoStorage instance
func NewMongoStorage(uri string, port int, database string) (*MongoStorage, error) {
	ctx := context.Background()
	opt := options.Client()
	opt.ApplyURI(fmt.Sprintf("mongodb://%s:%v", uri, port))
	client, err := mongo.NewClient(opt)
	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	db := client.Database(database)
	collections := map[string]*mongo.Collection{
		Connections:       db.Collection(Connections),
		ConnectionStreams: db.Collection(ConnectionStreams),
		ImportingSessions: db.Collection(ImportingSessions),
		Rules:             db.Collection(Rules),
		Searches:          db.Collection(Searches),
		Settings:          db.Collection(Settings),
		Services:          db.Collection(Services),
		Statistics:        db.Collection(Statistics),
	}

	if _, err := collections[Services].Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{"name", 1}},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}

	if _, err := collections[ConnectionStreams].Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"connection_id", -1}}, // descending
		},
		{
			Keys: bson.D{{"payload_string", "text"}},
		},
	}); err != nil {
		return nil, err
	}

	return &MongoStorage{
		client:      client,
		collections: collections,
	}, nil
}
