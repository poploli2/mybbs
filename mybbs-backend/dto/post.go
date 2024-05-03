package dto

type PostDTO struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type PostListQuery struct {
	Page     int    `json:"page" form:"page" binding:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" binding:"min=1"`
	Order    string `json:"order" form:"order"`
}

// UpdatePostDTO includes attributes needed to update a post
type UpdatePostDTO struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// PostIDParam includes attributes needed to identify a post
type PostIDParam struct {
	PostID int64 `uri:"post_id" binding:"required"`
}
