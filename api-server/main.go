package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KyleWS/blog-api/api-server/handlers"
	"github.com/KyleWS/blog-api/api-server/models"
	"github.com/KyleWS/blog-api/api-server/sessions"
	cache "github.com/patrickmn/go-cache"
	logrus "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	mgo "gopkg.in/mgo.v2"
)

const (
	apiSignIn = "/oauth/signin"
	apiReply  = "/oauth/reply"
)

func main() {
	// Logging options
	logLevel := os.Getenv("LOG_LEVEL")
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	if logLevel == "DEBUG" {
		logrus.SetLevel(logrus.DebugLevel)
	} else if logLevel == "INFO" {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	// Environment variables
	addr := os.Getenv("ADDR")
	gitRedir := os.Getenv("GITREDIR")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tls_cert := os.Getenv("TLS_CERT")
	tls_secret := os.Getenv("TLS_SECRET")
	dbaddr := os.Getenv("DBADDR")
	dbName := os.Getenv("POSTS_DB_NAME")
	colName := os.Getenv("POSTS_COLLECTION_NAME")

	sess, err := mgo.Dial(dbaddr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"dbaddr":  dbaddr,
			"dbName":  dbName,
			"colName": colName,
		}).Fatal("error connecting to db")
	}

	postStore := models.NewMongoStore(sess, dbName, colName)
	sessionStore := sessions.NewMemStore(120*time.Minute, time.Minute)

	// Used to authenticate with Github
	// below is so I can run locally and in deployment
	if len(addr) > 0 {
		gitRedir = addr
	}
	githubCtx := &sessions.GithubContext{
		OauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"read:user"},
			RedirectURL:  "https://" + gitRedir + apiReply,
			Endpoint:     github.Endpoint,
		},
		StateCache:   cache.New(5*time.Minute, 10*time.Second),
		SessionCache: sessionStore,
	}
	// Used to verify every request user makes to API
	reqCtx := handlers.ReqCtx{
		PostStore:    postStore,
		SessionStore: sessionStore,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(apiSignIn, githubCtx.OAuthSignInHandler)
	mux.HandleFunc(apiReply, githubCtx.OAuthReplyHandler)
	mux.HandleFunc("/post/", reqCtx.PostHandler)
	mux.HandleFunc("/all", reqCtx.AllPostsHandler)
	corsMux := handlers.NewCORS(mux)

	logrus.WithField("addr", addr).Info("blog api server now listening")
	log.Fatal(http.ListenAndServeTLS(addr, tls_cert, tls_secret, corsMux))

}
