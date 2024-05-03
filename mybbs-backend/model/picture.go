package model

type Picture struct {
	PictureID int64  `json:"picture_id,string"`
	UserID    int64  `gorm:"not null" json:"user_id,string"`
	Filename  string `gorm:"type:varchar(100);not null" json:"filename"`
	Filepath  string `gorm:"type:varchar(255);not null" json:"filepath"`
	User      User   `gorm:"foreignKey:UserID;references:UserID"` // This is the field name in Picture
	BaseModel
}
