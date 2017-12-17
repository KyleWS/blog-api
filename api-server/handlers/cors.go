package handlers

import "net/http"

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
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Add("Access-Control-Expose-Headers", "Authorization")
	w.Header().Add("Access-Control-Max-Age", "600")
	if r.Method != "OPTIONS" {
		c.Handler.ServeHTTP(w, r)
	}
}
