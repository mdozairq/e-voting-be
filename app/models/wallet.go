package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Wallet represents the structure of the Wallets document.
type Wallet struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UID        string             `bson:"uid"`
	PublicKey  string             `bson:"public_key"`
	PrivateKey string             `bson:"private_key"`
	Balance    string             `bson:"balance"`
	UpdatedAt  string             `bson:"updated_at"`
	CreatedAt  time.Time          `bson:"created_at"`
}
