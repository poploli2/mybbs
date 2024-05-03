package model

type Post struct {
	PostID     int64  `json:"post_id,string"`
	UserID     int64  `gorm:"not null" json:"user_id,string"`
	Title      string `gorm:"type:varchar(100);not null" json:"title"`
	Content    string `gorm:"type:text;not null" json:"content"`
	ClickCount int    `gorm:"type:int;default:0" json:"click_count"`
	User       User   `gorm:"foreignKey:UserID;references:UserID"` // This is the field name in Post
	BaseModel
}
