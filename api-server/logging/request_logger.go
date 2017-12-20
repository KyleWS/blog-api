package logging

import (
	"net/http"

	logrus "github.com/sirupsen/logrus"
)

const headerAuthorization = "Authorization"
const requestIdHeader = "RequestID"

func RequestLogger(w http.ResponseWriter, r *http.Request) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"method":      r.Method,
		"uri":         r.URL.RequestURI(),
		"auth_header": w.Header().Get(requestIdHeader),
		"agent":       r.UserAgent(),
		"RequestID":   w.Header().Get(requestIdHeader),
	})
}
