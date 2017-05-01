package cart

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"context"
)

var (
	ErrInvalidRequest = errors.New("Invalid request")
)

// MakeHTTPHandler mounts the endpoints into a REST-y HTTP handler.
func MakeHTTPHandler(ctx context.Context, e Endpoints, logger log.Logger) *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	//Get Cart
	r.Methods("GET").Path("/cart/{userID}").Handler(httptransport.NewServer(
		e.GetCartEndpoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	))

	//Post Cart
	r.Methods("POST").Path("/cart").Handler(httptransport.NewServer(
		e.PostCartEndpoint,
		decodePostCartRequest,
		encodeResponse,
		options...,
	))

	//Delete Cart
	r.Methods("DELETE").PathPrefix("/cart/delete/{userID}").Handler(httptransport.NewServer(
		e.DeleteCartEndpoint,
		decodeDeleteCartRequest,
		encodeResponse,
		options...,
	))

	//Post Product
	r.Methods("POST").PathPrefix("/cart/{userID}/products").Handler(httptransport.NewServer(
		e.PostProductEndpoint,
		decodePostProductRequest,
		encodeResponse,
		options...,
	))

	//Delete Product
	r.Methods("DELETE").PathPrefix("/cart/{prodID}/{userID}").Handler(httptransport.NewServer(
		e.DeleteProductEndpoint,
		decodeDeleteProductRequest,
		encodeResponse,
		options...,
	))

	return r
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/hal+json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":       err.Error(),
		"status_code": code,
		"status_text": http.StatusText(code),
	})
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	d := GetCartRequest{}
	userid := mux.Vars(r)["userID"]
	if len(userid) < 5 {
		return d, ErrInvalidRequest
	}
	d.ID = userid
	return d, nil
}
func decodePostCartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	reg := postCartRequest{}
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

func decodeDeleteCartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	d := deleteCartRequest{}
	userid := mux.Vars(r)["userID"]
	if len(userid) < 3 {
		return d, ErrInvalidRequest
	}
	d.UserID = userid
	return d, nil
}

func decodePostProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	reg := postProductRequest{}
	userid := mux.Vars(r)["userID"]
	if len(userid) < 5 {
		return reg, ErrInvalidRequest
	}
	err := json.NewDecoder(r.Body).Decode(&reg)
	reg.UserID = userid
	if err != nil {
		return nil, err
	}
	return reg, nil
}

func decodeDeleteProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	prodid := mux.Vars(r)["prodID"]
	userid := mux.Vars(r)["userID"]
	g := deleteProductRequest{}
	if len(prodid) < 5 && len(userid) < 5 {
		return g, ErrInvalidRequest
	}
	g.ProductID = prodid
	g.UserID = userid
	return g, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	// All of our response objects are JSON serializable, so we just do that.
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
