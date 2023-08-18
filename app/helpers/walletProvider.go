package helpers

import (
	"context"
	"log"
	"time"

	"github.com/mdozairq/e-voting-be/app/models"
	"github.com/mdozairq/e-voting-be/database"
	"go.mongodb.org/mongo-driver/bson"
)

var walletCollection = database.OpenCollection(database.Client, "wallet")

func GetUnusedPrivateKey(uid string) (Address string, PrivateKey string)  {
	var walletProvider models.Wallet
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := walletCollection.CountDocuments(ctx, bson.M{"uid": uid})

	if err != nil {
		log.Panic(err)
		return 
	}

	if count > 0 {
		log.Panic(err)
		return 
	}

	update := bson.D{
		{"$set", bson.D{
			{"uid", uid},
			{"updated_at", time.Now().UTC},
		}},
	}
	_, err = walletCollection.UpdateOne(ctx, bson.M{"uid": bson.TypeNull}, update)
	if err != nil {
		log.Panic(err)
		return 
	}

	walletCollection.FindOne(ctx, bson.M{"uid": uid}).Decode(&walletProvider)

	return walletProvider.PublicKey, walletProvider.PrivateKey
}
