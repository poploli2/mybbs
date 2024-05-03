package dto

import "mime/multipart"

type PictureUploadDTO struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type PictureIDParam struct {
	PictureID int64 `uri:"picture_id" binding:"required"`
}
