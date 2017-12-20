package sessions

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
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
		logrus.WithFields(logrus.Fields{
			"access_token": token.AccessToken,
			"token_expire": token.Expiry,
		}).Warn("invalid access token used in request")
		return fmt.Errorf("error validating access token")
	}
	return nil
}
