package handlers

import "github.com/KyleWS/blog-api/api-server/models"

type Ctx struct {
	PostStore *models.MongoStore
}
