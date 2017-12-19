package sessions

import (
	"fmt"
	"net/http"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

func CheckAuthToken(r *http.Request, ms *MemStore) error {
	authHeader := r.Header.Get(headerAuthorization)
	if len(authHeader) == 0 {
		authHeader = r.URL.Query().Get("auth")
	}
	if len(authHeader) == 0 {
		return fmt.Errorf("error access token header missing")
	}
	token, err := ms.Get(authHeader)
	if err != nil {
		return fmt.Errorf("error checking authorization in store: %v", err)
	}
	if !token.Valid() {
		return fmt.Errorf("error validating access token")
	}
	return nil
}
