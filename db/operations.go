package db

import (
	"errors"
	"flag"
	"time"

	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	mname     string
	mpassword string
	mhost     string
	db        = "carts"
	//ErrInvalidHexID represents a entity id that is not a valid bson ObjectID
	ErrInvalidHexID = errors.New("Invalid Id Hex")
)

func init() { // os.Getenv("MONGO_USER")
	flag.StringVar(&mname, "mongo-user", "", "Mongo user")
	flag.StringVar(&mpassword, "mongo-password", "", "Mongo password")
	flag.StringVar(&mhost, "mongo-host", "127.0.0.1:27017", "Mongo host")
}

type Mongo struct {
	Session *mgo.Session
}

// Init MongoDB
func (m *Mongo) Init() error {
	//u := getURL()
	var err error
	m.Session, err = mgo.DialWithTimeout("127.0.0.1", time.Duration(5)*time.Second)
	return err
}

//Ping checks db connection
func (m *Mongo) Ping() error {
	s := m.Session.Copy()
	defer s.Close()
	return s.Ping()
}

// CRUD Op

//CreateCart creates or updates a cart
func (m *Mongo) CreateCart(c *Cart) error {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	dbc := NewCartDB()
	dbc.ID = id
	dbc.Cart = *c
	cdb := s.DB("").C("carts")

	_, err := cdb.UpsertId(dbc.ID, dbc)
	if err != nil {
		return err
	}
	dbc.Cart.ID = dbc.ID.Hex()
	*c = dbc.Cart
	return nil
}

//GetCartForUser returns cart for user id
func (m *Mongo) GetCartForUser(userID string) (Cart, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("carts")
	dbc := NewCartDB()
	fmt.Println(userID)
	err := c.Find(bson.M{"userID": userID}).One(&dbc)
	dbc.ConvertObjectsIds()
	return dbc.Cart, err
}

//GetCart returns a cart with given id
func (m *Mongo) GetCart(id string) (Cart, error) {
	s := m.Session.Copy()
	defer s.Close()
	if !bson.IsObjectIdHex(id) {
		return NewCart(), ErrInvalidHexID
	}
	c := s.DB("").C("carts")
	dbc := NewCartDB()
	err := c.FindId(bson.ObjectIdHex(id)).One(&dbc)
	dbc.ConvertObjectsIds()
	return dbc.Cart, err
}

//DeleteCart deletes the cart
func (m *Mongo) DeleteCart(id string) error {
	if !bson.IsObjectIdHex(id) {
		return ErrInvalidHexID
	}
	fmt.Println("delete cart")
	s := m.Session.Copy()
	defer s.Close()
	cdb := s.DB("").C("carts")
	c, err := m.GetCart(id)
	if err != nil {
		return err
	}
	prodsIds := make([]bson.ObjectId, 0)
	for _, p := range c.Products {
		prodsIds = append(prodsIds, bson.ObjectIdHex(p.ID))
	}
	pdb := s.DB("").C("products")
	pdb.RemoveAll(bson.M{"_id": bson.M{"$in": prodsIds}})
	errc := cdb.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return errc
}

//PopulateProductsForCart adds products to a cart struct
func (m *Mongo) PopulateProductsForCart(c *Cart) error {
	s := m.Session.Copy()
	defer s.Close()
	ids := make([]bson.ObjectId, 0)
	for _, p := range c.Products {
		if !bson.IsObjectIdHex(p.ID) {
			return ErrInvalidHexID
		}
		ids = append(ids, bson.ObjectIdHex(p.ID))
	}

	var prodDB []ProductDB
	pdb := s.DB("").C("products")
	err := pdb.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&prodDB)
	if err != nil {
		return err
	}
	pra := make([]Product, 0)
	for _, prd := range prodDB {
		prd.Product.ID = prd.ID.Hex()
		pra = append(pra, prd.Product)
	}
	c.Products = pra
	return nil
}

//CreateProduct creates or updates a product
func (m *Mongo) CreateProduct(p *Product, cartID string) error {
	if cartID == "" || !bson.IsObjectIdHex(cartID) {
		return ErrInvalidHexID
	}
	if p.ID == "" || !bson.IsObjectIdHex(p.ID) {
		return ErrInvalidHexID
	}

	s := m.Session.Copy()
	defer s.Close()
	pdb := s.DB("").C("products")
	id := bson.ObjectIdHex(p.ID)
	prodDB := ProductDB{Product: *p, ID: id}
	_, err := pdb.UpsertId(prodDB.ID, prodDB)
	if err != nil {
		return err
	}
	if cartID != "" && p.ID != "" {
		err = m.addProductToCart(id, cartID)
		if err != nil {
			return err
		}
	}
	prodDB.Product.ID = prodDB.ID.Hex()
	*p = prodDB.Product

	return err
}

//GetProduct returns a product with given id
func (m *Mongo) GetProduct(id string) (Product, error) {
	if id == "" || !bson.IsObjectIdHex(id) {
		return Product{}, ErrInvalidHexID
	}
	s := m.Session.Copy()
	defer s.Close()
	pdb := s.DB("").C("products")
	prodDB := ProductDB{}
	err := pdb.FindId(bson.ObjectIdHex(id)).One(&prodDB)
	prodDB.Product.ID = prodDB.ID.Hex()
	return prodDB.Product, err
}

//DeleteProduct removes a product from db and the current cart
func (m *Mongo) DeleteProduct(cartID, prodID string) error {
	if cartID == "" || !bson.IsObjectIdHex(cartID) {
		return ErrInvalidHexID
	}
	if prodID == "" || !bson.IsObjectIdHex(prodID) {
		return ErrInvalidHexID
	}
	fmt.Println("delete prod")
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("products")
	err := c.Remove(bson.M{"_id": bson.ObjectIdHex(prodID)})
	if err != nil {
		return err
	}
	errc := m.removeProductFromCart(bson.ObjectIdHex(prodID), cartID)
	return errc
}

//UpdateProduct updates a product with new quantit and unit price
func (m *Mongo) UpdateProduct(p Product) error {
	if p.ID == "" || !bson.IsObjectIdHex(p.ID) {
		return ErrInvalidHexID
	}
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("products")
	err := c.Update(bson.M{"_id": bson.ObjectIdHex(p.ID)},
		bson.M{"$set": bson.M{
			"quantity":  p.Quantity,
			"unitPrice": p.UnitPrice,
		}})
	if err != nil {
		return err
	}
	return nil
}

func (m *Mongo) addProductToCart(id bson.ObjectId, cartID string) error {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("carts")
	return c.Update(bson.M{"_id": bson.ObjectIdHex(cartID)},
		bson.M{"$addToSet": bson.M{"products": id}})
}

func (m *Mongo) removeProductFromCart(id bson.ObjectId, cartID string) error {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("carts")
	return c.Update(bson.M{"_id": bson.ObjectIdHex(cartID)},
		bson.M{"$pull": bson.M{"products": id}})
}
