package db

import (
	"gopkg.in/mgo.v2/bson"
)

// Cart describes a cart
type Cart struct {
	ID       string    `json:"cartID" bson:"-"`
	UserID   string    `json:"userID" bson:"userID"`
	Products []Product `json:"-,omitempty" bson:"-"`
}

// CartDB blabla
type CartDB struct {
	Cart       `bson:",inline"`
	ID         bson.ObjectId   `bson:"_id"`
	ProductIDs []bson.ObjectId `bson:"products"`
}

//NewCart return a cart
func NewCart() Cart {
	c := Cart{Products: make([]Product, 0)}
	return c
}

//NewCartDB returns cartDB
func NewCartDB() CartDB {
	c := NewCart()
	return CartDB{
		Cart:       c,
		ProductIDs: make([]bson.ObjectId, 0),
	}
}

// ConvertObjectsIds converts mongo ids to hex
func (dbc *CartDB) ConvertObjectsIds() {
	if dbc.Cart.Products == nil {
		dbc.Cart.Products = make([]Product, 0)
	}
	for _, id := range dbc.ProductIDs {
		dbc.Cart.Products = append(dbc.Cart.Products, Product{ID: id.Hex()})
	}
	dbc.Cart.ID = dbc.ID.Hex()
}
