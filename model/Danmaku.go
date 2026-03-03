package model

import "time"

type Danmaku struct {
	Id         int       `gorm:"primary_key;AUTO_INCREMENT;column:id"`
	VideoId    int       `gorm:"not null;column:video_id"`
	UserId     int       `gorm:"not null;column:user_id"`
	Content    string    `gorm:"not null;column:content"`
	VideoTime  int       `gorm:"not null;column:video_time"`
	CreateTime time.Time `gorm:"not null;column:create_time"`
}

func (Danmaku) TableName() string {
	return "Danmaku"
}
