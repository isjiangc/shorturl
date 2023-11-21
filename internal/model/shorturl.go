package model

import (
	"gorm.io/gorm"
	"reflect"
)

type OpenType int

const (
	OpenInAll OpenType = iota
	OpenInWechat
	OpenInDingTalk
	OpenInIPhone
	OpenInAndroid
	OpenInIPad
	OpenInSafari
	OpenInChrome
	OpenInFirefox
)

type MemShortUrl struct {
	DestUrl  string
	OpenType OpenType
}

type ShortUrl struct {
	Id       uint   `gorm:"primarykey"`
	ShortUrl string `gorm:"not null"`
	DestUrl  string `gorm:"not null"`
	Valid    bool   `gorm:"not null;default:true"`
	Memo     string
	OpenType OpenType `gorm:"not null;default:0"`
	gorm.Model
}

func (url ShortUrl) IsEmpty() bool {
	return reflect.DeepEqual(url, ShortUrl{})
}

func (u *ShortUrl) TableName() string {
	return "shorturl"
}
