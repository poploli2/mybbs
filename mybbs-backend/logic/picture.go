package logic

import (
	"context"
	"fmt"
	"mime/multipart"
	"mybbs-backend/config"
	"mybbs-backend/model"
	"mybbs-backend/pkg/snowflake"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadPicture(ctx *gin.Context, userID int64, fileHeader *multipart.FileHeader) (model.Picture, error) {
	db := config.GetDB()

	// 使用雪花算法生成的唯一ID作为文件名
	filename := snowflake.GenerateID()
	filepath := fmt.Sprintf("static/picture/%d.jpg", filename)

	// 保存图片文件到服务器文件系统
	if err := ctx.SaveUploadedFile(fileHeader, filepath); err != nil {
		return model.Picture{}, err
	}

	// 创建图片记录并保存到数据库
	picture := model.Picture{
		UserID:   userID,
		Filename: fileHeader.Filename,
		Filepath: filepath,
	}

	if err := db.Create(&picture).Error; err != nil {
		return model.Picture{}, err
	}

	return picture, nil
}

func DeletePicture(ctx context.Context, pictureID int64) error {
	db := config.GetDB()
	var picture model.Picture

	// 查找并删除数据库记录
	if err := db.Where("picture_id = ?", pictureID).First(&picture).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // 图片不存在则跳过删除
		}
		return err
	}

	// 从文件系统中删除图片文件
	if err := os.Remove(picture.Filepath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// 删除数据库中的图片记录
	if err := db.Delete(&picture).Error; err != nil {
		return err
	}

	return nil
}

func GetPictures(ctx context.Context) ([]model.Picture, error) {
	db := config.GetDB()
	var pictures []model.Picture

	// 获取所有图片记录
	if err := db.Find(&pictures).Error; err != nil {
		return nil, err
	}

	return pictures, nil
}
