package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	cache "github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
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

// GithubContext allows us to log in with our github
// account.
type GithubContext struct {
	//oauthConfig is the OAuth configuration for GitHub
	OauthConfig *oauth2.Config
	//stateCache is a cache of previously-generated OAuth state values
	StateCache *cache.Cache
	// sessionCache lets us save the newly authenticated user's
	// access token and github token so we can verify them on the
	// fly
	SessionCache *MemStore
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
	ctx.StateCache.Add(state, nil, cache.DefaultExpiration)
	redirURL := ctx.OauthConfig.AuthCodeURL(state)
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
	if _, found := ctx.StateCache.Get(stateReturned); !found {
		http.Error(w, fmt.Sprintf("invalid state value returned from oauth provider"), http.StatusBadRequest)
		return
	}

	// exchange our code for an access token
	ctx.StateCache.Delete(stateReturned)
	token, err := ctx.OauthConfig.Exchange(oauth2.NoContext, qsParams.Get("code"))
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting access token: %v", err), http.StatusInternalServerError)
		return
	}

	client := ctx.OauthConfig.Client(oauth2.NoContext, token)
	profileRequest, _ := http.NewRequest(http.MethodGet, githubCurrentUserAPI, nil)
	profileRequest.Header.Add(headerAccept, acceptGitHubV3JSON)
	profileResponse, err := client.Do(profileRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting profile : %v", err), http.StatusInternalServerError)
		return
	}
	defer profileResponse.Body.Close()
	// If we have gotten this far, time to save that access token
	ctx.SessionCache.Save(token.AccessToken, token)
	// Give them the access token back so they can start submitting it
	// in their request as an Authoritzation header
	w.Header().Add(headerContentType, profileResponse.Header.Get(headerContentType))
	w.Header().Add(headerAuthorization, token.AccessToken)
	w.WriteHeader(http.StatusAccepted)
	/////////////////////////////////////////////////////////
}
