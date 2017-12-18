package handlers

import "net/http"

func (ctx *ReqCtx) PostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// handle getting specific post with provided id
	case http.MethodPost:
		// handle creating a new message
	case http.MethodPatch:
		// handle updating a post
	case http.MethodDelete:
		// handle deleting a specific post
		// this should be locked down
		// because deleting everything is bad.
	}
}

func (ctx *ReqCtx) AllPostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Handles all posts fetch
	}
}
