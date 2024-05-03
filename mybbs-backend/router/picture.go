package router

import (
	"github.com/gin-gonic/gin"
	"mybbs-backend/api"
	"mybbs-backend/middleware"
)

func GetPictureRoutes(router *gin.RouterGroup) {
	pictureGroup := router.Group("/picture")
	{
		pictureGroup.POST("/upload", middleware.AuthRequired(), api.UploadPicture)
		pictureGroup.DELETE("/:picture_id", middleware.AuthRequired(), api.DeletePicture)
		pictureGroup.GET("/", middleware.AuthRequired(), api.GetAllPictures)
	}
}