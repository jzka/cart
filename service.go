package cart

import (
	"fmt"

	"github.com/cart/db"
	"github.com/go-kit/kit/log"
)

//Service is the interface for bussines methods
type Service interface {
	GetCart(userID string) (db.Cart, error)
	PostCart(c db.Cart) (string, error)
	DeleteCart(userID string) error
	DeleteProduct(prodID, userID string) error
	PostProduct(p db.Product, userID string) (string, error)
	//Health() []Health // GET /health
}

type cartService struct {
	db     *db.Mongo
	logger log.Logger
}

//NewCartService implements interface Service
func NewCartService(db *db.Mongo, logger log.Logger) Service {
	return &cartService{
		db:     db,
		logger: logger,
	}
}

func (s *cartService) GetCart(userID string) (db.Cart, error) {
	c, errc := s.db.GetCartForUser(userID)
	if errc != nil {
		return db.NewCart(), errc
	}
	errp := s.db.PopulateProductsForCart(&c)
	return c, errp
}

func (s *cartService) PostCart(c db.Cart) (string, error) {
	err := s.db.CreateCart(&c)
	return c.ID, err
}

func (s *cartService) DeleteCart(userID string) error {
	c, errc := s.db.GetCartForUser(userID)
	if errc != nil {
		return errc
	}
	return s.db.DeleteCart(c.ID)
}

func (s *cartService) DeleteProduct(prodID, userID string) error {
	fmt.Println("DeleteProduct service")
	c, errc := s.db.GetCartForUser(userID)
	if errc != nil {
		return errc
	}
	return s.db.DeleteProduct(c.ID, prodID)
}

func (s *cartService) PostProduct(p db.Product, userID string) (string, error) {
	c, errc := s.db.GetCartForUser(userID)
	if errc != nil {
		return "", errc
	}
	for _, exP := range c.Products {
		if exP.ID == p.ID {
			prod, errp := s.db.GetProduct(p.ID)
			if errp != nil {
				return "", errp
			}
			prod.Quantity += p.Quantity
			prod.UnitPrice = p.UnitPrice
			err := s.db.UpdateProduct(prod)
			return prod.ID, err
		}
	}
	err := s.db.CreateProduct(&p, c.ID)
	return p.ID, err
}

func (s *cartService) updateProduct(p db.Product) (string, error) {
	_, err := s.db.GetProduct(p.ID)
	if err != nil {
		return p.ID, err
	}
	erru := s.db.UpdateProduct(p)
	return p.ID, erru
}
