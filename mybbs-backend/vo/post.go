package vo

import "mybbs-backend/model"

type PostDetail struct {
	Username string `json:"user_name"`
	model.Post
}
