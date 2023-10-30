package categorizer

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// RowID : type used to represent the ID of a connection in the database
// Service : type used to represent the service active on the vulnerable machine
// Connection : type used to represent a connection, it contains all the informations related to the connection itself
// those are all necessary to interact with caronte, which defines the following types and organizes data in the database accordingly
type RowID primitive.ObjectID

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
