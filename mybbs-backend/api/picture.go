package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mybbs-backend/common"
	"mybbs-backend/config"
	"mybbs-backend/dto"
	"mybbs-backend/logic"
)

func UploadPicture(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	var pictureDTO dto.PictureUploadDTO
	if err := ctx.ShouldBind(&pictureDTO); err != nil {
		config.ValidateError(ctx, err)
		return
	}

	fileHeader := pictureDTO.File
	picture, err := logic.UploadPicture(ctx, userID.(int64), fileHeader)
	if err != nil {
		zap.L().Error("picture upload failed", zap.Error(err))
		common.FailByMsg(ctx, "图片上传失败")
		return
	}

	common.Success(ctx, gin.H{
		"picture_id": picture.PictureID,
		"file_path":  picture.Filepath,
	})
}

func DeletePicture(ctx *gin.Context) {
	var pictureIDParam dto.PictureIDParam
	if err := ctx.ShouldBindUri(&pictureIDParam); err != nil {
		config.ValidateError(ctx, err)
		return
	}

	if err := logic.DeletePicture(ctx, pictureIDParam.PictureID); err != nil {
		zap.L().Error("picture deletion failed", zap.Error(err))
		common.FailByMsg(ctx, "图片删除失败")
		return
	}

	common.Success(ctx, "图片删除成功")
}

func GetAllPictures(ctx *gin.Context) {
	pictures, err := logic.GetPictures(ctx)
	if err != nil {
		zap.L().Error("retrieving pictures failed", zap.Error(err))
		common.FailByMsg(ctx, "获取图片失败")
		return
	}

	common.Success(ctx, pictures)
}
