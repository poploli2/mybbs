package router

import (
	"github.com/gin-gonic/gin"
	"mybbs-backend/middleware"

	"mybbs-backend/api"
)

func GetPostRoutes(router *gin.RouterGroup) {
	communityGroup := router.Group("/post")
	{
		communityGroup.GET("/:post_id", api.GetPostDetail)
		communityGroup.GET("", api.GetTopPostList)
		communityGroup.Use(middleware.AuthRequired())
		communityGroup.POST("", api.CreatePost)
		communityGroup.PATCH("/:post_id", api.UpdatePost)
		communityGroup.DELETE("/:post_id", api.DeletePost)
	}
}
