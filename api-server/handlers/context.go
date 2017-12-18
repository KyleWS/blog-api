package handlers

import (
	"github.com/KyleWS/blog-api/api-server/models"
	"github.com/KyleWS/blog-api/api-server/sessions"
)

type ReqCtx struct {
	PostStore    *models.MongoStore
	SessionStore *sessions.MemStore
}
