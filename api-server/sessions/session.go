package sessions

import "net/http"

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

func GetAuthorization(r *http.Request, ms *MemStore) {

}
