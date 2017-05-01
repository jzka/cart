package db

import "gopkg.in/mgo.v2/bson"

// Product represents a product in a cart
type Product struct {
	ID        string  `json:"productID" bson:"-"`
	Quantity  int     `json:"quantity" bson:"quantity"`
	UnitPrice float64 `json:"unitPrice" bson:"unitPrice"`
}

//ProductDB representes db product
type ProductDB struct {
	Product `bson:",inline"`
	ID      bson.ObjectId `bson:"_id"`
}
