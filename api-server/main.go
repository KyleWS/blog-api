package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	cache "github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	headerContentType = "Content-Type"
	headerAccept      = "Accept"
)

const (
	//githubCurrentUserAPI is the URL for GitHub's current user API
	githubCurrentUserAPI = "https://api.github.com/user"

	//acceptGitHubV3JSON is the value you should include in
	//the Accept header when making requests to the GitHub API
	acceptGitHubV3JSON = "application/vnd.github.v3+json"
)

const (
	apiSignIn = "/oauth/signin"
	apiReply  = "/oauth/reply"
)

// GithubContext allows us to log in with our github
// account.
type GithubContext struct {
	//oauthConfig is the OAuth configuration for GitHub
	oauthConfig *oauth2.Config
	//stateCache is a cache of previously-generated OAuth state values
	stateCache *cache.Cache
}

// random value to use as state for oauth
func newStateValue() string {
	buf := make([]byte, 0, 32)
	if _, err := rand.Read(buf); err != nil {
		panic("error generating random bytes")
	}
	return base64.URLEncoding.EncodeToString(buf)
}

// OAuthSignInHandler handles requests for the oauth sign-on API
func (ctx *GithubContext) OAuthSignInHandler(w http.ResponseWriter, r *http.Request) {
	state := newStateValue()
	ctx.stateCache.Add(state, nil, cache.DefaultExpiration)
	redirURL := ctx.oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, redirURL, http.StatusSeeOther)
}

//OAuthReplyHandler handles requests made after authenticating
//with the OAuth provider, and authorizing our application
func (ctx *GithubContext) OAuthReplyHandler(w http.ResponseWriter, r *http.Request) {
	// handle OAutho errors if they ovvured
	qsParams := r.URL.Query()
	if len(qsParams.Get("error")) > 0 {
		errorDescription := qsParams.Get("error_description")
		if len(errorDescription) == 0 {
			errorDescription = "error signing in: " + qsParams.Get("error")
		}
		http.Error(w, fmt.Sprintf("error signing in: %s", errorDescription), http.StatusInternalServerError)
		return
	}

	// check the returned state to make sure it matches
	stateReturned := qsParams.Get("state")
	if _, found := ctx.stateCache.Get(stateReturned); !found {
		http.Error(w, fmt.Sprintf("invalid state value returned from oauth provider"), http.StatusBadRequest)
		return
	}

	// exchange our code for an access token
	ctx.stateCache.Delete(stateReturned)
	token, err := ctx.oauthConfig.Exchange(oauth2.NoContext, qsParams.Get("code"))
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting access token: %v", err), http.StatusInternalServerError)
		return
	}

	client := ctx.oauthConfig.Client(oauth2.NoContext, token)
	profileRequest, _ := http.NewRequest(http.MethodGet, githubCurrentUserAPI, nil)
	profileRequest.Header.Add(headerAccept, acceptGitHubV3JSON)
	profileResponse, err := client.Do(profileRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting profile : %v", err), http.StatusInternalServerError)
		return
	}
	defer profileResponse.Body.Close()
	// TODO: Attach this information so that I can track my being logged in
	// and verify it.
	w.Header().Add(headerContentType, profileResponse.Header.Get(headerContentType))
	io.Copy(w, profileResponse.Body)
}

func main() {
	addr := os.Getenv("ADDR")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tls_cert := os.Getenv("TLS_CERT")
	tls_secret := os.Getenv("TLS_SECRET")

	ctx := &GithubContext{
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"read:user"},
			RedirectURL:  "https://" + addr + apiReply,
			Endpoint:     github.Endpoint,
		},
		stateCache: cache.New(5*time.Minute, 10*time.Second),
	}

	mux := http.NewServeMux()
	mux.HandleFunc(apiSignIn, ctx.OAuthSignInHandler)
	mux.HandleFunc(apiReply, ctx.OAuthReplyHandler)

	log.Printf("blog api server not listening on https://%s...", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tls_cert, tls_secret, mux))

}
