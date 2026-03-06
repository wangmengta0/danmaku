package model

import "time"

type SendMessageMQ struct {
	MsgId      string    `gorm:"not null;column:msg_id"`
	VideoId    int       `gorm:"not null;column:video_id"`
	UserId     int       `gorm:"not null;column:user_id"`
	Content    string    `gorm:"not null;column:content"`
	VideoTime  int       `gorm:"not null;column:video_time"`
	CreateTime time.Time `gorm:"not null;column:create_time"`
}
