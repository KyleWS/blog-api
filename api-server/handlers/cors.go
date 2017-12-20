package handlers

import (
	"net/http"

	"github.com/KyleWS/blog-api/api-server/logging"
	"gopkg.in/mgo.v2/bson"
)

// CORS struct contains handler that will attach proper Access-Control headers
type CORS struct {
	Handler http.Handler
}

// NewCORS returns cors object with given handler assigned
func NewCORS(handler http.Handler) *CORS {
	return &CORS{
		Handler: handler,
	}
}

func (c *CORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, RequestID")
	w.Header().Add("Access-Control-Expose-Headers", "Authorization, RequestID")
	w.Header().Add("Access-Control-Max-Age", "600")
	if r.Method != "OPTIONS" {
		requestID := bson.NewObjectId().Hex()
		w.Header().Add("RequestID", requestID)
		logging.RequestLogger(w, r).Info("serving request")
		c.Handler.ServeHTTP(w, r)
	}
}
