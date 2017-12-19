package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/KyleWS/blog-api/api-server/models"
	"github.com/KyleWS/blog-api/api-server/sessions"
	"gopkg.in/mgo.v2/bson"
)

func (ctx *ReqCtx) PostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// no authentication required
		// handle getting specific post with provided id
		///// should probably make method for this
		fullPath := r.URL.RequestURI()
		splitPath := strings.Split(fullPath, "/")
		path := splitPath[len(splitPath)-1]
		if path == "" {
			http.Error(w, fmt.Sprintf("error cannot fetch empty path"), http.StatusBadRequest)
			return
		}
		if !bson.IsObjectIdHex(path) {
			http.Error(w, fmt.Sprintf("error path is not valid id"), http.StatusBadRequest)
			return
		}
		bsonID := bson.ObjectIdHex(path)
		post, err := ctx.PostStore.GetTextPostByID(bsonID)
		if err != nil {
			http.Error(w, fmt.Sprintf("error cannot find post with given ID: %v", err), http.StatusBadRequest)
			return
		}
		////// fetch post from db /////////
		json.NewEncoder(w).Encode(post)
	case http.MethodPost:
		// require authenticated user
		if err := sessions.CheckAuthToken(r, ctx.SessionStore); err != nil {
			http.Error(w, fmt.Sprintf("error access token required: %v", err), http.StatusUnauthorized)
			return
		}
		decodedUserTextPost := &models.UserTextPost{}
		if err := json.NewDecoder(r.Body).Decode(decodedUserTextPost); err != nil {
			http.Error(w, fmt.Sprintf("error decoding received json: %v", err), http.StatusBadRequest)
			return
		}
		newTextPost := decodedUserTextPost.GenPostMetaData()
		if err := ctx.PostStore.InsertTextPost(newTextPost); err != nil {
			http.Error(w, fmt.Sprintf("error inserting new post into store: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTextPost)
	case http.MethodPatch:
		// require authenticated user
		if err := sessions.CheckAuthToken(r, ctx.SessionStore); err != nil {
			http.Error(w, fmt.Sprintf("error access token required: %v", err), http.StatusUnauthorized)
			return
		}
		// check that post they want to update exists
		///// should probably make method for this
		fullPath := r.URL.RequestURI()
		splitPath := strings.Split(fullPath, "/")
		path := splitPath[len(splitPath)-1]
		if path == "" {
			http.Error(w, fmt.Sprintf("error cannot fetch empty path"), http.StatusBadRequest)
			return
		}
		if !bson.IsObjectIdHex(path) {
			http.Error(w, fmt.Sprintf("error path is not valid id"), http.StatusBadRequest)
			return
		}
		bsonID := bson.ObjectIdHex(path)
		_, err := ctx.PostStore.GetTextPostByID(bsonID)
		if err != nil {
			http.Error(w, fmt.Sprintf("error cannot find post with given ID: %v", err), http.StatusBadRequest)
			return
		}
		////// fetch post from db /////////
		// handle updating a post
		updates := &models.TextPostUpdates{}
		if err := json.NewDecoder(r.Body).Decode(updates); err != nil {
			http.Error(w, fmt.Sprintf("error decoding received json: %v", err), http.StatusBadRequest)
			return
		}
		updatedPost, err := ctx.PostStore.UpdateTextPost(bsonID, updates)
		if err != nil {
			http.Error(w, fmt.Sprintf("error updating post: %v", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(updatedPost)
	case http.MethodDelete:
		// require authenticated user
		// require authenticated user
		if err := sessions.CheckAuthToken(r, ctx.SessionStore); err != nil {
			http.Error(w, fmt.Sprintf("error access token required: %v", err), http.StatusUnauthorized)
			return
		}
		// check that post they want to update exists
		///// should probably make method for this
		fullPath := r.URL.RequestURI()
		splitPath := strings.Split(fullPath, "/")
		path := splitPath[len(splitPath)-1]
		if path == "" {
			http.Error(w, fmt.Sprintf("error cannot fetch empty path"), http.StatusBadRequest)
			return
		}
		if !bson.IsObjectIdHex(path) {
			http.Error(w, fmt.Sprintf("error path is not valid id"), http.StatusBadRequest)
			return
		}
		bsonID := bson.ObjectIdHex(path)
		_, err := ctx.PostStore.GetTextPostByID(bsonID)
		if err != nil {
			http.Error(w, fmt.Sprintf("error cannot find post with given ID: %v", err), http.StatusBadRequest)
			return
		}
		////// fetch post from db /////////
		// handle deleting a specific post
		// this should be locked down
		// because deleting everything is bad.
		if err := ctx.PostStore.DeletePost(bsonID); err != nil {
			http.Error(w, fmt.Sprintf("error handling delete: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("only accepts GET, POST, PATCH and DELETE"), http.StatusMethodNotAllowed)
	}
}

func (ctx *ReqCtx) AllPostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var allPostsShort []*models.PostShort
		var err error
		if err := sessions.CheckAuthToken(r, ctx.SessionStore); err != nil {
			// get all non-drafts
			allPostsShort, err = ctx.PostStore.FetchAllShort(false)
		} else {
			//get all posts including drafts
			allPostsShort, err = ctx.PostStore.FetchAllShort(true)
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("error fetching all posts: %v", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(allPostsShort)
	default:
		http.Error(w, fmt.Sprintf("only accepts GET"), http.StatusMethodNotAllowed)
	}
}
