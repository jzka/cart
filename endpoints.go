package cart

import (
	"context"

	"github.com/cart/db"
	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects the endpoints that comprise the Service.
type Endpoints struct {
	GetCartEndpoint     endpoint.Endpoint
	PostCartEndpoint    endpoint.Endpoint
	DeleteCartEndpoint  endpoint.Endpoint
	PostProductEndpoint endpoint.Endpoint
	// UpdateProductEndPoint endpoint.Endpoint
	DeleteProductEndpoint endpoint.Endpoint
}

// MakeEndpoints returns an Endpoints structure, where each endpoint is
// backed by the given service.
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetCartEndpoint:     MakeGetCartEndpoint(s),
		PostCartEndpoint:    MakePostCartEndpoint(s),
		DeleteCartEndpoint:  MakeDeleteCartEndpoint(s),
		PostProductEndpoint: MakePostProductEndpoint(s),
		// UpdateProductEndPoint: MakeUpdateProductEndPoint(s),
		DeleteProductEndpoint: MakeDeleteProductEndpoint(s),
	}
}

// MakeGetCartEndpoint returns an endpoint via the given service.
func MakeGetCartEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetCartRequest)
		c, err := s.GetCart(req.ID)
		return getCartResponse{Cart: c}, err
	}
}

// MakePostCartEndpoint returns an endpoint via the given service.
func MakePostCartEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postCartRequest)
		id, err := s.PostCart(req.Cart)
		return postResponse{ID: id}, err
	}
}

// MakeDeleteCartEndpoint returns an endpoint via the given service.
func MakeDeleteCartEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteCartRequest)
		err := s.DeleteCart(req.UserID)
		if err == nil {
			return statusResponse{Status: true}, err
		}
		return statusResponse{Status: false}, err
	}
}

// MakePostProductEndpoint returns an endpoint via the given service.
func MakePostProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postProductRequest)
		id, err := s.PostProduct(req.Product, req.UserID)
		return postResponse{ID: id}, err
	}
}

// // MakeUpdateProductEndPoint returns an endpoint via the given service.
// func MakeUpdateProductEndPoint(s Service) endpoint.Endpoint {
// 	return func(ctx context.Context, request interface{}) (interface{}, error) {
// 		req := request.(updateProductRequest)
// 		id, err := s.UpdateProduct(req.Product)
// 		return postResponse{ID: id}, err
// 	}
// }

// MakeDeleteProductEndpoint returns an endpoint via the given service.
func MakeDeleteProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteProductRequest)
		err := s.DeleteProduct(req.ProductID, req.UserID)
		if err == nil {
			return statusResponse{Status: true}, err
		}
		return statusResponse{Status: false}, err
	}
}

type GetCartRequest struct {
	ID string
}

type getCartResponse struct {
	Cart db.Cart `json:"cart"`
}

type postCartRequest struct {
	db.Cart
}

type postResponse struct {
	ID string `json:"id"`
}

type deleteCartRequest struct {
	UserID string
}

type postProductRequest struct {
	db.Product
	UserID string
}

type updateProductRequest struct {
	db.Product
}

type deleteProductRequest struct {
	ProductID string
	UserID    string
}

type statusResponse struct {
	Status bool `json:"status"`
}

type healthRequest struct {
	//
}

type EmbedStruct struct {
	Embed interface{} `json:"_embedded"`
}
